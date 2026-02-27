import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'
go_files = glob.glob(os.path.join(base_dir, 'fndm*', '*.go'))

modified_count = 0

for filepath in go_files:
    if 'fndm001' in filepath:
        continue # Skip FNDM001

    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    original_content = content
    
    lines = content.split('\n')
    new_lines = []
    
    for line in lines:
        match = re.search(r'([a-zA-Z0-9_\[\].]+)\s*:\s*decimal\.NewFromFloat\(([^)]+)\)', line)
        
        if match and '.RoundBank' not in line:
            left_side = match.group(1)
            
            precision = 2
            lower_name = left_side.lower()
            if 'rate' in lower_name or 'unit' in lower_name or 'ratio' in lower_name:
                precision = 4
                
            new_line = line.replace(
                f'decimal.NewFromFloat({match.group(2)})', 
                f'decimal.NewFromFloat({match.group(2)}).RoundBank({precision})'
            )
            new_lines.append(new_line)
        else:
            new_lines.append(line)
            
    content = '\n'.join(new_lines)
    
    if content != original_content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        modified_count += 1
        print(f"Updated {filepath}")

print(f'Modified {modified_count} go files to apply decimal RoundBank')
