import os
import glob

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'

# 1. Clean all shadow fields robustly
model_files = glob.glob(os.path.join(base_dir, 'fndm*', 'db', 'models.go'))
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
            
        # check if line is defining one of the 4 fields
        if ('MakerID string' in line or 'MakerID *string' in line) and 'gorm:' in line:
            continue
        if ('CheckerID string' in line or 'CheckerID *string' in line) and 'gorm:' in line:
            continue
        if ('Status string' in line or 'Status *string' in line) and 'gorm:' in line:
            continue
        if ('Remark string' in line or 'Remark *string' in line) and 'gorm:' in line:
            continue
            
        new_lines.append(line)

    new_content = '\n'.join(new_lines)
    if new_content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print(f"Removed shadow fields in {filepath}")

# 2. Fix FNDM005
svc5 = os.path.join(base_dir, 'fndm005', 'service.go')
with open(svc5, 'r', encoding='utf-8') as f:
    c = f.read()
if 'db.TAFNDFundCalDateEdit' in c:
    c = c.replace('db.TAFNDFundCalDateEdit', 'db.TAFNDFundCalDtlEdit')
    with open(svc5, 'w', encoding='utf-8') as f:
        f.write(c)
    print("Fixed FNDM005 struct name")


