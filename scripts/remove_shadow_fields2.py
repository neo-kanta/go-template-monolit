import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'

# 1. Clean all shadow fields robustly
model_files = glob.glob(os.path.join(base_dir, 'fndm*', 'db', 'models.go'))
for filepath in model_files:
    if 'fndm001' in filepath:
        continue

    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    original = content
    
    # Remove lines declaring MakerID, CheckerID, Status, Remark with string or *string
    content = re.sub(r'\n\s+(MakerID|CheckerID|Status|Remark)\s+(?:\*?string)\s+gorm:[^]+', '', content)

    if content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Removed shadow fields via regex in {filepath}")


# 2. Fix FNDM005
svc5 = os.path.join(base_dir, 'fndm005', 'service.go')
with open(svc5, 'r', encoding='utf-8') as f:
    c = f.read()
if 'db.TAFNDFundCalDateEdit' in c:
    c = c.replace('db.TAFNDFundCalDateEdit', 'db.TAFNDFundCalDtlEdit')
    with open(svc5, 'w', encoding='utf-8') as f:
        f.write(c)
    print("Fixed FNDM005 struct name")


