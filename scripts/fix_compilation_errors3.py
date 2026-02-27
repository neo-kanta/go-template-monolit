import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'

# Fix Handler.go errors: h.Service. -> h.svc.
handler_files = glob.glob(os.path.join(base_dir, 'fndm*', 'handler.go'))
for filepath in handler_files:
    if 'fndm001' in filepath:
        continue
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    new_content = content.replace('h.Service.', 'h.svc.')
    new_content = new_content.replace('h.service.', 'h.svc.')
    
    if new_content != content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print(f"Fixed handler in {filepath}")

# Fix Service.go errors: db.DTA... vs db.TA...
service_files = glob.glob(os.path.join(base_dir, 'fndm*', 'service.go'))

def extract_edit_struct(filepath):
    # Peek into db/models.go to see what the actual Edit struct is called
    module_name = filepath.split(os.sep)[-2] # e.g. fndm007
    models_path = os.path.join(base_dir, module_name, 'db', 'models.go')
    if os.path.exists(models_path):
        with open(models_path, 'r', encoding='utf-8') as mf:
            mcontent = mf.read()
            # find first Edit struct
            match = re.search(r'type\s+([A-Za-z0-9_]+Edit)\s+struct', mcontent)
            if match:
                return match.group(1)
    return None


for filepath in service_files:
    if 'fndm001' in filepath:
        continue
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original = content
        
    edit_struct = extract_edit_struct(filepath)
    if edit_struct:
        # Our script blindly injected db.D{entity_name}Edit
        # Let's find any db.D...Edit and replace it with db.{edit_struct}
        # e.g., db.DTAFNDCustGroupEdit -> db.TAFNDCustGroupEdit
        
        # In service.go we injected: var record db.D...Edit
        # We can just look for 'var record db.' and replace whatever Edit struct is there
        
        content = re.sub(r'var\s+record\s+db\.[A-Za-z0-9_]+Edit', f'var record db.{edit_struct}', content)

    if content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Fixed service db edit struct in {filepath}")

