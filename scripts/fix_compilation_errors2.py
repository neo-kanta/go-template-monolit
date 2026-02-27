import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'

# Fix Service.go errors
service_files = glob.glob(os.path.join(base_dir, 'fndm*', 'service.go'))
for filepath in service_files:
    if 'fndm001' in filepath:
        continue
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    original = content
        
    # 1. Fix string(models.StatusXXX) to models.StatusXXX
    content = content.replace('string(models.StatusPendingApproval)', 'models.StatusPendingApproval')
    content = content.replace('string(models.StatusApproved)', 'models.StatusApproved')
    content = content.replace('string(models.StatusRejected)', 'models.StatusRejected')
    
    # 2. Fix pointers for CheckerID and Remark
    # record.CheckerID = checkerID
    # record.Remark = req.GetRemark()
    content = content.replace('record.CheckerID = checkerID', 'record.CheckerID = &checkerID')
    
    # For Remark, we need a local variable to take address of
    remark_fix = '''	remark := req.GetRemark()
	record.Remark = &remark'''
    content = content.replace('record.Remark = req.GetRemark()', remark_fix)
    
    # 3. Fix missing Edit models:
    # FNDM007: DTAFNDCustGroupEdit -> DTAFNDCustGroup
    # FNDM012: DTAFNDPGFundFeeRdmEdit -> DTAFNDPGFundFeeRdm
    # FNDM005: DTAFNDFundCalDateEdit -> DTAFNDFundCalDate
    # FNDM009: DTAFNDFavDiscEdit -> DTAFNDFavDisc
    # FNDM013: DTAFNDTMFundFeeRdmEdit -> DTAFNDTMFundFeeRdm
    
    if 'fndm007' in filepath:
        content = content.replace('db.DTAFNDCustGroupEdit', 'db.DTAFNDCustGroup')
        
    if 'fndm012' in filepath:
        content = content.replace('db.DTAFNDPGFundFeeRdmEdit', 'db.DTAFNDPGFundFeeRdm')
        
    if 'fndm005' in filepath:
        content = content.replace('db.DTAFNDFundCalDateEdit', 'db.DTAFNDFundCalDate')
        
    if 'fndm009' in filepath:
        content = content.replace('db.DTAFNDFavDiscEdit', 'db.DTAFNDFavDisc')
        
    if 'fndm013' in filepath:
        content = content.replace('db.DTAFNDTMFundFeeRdmEdit', 'db.DTAFNDTMFundFeeRdm')

    if content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Fixed {filepath}")

