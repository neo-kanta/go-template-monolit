import os
import re

proto_dir = r'common\proto\fnd\v1'

# We know the main RPC names and tags for each module from our previous scripts.
# Let's map module -> Main Entity Name (like TAFNDFundFee)
modules = {
    'fndm001': 'TAFNDFundInfo',
    'fndm002': 'TAFNDFundFee',
    'fndm003': 'TAFNDSwitch',
    'fndm004': 'TAFNDFundAgent',
    'fndm005': 'TAFNDFundCalDate',
    'fndm006': 'TAFNDRPFeeChgType',
    'fndm007': 'TAFNDCustGroup',
    'fndm008': 'TAFNDPauseTxn',
    'fndm009': 'TAFNDFavDisc',
    'fndm010': 'TAFNDIShareFundFee',
    'fndm011': 'TAFNDIShareFundFeeRdm',
    'fndm012': 'TAFNDPGFundFeeRdm',
    'fndm013': 'TAFNDTMFundFeeRdm'
}

for i in range(1, 14):
    mod = f"fndm{i:03d}"
    entity = modules[mod]
    file_path = os.path.join(proto_dir, f"{mod}.proto")
    
    if not os.path.exists(file_path):
        continue
        
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
        
    # Find the tag
    tag_match = re.search(r'tags:\s*\["([^"]+)"\]', content)
    tag = tag_match.group(1) if tag_match else f"APIFNDM{i:03d}：{entity}"
    
    # We want to add:
    # rpc GetDataByDataID(GetDataRequest) returns (<Entity>Request) {
    #   option (google.api.http) = { get: "/TAapi/Fund/<Entity>/GetData/{data_id}" };
    #   option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
    #     tags: ["<tag>"]
    #   };
    # }
    
    # For FNDM001, we already have some QueryFundInfoByDataID. Let's just add GetDataByDataID anyway.
    
    # Check if already added
    if "GetDataByDataID" in content:
        print(f"Already applied to {mod}")
        continue
        
    # The return type for FNDM001 should probably be SaveFundInfoRequest or QueryFundInfoByDataIDResponse
    # But FNDM001 is special, it uses SaveFundInfoRequest for AUD.
    return_type = f"{entity}Request"
    if mod == 'fndm001':
        return_type = "SaveFundInfoRequest"
        
    rpc_block = f"""
  // GET Data by DataID (Maintain Query)
  rpc GetDataByDataID(GetDataRequest) returns ({return_type}) {{
    option (google.api.http) = {{ get: "/TAapi/Fund/{entity}/GetData/{{data_id}}" }};
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {{
      tags: ["{tag}"]
    }};
  }}
"""
    
    # BaseFund for 004, 005, 006, 007
    if mod in ['fndm004', 'fndm005', 'fndm006', 'fndm007']:
        rpc_block = rpc_block.replace("/TAapi/Fund/", "/TAapi/BaseFund/")
        
    # Insert before the closing brace of the service block
    service_pattern = re.compile(r'(service\s+\w+\s*\{)(.*?)(\n\})', re.DOTALL)
    
    def replacer(match):
        return match.group(1) + match.group(2) + "\n" + rpc_block + match.group(3)
        
    new_content = service_pattern.sub(replacer, content)
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(new_content)
        
    print(f"Injected GetDataByDataID into {mod}")
