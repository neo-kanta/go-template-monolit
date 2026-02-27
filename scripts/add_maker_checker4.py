import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'
model_files = glob.glob(os.path.join(base_dir, 'fndm*', 'db', 'models.go'))

modified_count = 0

for filepath in model_files:
    if 'fndm001' in filepath:
        continue # Skip FNDM001

    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    original = content

    # Add MakerCheckerFields if not present, right after DiffColumns
    if 'models.MakerCheckerFields' not in content:
        # Match DiffColumns field and the closing brace
        content = re.sub(
            r'(DiffColumns\s+string\s+gorm:"[^"]+"\s*)\}',
            r'\1\tmodels.MakerCheckerFields\n}',
            content
        )

    if content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        modified_count += 1

print(f"Modified {modified_count} models files")
