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

    # Rollback the broken insert from outside the block
    broken_outside_pattern = r'\}\n\n  \n  rpc Approve([A-Za-z0-9_]+)\(Approve\1Request\) returns \(SaveResponse\) \{\n    option \(google\.api\.http\) = \{ post: "/TAapi/Fund/\1/Approve" body: "\*" \};\n    option \(grpc\.gateway\.protoc_gen_openapiv2\.options\.openapiv2_operation\) = \{\n      tags: \["API_REPLACE"\]\n    \};\n  \}\n'
    content = re.sub(broken_outside_pattern, r'}', content)

    # We also need to extract entity_name again to place it correctly inside
    rpc_name_match = re.search(r'rpc Save([A-Za-z0-9_]+)', content)
    if not rpc_name_match:
        continue
        
    entity_name = rpc_name_match.group(1)

    # Put it back INSIDE the service block before the closing brace '}' that precedes the messages section
    # Usually the end looks like:
    #   }
    # 
    # }
    # 
    # // -------------------------------------------------------------------
    # // FNDM002 Messages
    
    # Let's find the closing brace that comes right before the messages section
    messages_marker = re.search(r'\}\s*// -------------------------------------------------------------------\s*// FNDM[0-9]{3} Messages', content)
    
    if messages_marker:
        approve_rpc = f'''
  rpc Approve{entity_name}(Approve{entity_name}Request) returns (SaveResponse) {{
    option (google.api.http) = {{ post: "/TAapi/Fund/{entity_name}/Approve" body: "*" }};
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {{
      tags: ["API{module_name}:{entity_name}"]
    }};
  }}
'''
        # Replace the } with pprove_rpc followed by }
        # The start of messages_marker is precisely the index of that closing brace.
        content = content[:messages_marker.start()] + approve_rpc + content[messages_marker.start():]

    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f"Fixed {filepath}")

