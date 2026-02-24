import os
import re

proto_dir = r'c:\Users\kanta\source\repos\go-transfer-agent\common\proto\fnd\v1'
for filename in os.listdir(proto_dir):
    if not filename.endswith('.proto'):
        continue
    filepath = os.path.join(proto_dir, filename)
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Fix option go_package regex
    # It looks like: option go_package = "go-transfer-agent/common/gen/fnd/v1 [json_name="GoPackage"];fndv1";
    # We want to revert it to option go_package = "go-transfer-agent/common/gen/fnd/v1;fndv1";
    content = content.replace(' [json_name="GoPackage"]', '')
    content = content.replace(' [json_name="V1"]', '')
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
print('Fixed go_package options!')
