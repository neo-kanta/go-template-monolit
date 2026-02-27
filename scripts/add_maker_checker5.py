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

    # Use regex to find DiffColumns ... } and replace with DiffColumns ... \n\tmodels.MakerCheckerFields\n}
    # It looks like:
    # 	DiffColumns string    gorm:"column:DiffColumns;type:text;not null;default:''"
    # }
    
    # Let's find any closing brace that follows DiffColumns and insert if models.MakerCheckerFields isn't already there
    
    # We want to do this for every struct that has DiffColumns.
    def repl(m):
        full_match = m.group(0)
        diff_col_line = m.group(1) # up to the end of the line
        if 'models.MakerCheckerFields' not in full_match:
            return diff_col_line + '\n\tmodels.MakerCheckerFields\n}'
        return full_match

    content = re.sub(
        r'(DiffColumns\s+string\s+gorm:"[^"]+")\s*\}',
        repl,
        content
    )

    if content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        modified_count += 1
        print(f"Injected into {filepath}")

print(f"Modified {modified_count} models files")
