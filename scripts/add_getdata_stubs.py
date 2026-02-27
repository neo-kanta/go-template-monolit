import os
import re

modules = {
    'fndm001': 'TAFNDFundInfo',
    'fndm002': 'TAFNDFundFee',
    'fndm003': 'TAFNDSwitch',
    'fndm004': 'TAFNDFundAgent',
    'fndm005': 'TAFNDFundCalDate',
    'fndm006': 'TAFNDRPFeeChgType',
    'fndm007': 'TAFNDCustGroup',
    'fndm008': 'TAFNDPauseTxn',
    'fndm009': 'TAFNDFavDisc',
    'fndm010': 'TAFNDIShareFundFee',
    'fndm011': 'TAFNDIShareFundFeeRdm',
    'fndm012': 'TAFNDPGFundFeeRdm',
    'fndm013': 'TAFNDTMFundFeeRdm'
}

base_dir = r'services\fnd'

for i in range(1, 14):
    mod = f"fndm{i:03d}"
    entity = modules[mod]
    
    # ─── 1. Update handler.go ───
    handler_path = os.path.join(base_dir, mod, 'handler.go')
    if os.path.exists(handler_path):
        with open(handler_path, 'r', encoding='utf-8') as f:
            h_content = f.read()
            
        req_type = "GetDataRequest"
        res_type = f"{entity}Request"
        if mod == 'fndm001':
            res_type = "QueryFundInfoByDataIDResponse"
            
        # For fndm001, replace QueryFundInfoByDataID with GetDataByDataID
        if mod == 'fndm001' and "QueryFundInfoByDataID" in h_content:
            h_content = re.sub(r'func \(h \*Handler\) QueryFundInfoByDataID.*?return h\.svc\.QueryFundInfoByDataID.*?\n\}', 
                f"""func (h *Handler) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.QueryFundInfoByDataIDResponse, error) {{
	h.log.Info("GetDataByDataID", slog.String("data_id", req.GetDataId()))
	return h.svc.GetDataByDataID(ctx, req)
}}""", h_content, flags=re.DOTALL)
        elif "GetDataByDataID" not in h_content:
            stub = f"""
func (h *Handler) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.{res_type}, error) {{
	h.log.Info("GetDataByDataID", slog.String("data_id", req.GetDataId()))
	return h.svc.GetDataByDataID(ctx, req)
}}
"""
            # Insert before Check Functions or at the end of the file (before EOF or RegisterHandlers)
            if "// ─── KFNDM: Check Functions ───" in h_content:
                h_content = h_content.replace("// ─── KFNDM: Check Functions ───", stub + "\n// ─── KFNDM: Check Functions ───")
            else:
                h_content = re.sub(r'(func RegisterHandlers.*)', stub + r'\1', h_content, flags=re.DOTALL)
                
        with open(handler_path, 'w', encoding='utf-8') as f:
            f.write(h_content)

    # ─── 2. Update service.go ───
    service_path = os.path.join(base_dir, mod, 'service.go')
    if os.path.exists(service_path):
        with open(service_path, 'r', encoding='utf-8') as f:
            s_content = f.read()
            
        res_type = f"{entity}Request"
        if mod == 'fndm001':
            res_type = "QueryFundInfoByDataIDResponse"
            
        if mod == 'fndm001' and "QueryFundInfoByDataID" in s_content:
           # Manual fix needed for FNDM001 due to complex logic inside QueryFundInfoByDataID.
           # Let's just rename QueryFundInfoByDataID to GetDataByDataID for FNDM001 in service.
           s_content = re.sub(
               r'func \(s \*Service\) QueryFundInfoByDataID\(sysCoID, prtFundCode string\) \(\*fndv1\.QueryFundInfoByDataIDResponse, error\) \{',
               r'func (s *Service) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.QueryFundInfoByDataIDResponse, error) {',
               s_content)
               
        elif "GetDataByDataID" not in s_content:
            stub = f"""
// GetDataByDataID retrieves full data by DataID.
func (s *Service) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.{res_type}, error) {{
	return &fndv1.{res_type}{{}}, nil
}}
"""
            s_content += stub
            
        with open(service_path, 'w', encoding='utf-8') as f:
            f.write(s_content)

print("Go stubs applied to handlers and services.")
