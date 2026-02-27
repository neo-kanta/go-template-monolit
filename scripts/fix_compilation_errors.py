import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'

# Fix Handler.go errors: h.service -> h.Service
handler_files = glob.glob(os.path.join(base_dir, 'fndm*', 'handler.go'))
for filepath in handler_files:
    if 'fndm001' in filepath:
        continue
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # In C#, it might be a public property, wait go-transfer-agent handlers use h.Service.
    new_content = content.replace('h.service.', 'h.Service.')
    
    if new_content != content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print(f"Fixed handler handler {filepath}")


# Fix Service.go errors: model.StatusPendingApproval -> models.StatusPendingApproval
service_files = glob.glob(os.path.join(base_dir, 'fndm*', 'service.go'))
for filepath in service_files:
    if 'fndm001' in filepath:
        continue
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    new_content = content.replace('model.StatusPendingApproval', 'models.StatusPendingApproval')
    new_content = new_content.replace('model.StatusApproved', 'models.StatusApproved')
    new_content = new_content.replace('model.StatusRejected', 'models.StatusRejected')
    
    # Fix undefined: db.DTAFNDCustGroupEdit in fndm007 and defined db.fooEdit in others
    # Look for "var record db.D[A-Za-z0-9_]+Edit" and see if we can just remove Edit if it's missing from db layer
    
    if new_content != content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print(f"Fixed service models alias in {filepath}")

