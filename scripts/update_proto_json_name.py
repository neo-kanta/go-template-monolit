import os
import re

def to_pascal_case(snake_str):
    components = snake_str.split('_')
    res = ''
    for c in components:
        if c.lower() == 'id':
            res += 'ID'
        elif c.lower() == 'co':
            res += 'Co'
        elif c.lower() == 'prt':
            res += 'Prt'
        elif c.lower() == 'dtl':
            res += 'Dtl'
        elif c.lower() == 'ag':
            res += 'AG'
        elif c.lower() == 'fh':
            res += 'FH'
        elif c.lower() == 'cdsc':
            res += 'CDSC'
        elif c.lower() == 'rp':
            res += 'RP'
        elif c.lower() == 'of':
            res += 'Of'
        elif c.lower() == 'is':
            res += 'Is'
        elif c.lower() == 'fee':
            res += 'Fee'
        elif c.lower() == 'amt':
            res += 'Amt'
        elif c.lower() == 'nm':
            res += 'Nm'
        elif c.lower() == 'op':
            res += 'OP'
        elif c.lower() == 'subs':
            res += 'Subs'
        elif c.lower() == 'rcv':
            res += 'Rcv'
        elif c.lower() == 'txn':
            res += 'Txn'
        elif c.lower() == 'pct':
            res += 'Pct'
        elif c.lower() == 'sw':
            res += 'Sw'
        elif c.lower() == 'pg':
            res += 'PG'
        elif c.lower() == 'tm':
            res += 'TM'
        else:
            res += c.capitalize()
    return res

proto_dir = r'c:\Users\kanta\source\repos\go-transfer-agent\common\proto\fnd\v1'
for filename in os.listdir(proto_dir):
    if not filename.endswith('.proto'):
        continue
    filepath = os.path.join(proto_dir, filename)
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    def repl(match):
        indent = match.group(1)
        attr_type = match.group(2)
        attr_name = match.group(3)
        equal_num = match.group(4)
        
        if 'json_name' in equal_num:
            return match.group(0)
            
        pascal = to_pascal_case(attr_name)
        if attr_name == 'is_switch_in': pascal = 'IsSwitchIn'
        if attr_name == 'sw_i_fund_cry': pascal = 'SwIFundCry'
        if attr_name == 'sw_o_fund_cry_set': pascal = 'SwOFundCrySet'
        if attr_name == 'is_m_rcv_txn_type': pascal = 'IsMRcvTxnType'
        if attr_name == 'is_m_agent_op_type': pascal = 'IsMAgentOPType'
        
        # Remove trailing semicolon
        equal_num = equal_num.rstrip(';').strip()
        return f'{indent}{attr_type} {attr_name} {equal_num} [json_name="{pascal}"];'

    new_content = re.sub(r'(?m)^(\s+)(repeated\s+[\w\.]+|optional\s+[\w\.]+|[\w\.]+)\s+(\w+)\s*(=[^;]+;)', repl, content)
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(new_content)
    print(f"Updated {filename}")
print('Done matching proto fields!')
