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

    # Rollback the broken insert from the previous script
    broken_pattern = r'\n\n  rpc Approve[A-Za-z0-9_]+\(Approve[A-Za-z0-9_]+Request\) returns \(SaveResponse\) \{.*?\}\s*\}'
    if re.search(broken_pattern, content, re.DOTALL):
        # We need to target specifically what we injected and remove the extra brace.
        content = re.sub(
            r'  rpc Approve([A-Za-z0-9_]+)\(Approve\1Request\) returns \(SaveResponse\) \{\n    option \(google\.api\.http\) = \{ post: "/TAapi/Fund/\1/Approve" body: "\*" \};\n    option \(grpc\.gateway\.protoc_gen_openapiv2\.options\.openapiv2_operation\) = \{\n      tags: \["API[0-9A-Z]{7}:\1"\]\n    \};\n  \}',
            r'',
            content,
            flags=re.DOTALL
        )
        content = re.sub(r'message Approve[A-Za-z0-9_]+Request \{.*?\}', '', content, flags=re.DOTALL)
        
    # Now let's just re-apply it safely by finding the END of the service block
    rpc_name_match = re.search(r'rpc Save([A-Za-z0-9_]+)', content)
    if not rpc_name_match:
        continue
        
    entity_name = rpc_name_match.group(1)
    
    # 1. Add ApproveRequest message at the very end of the file
    req_msg = f'''
message Approve{entity_name}Request {{
  string sys_co_id = 1 [json_name="SysCoID"];
  string data_id = 2 [json_name="DataID"];
  string checker_id = 3 [json_name="CheckerID"];
  bool is_approved = 4 [json_name="IsApproved"];
  string remark = 5 [json_name="Remark"];
}}
'''
    if f'message Approve{entity_name}' not in content:
        content += req_msg

    # 2. Add Approve RPC to the service block, finding the closing brace of the service
    if f'rpc Approve{entity_name}' not in content:
        service_end = re.search(r'\}\s*// -------------------------------------------------------------------', content)
        if service_end:
            approve_rpc = f'''
  rpc Approve{entity_name}(Approve{entity_name}Request) returns (SaveResponse) {{
    option (google.api.http) = {{ post: "/TAapi/Fund/{entity_name}/Approve" body: "*" }};
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {{
      tags: ["API{module_name}:{entity_name}"]
    }};
  }}
'''
            new_service_end = approve_rpc + service_end.group(0)
            content = content[:service_end.start()] + new_service_end + content[service_end.end():]

    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f"Fixed {filepath}")

