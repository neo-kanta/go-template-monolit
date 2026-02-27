import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'
handler_files = glob.glob(os.path.join(base_dir, 'fndm*', 'handler.go'))

modified_count = 0

for filepath in handler_files:
    if 'fndm001' in filepath:
        continue # Skip FNDM001

    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
        
    original = content
        
    rpc_name_match = re.search(r'func \(h \*Handler\) Save([A-Za-z0-9_]+)', content)
    if rpc_name_match:
        entity_name = rpc_name_match.group(1)
        
        if f'func (h *Handler) Approve{entity_name}' not in content:
            approve_func = f'''
// Approve{entity_name} handles the Approval POST endpoint.
func (h *Handler) Approve{entity_name}(ctx context.Context, req *fndv1.Approve{entity_name}Request) (*fndv1.SaveResponse, error) {{
	h.log.Info("Approve{entity_name} called", slog.String("data_id", req.GetDataId()))
	return h.service.Approve{entity_name}(ctx, req)
}}
'''
            content += approve_func
            
            # Note: The context used to not be passed to the service
            # We already changed service.go to accept ctx. We need to ensure handlers pass ctx!
            
            # Fix Save passes ctx
            content = re.sub(
                r'return h\.service\.Save([A-Za-z0-9_]+)\((req)\)',
                r'return h.service.Save\1(ctx, \2)',
                content
            )
            # Fix Update passes ctx
            content = re.sub(
                r'return h\.service\.Update([A-Za-z0-9_]+)\((req)\)',
                r'return h.service.Update\1(ctx, \2)',
                content
            )
            # Fix Delete passes ctx
            content = re.sub(
                r'return h\.service\.Delete([A-Za-z0-9_]+)\((req)\)',
                r'return h.service.Delete\1(ctx, \2)',
                content
            )
            # Fix GetData passes ctx
            content = re.sub(
                r'return h\.service\.GetData([A-Za-z0-9_]+)\((req)\)',
                r'return h.service.GetData\1(ctx, \2)',
                content
            )

    if content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        modified_count += 1
        print(f"Updated {{filepath}}")

print(f"Modified {{modified_count}} handler files")
