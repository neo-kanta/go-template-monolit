import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\common\proto\fnd\v1'
proto_files = glob.glob(os.path.join(base_dir, 'fndm*.proto'))

for filepath in proto_files:
    if 'fndm001' in filepath:
        continue # Skip FNDM001
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
        
    original = content

    # The previous script had this block:
    # rpc Approve{entity_name}(Approve{entity_name}Request) returns (SaveResponse) {{
    #   option (google.api.http) = {{ post: "/TAapi/Fund/{entity_name}/Approve" body: "*" }};
    #   option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {{
    #     tags: ["API{module_name}:{entity_name}"]
    #   }};
    # }}
    
    # We accidentally injected bad syntax if the injection site wasn't perfect. Let's fix it by completely replacing the corrupted block.
    # Searching for the broken approve RPC and replacing it with valid syntax.
    
    content = re.sub(
        r'rpc Approve([A-Za-z0-9_]+)\(Approve[A-Za-z0-9_]+Request\) returns \(SaveResponse\) {\s*option \(google\.api\.http\) = { post: "/TAapi/Fund/[A-Za-z0-9_]+/Approve" body: "\*" };\s*option \(grpc\.gateway\.protoc_gen_openapiv2\.options\.openapiv2_operation\) = {\s*tags: \[[^\]]+\]\s*};\s*}',
        r'''
  rpc Approve\1(Approve\1Request) returns (SaveResponse) {
    option (google.api.http) = { post: "/TAapi/Fund/\1/Approve" body: "*" };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      tags: ["API_REPLACE"]
    };
  }''',
        content
    )
    
    # Actually, it's safer to just replace it generally:
    # The syntax error was "common\proto\fnd\v1\fndm002.proto:53:3:syntax error: unexpected "rpc"" which means the previous RPC didn't close its brace properly before the insertion, OR we inserted outside the service block.
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)

