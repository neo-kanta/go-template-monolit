import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'
service_files = glob.glob(os.path.join(base_dir, 'fndm*', 'service.go'))

for filepath in service_files:
    if 'fndm001' in filepath:
        continue
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    original = content
    
    # 1. Restore pointer fields
    content = content.replace('record.CheckerID = checkerID', 'record.CheckerID = &checkerID')
    
    remark_fix = '''	remark := req.GetRemark()
	record.Remark = &remark'''
    content = content.replace('record.Remark = req.GetRemark()', remark_fix)
    
    # 2. Fix the DTO mapping in FNDM008
    # In FNDM008, m.Remark (which is now *string) is mapped to Remark: (which is string)
    # 	Remark:        m.Remark, -> Remark: func() string { if m.Remark != nil { return *m.Remark } return "" }(),
    if 'fndm008' in filepath:
        content = content.replace('Remark:        m.Remark,', 'Remark: func() string { if m.Remark != nil { return *m.Remark }; return "" }(),')
        
    # Same thing in FNDM005
    if 'fndm005' in filepath:
        content = content.replace('Remark:     m.Remark,', 'Remark: func() string { if m.Remark != nil { return *m.Remark }; return "" }(),')

    if content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Fixed {filepath}")

