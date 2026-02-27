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

    # Skip if already embedded
    if 'models.MakerCheckerFields' in content:
        continue

    # 1. Add import
    if 'github.com/shopspring/decimal' in content and 'go-transfer-agent/common/platform/model' not in content:
        content = content.replace(
            '\"time\"\n)',
            'models \"go-transfer-agent/common/platform/model\"\n\n\t\"time\"\n)'
        )
    elif 'go-transfer-agent/common/platform/model' not in content:
        content = content.replace(
            'import (\n',
            'import (\n\tmodels "go-transfer-agent/common/platform/model"\n'
        )

    # 2. Add to AuditFields
    target_str = 'DiffColumns string    gorm:"column:DiffColumns;type:text;not null;default:\'\'"\n}'
    replacement_str = 'DiffColumns string    gorm:"column:DiffColumns;type:text;not null;default:\'\'"\n\tmodels.MakerCheckerFields\n}'
    
    if target_str in content:
        content = content.replace(target_str, replacement_str)

    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    
    modified_count += 1

print(f'Modified {modified_count} models.go files to embed MakerCheckerFields in AuditFields')
