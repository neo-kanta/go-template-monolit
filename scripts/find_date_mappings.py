import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'
go_files = glob.glob(os.path.join(base_dir, 'fndm*', '*.go'))

for filepath in go_files:
    if 'fndm001' in filepath:
        continue # Skip FNDM001

    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
        
    lines = content.split('\n')
    
    for i, line in enumerate(lines):
        if 'time.Parse(' in line or 'Get' in line and ('Date' in line or 'Time' in line):
            print(f"{filepath}:{i+1}: {line.strip()}")

