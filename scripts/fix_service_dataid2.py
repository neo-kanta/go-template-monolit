import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'
service_files = glob.glob(os.path.join(base_dir, 'fndm*', 'service.go'))

modified_count = 0

for filepath in service_files:
    if 'fndm001' in filepath:
        continue # Skip FNDM001

    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
        
    original = content

    # Replace the broken string
    broken_str = 'Where(""DataID" = ?", req.GetDataId())'
    fixed_str = 'Where("\\"DataID\\" = ?", req.GetDataId())'
    
    if broken_str in content:
        content = content.replace(broken_str, fixed_str)

    if content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        modified_count += 1
        print(f"Updated {filepath}")

print(f"Modified {modified_count} service files")
