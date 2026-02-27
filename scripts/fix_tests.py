import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'

# Fix FNDM005 test
fndm005_test = os.path.join(base_dir, 'fndm005', 'service_test.go')
if os.path.exists(fndm005_test):
    with open(fndm005_test, 'r', encoding='utf-8') as f:
        content = f.read()
    # Let's just remove the Remark field from the struct literal
    content = re.sub(r'Remark:\s+[^,]+,', '', content)
    with open(fndm005_test, 'w', encoding='utf-8') as f:
        f.write(content)
    print("Fixed FNDM005 test")

# Fix FNDM008 test
fndm008_test = os.path.join(base_dir, 'fndm008', 'service_test.go')
if os.path.exists(fndm008_test):
    with open(fndm008_test, 'r', encoding='utf-8') as f:
        content = f.read()
    # Let's just remove the Remark field from the struct literal
    content = re.sub(r'Remark:\s+[^,]+,', '', content)
    with open(fndm008_test, 'w', encoding='utf-8') as f:
        f.write(content)
    print("Fixed FNDM008 test")

