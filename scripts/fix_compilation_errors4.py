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
    
    # 1. Revert pointer fields
    content = content.replace('record.CheckerID = &checkerID', 'record.CheckerID = checkerID')
    
    remark_fix = '''	remark := req.GetRemark()
	record.Remark = &remark'''
    content = content.replace(remark_fix, 'record.Remark = req.GetRemark()')

    # 2. Fix struct names
    if 'fndm005' in filepath:
        content = content.replace('db.DTAFNDFundCalDate', 'db.TAFNDFundCalDateEdit')
        
    if 'fndm007' in filepath:
        content = content.replace('db.DTAFNDCustGroup', 'db.TAFNDCustGroupEdit')
        
    if 'fndm009' in filepath:
        content = content.replace('db.DTAFNDFavDisc', 'db.TAFNDFavDiscEdit')

    if content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Fixed {filepath}")

