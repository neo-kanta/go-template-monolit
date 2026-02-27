import os
import re

proto_dir = r"c:\Users\kanta\source\repos\go-transfer-agent\common\proto\fnd\v1"

tags_map = {
    "fnd.proto": "APIFNDM001：TAFNDFundInfo",
    "fndm002.proto": "APIFNDM002：TAFNDFundFee",
    "fndm003.proto": "APIFNDM003：TAFNDSwitch",
    "fndm004.proto": "APIFNDM004：TAFNDFundAgent",
    "fndm005.proto": "APIFNDM005：TAFNDFundCalDate",
    "fndm006.proto": "APIFNDM006：TAFNDRPFeeChgType",
    "fndm007.proto": "APIFNDM007：TAFNDFundGroup",
    "fndm008.proto": "APIFNDM008：TAFNDPauseTxn",
    "fndm009.proto": "APIFNDM009：TAFNDFavDisc",
    "fndm010.proto": "APIFNDM010：TAFNDIShareFundFee",
    "fndm011.proto": "APIFNDM011：TAFNDIShareFundFeeRdm",
    "fndm012.proto": "APIFNDM012：TAFNDPGFundFeeRdm",
    "fndm013.proto": "APIFNDM013：TAFNDTMFundFeeRdm",
}

for filename in os.listdir(proto_dir):
    if not filename.endswith(".proto"):
        continue

    tag = tags_map.get(filename)
    if not tag:
        continue

    filepath = os.path.join(proto_dir, filename)
    with open(filepath, "r", encoding="utf-8") as f:
        content = f.read()

    # Add import if missing
    if "protoc-gen-openapiv2/options/annotations.proto" not in content:
        content = content.replace(
            'import "google/api/annotations.proto";',
            'import "google/api/annotations.proto";\nimport "protoc-gen-openapiv2/options/annotations.proto";',
        )

    if filename == "fnd.proto" and "openapiv2_swagger" not in content:
        content = content.replace(
            'option go_package = "go-transfer-agent/common/gen/fnd/v1;fndv1";',
            'option go_package = "go-transfer-agent/common/gen/fnd/v1;fndv1";\n\noption (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {\n  info: {\n    title: "TaBaseAPI";\n    version: "1.0";\n  };\n};',
        )

    def replacer(match):
        http_opt = match.group(1)
        insertion = f'\n    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {{\n      tags: ["{tag}"]\n    }};'
        return http_opt + insertion

    if "openapiv2_operation" not in content:
        content = re.sub(
            r"(option\s*\(google\.api\.http\)\s*=\s*\{.*?\};)",
            replacer,
            content,
            flags=re.DOTALL,
        )

    with open(filepath, "w", encoding="utf-8") as f:
        f.write(content)
print("Tags injected into all proto files with array syntax.")
