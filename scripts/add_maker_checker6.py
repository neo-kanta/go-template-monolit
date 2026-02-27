import os
import glob

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'
model_files = glob.glob(os.path.join(base_dir, 'fndm*', 'db', 'models.go'))

modified_count = 0

for filepath in model_files:
    if 'fndm001' in filepath:
        continue # Skip FNDM001

    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    original = content
    lines = content.split('\n')
    new_lines = []
    
    for i, line in enumerate(lines):
        if 'models.MakerCheckerFields' in line:
            # If already there, we might run into it
            pass
            
        new_lines.append(line)
        
        # Inject right after DiffColumns
        if 'DiffColumns string' in line and 'gorm:' in line:
            # Check if next line is already MakerCheckerFields
            if i + 1 < len(lines) and 'models.MakerCheckerFields' in lines[i+1]:
                continue
            new_lines.append('\tmodels.MakerCheckerFields')

    new_content = '\n'.join(new_lines)
    if new_content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(new_content)
        modified_count += 1
        print(f"Updated {filepath}")

print(f"Modified {modified_count} models files")
