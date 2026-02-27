import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\common\proto\fnd\v1'
proto_files = glob.glob(os.path.join(base_dir, 'fndm*.proto'))

for filepath in proto_files:
    if 'fndm001' in filepath:
        continue # Skip FNDM001

    filename = os.path.basename(filepath)
    module_name = filename.split('.')[0].upper() # e.g. FNDM002
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
        
    original = content

    # 1. Add Approve RPC if not exists
    if 'rpc Approve' not in content:
        # Find the Delete rpc to inject after it
        delete_match = re.search(r'rpc Delete[A-Za-z0-9_]+\s*\(.*?\s*\}\s*\}', content, re.DOTALL)
        if delete_match:
            # Extract the actual entity name (e.g. TAFNDFundFee for FNDM002)
            rpc_name_match = re.search(r'rpc Save([A-Za-z0-9_]+)', content)
            if rpc_name_match:
                entity_name = rpc_name_match.group(1)
                
                approve_rpc = f'''

  rpc Approve{entity_name}(Approve{entity_name}Request) returns (SaveResponse) {{
    option (google.api.http) = {{ post: "/TAapi/Fund/{entity_name}/Approve" body: "*" }};
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {{
      tags: ["API{module_name}:{entity_name}"]
    }};
  }}'''
                content = content[:delete_match.end()] + approve_rpc + content[delete_match.end():]

    # 2. Add ApproveRequest message if not exists
    if f'message Approve' not in content:
        rpc_name_match = re.search(r'rpc Save([A-Za-z0-9_]+)', content)
        if rpc_name_match:
            entity_name = rpc_name_match.group(1)
            
            req_msg = f'''
message Approve{entity_name}Request {{
  string sys_co_id = 1 [json_name="SysCoID"];
  string data_id = 2 [json_name="DataID"];
  string checker_id = 3 [json_name="CheckerID"];
  bool is_approved = 4 [json_name="IsApproved"];
  string remark = 5 [json_name="Remark"];
}}
'''
            content += req_msg

    if content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Updated {filepath}")

