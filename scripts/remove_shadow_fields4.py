import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'

model_files = glob.glob(os.path.join(base_dir, 'fndm*', 'db', 'models.go'))
modified_count = 0

for filepath in model_files:
    if 'fndm001' in filepath:
        continue

    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    original = content
    
    lines = content.split('\n')
    new_lines = []
    
    for line in lines:
        stripped = line.strip()
        # skip MakerCheckerFields
        if stripped.startswith('models.MakerCheckerFields'):
            new_lines.append(line)
            continue
            
        # Regex match: start with MakerID, 1+ spaces, optionally *, string, 1+ spaces, gorm
        if re.match(r'^(MakerID|CheckerID|Status|Remark)\s+\*?string\s+gorm', stripped):
            continue
            
        new_lines.append(line)

    new_content = '\n'.join(new_lines)
    if new_content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(new_content)
        modified_count += 1
        print(f"Removed shadow fields in {filepath}")

