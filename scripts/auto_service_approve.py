import os
import glob
import re

base_dir = r'C:\Users\kanta\source\repos\go-transfer-agent\services\fnd'
service_files = glob.glob(os.path.join(base_dir, 'fndm*', 'service.go'))

for filepath in service_files:
    if 'fndm001' in filepath:
        continue # Skip FNDM001

    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
        
    original = content
        
    # Inject context and models imports
    if 'go-transfer-agent/services/fnd/shared' not in content:
        content = content.replace(
            '"go-transfer-agent/common/gen/fnd/v1"',
            '"go-transfer-agent/common/gen/fnd/v1"\n\t"go-transfer-agent/services/fnd/shared"'
        )
    if 'go-transfer-agent/common/platform/model' not in content:
        content = content.replace(
            '"gorm.io/gorm"',
            '"go-transfer-agent/common/platform/model"\n\t"gorm.io/gorm"'
        )
        
    # Fix Save endpoints to set MakerID and Status
    # Example signature: func (s *Service) SaveTAFNDFundFee(ctx context.Context, req *fndv1.TAFNDFundFeeRequest) (*fndv1.SaveResponse, error) {
    save_matches = re.finditer(r'func \(s \*Service\) Save([A-Za-z0-9_]+)\(ctx context\.Context, req \*fndv1\.([A-Za-z0-9_]+Request)\) \(\*fndv1\.SaveResponse, error\) \{', content)
    
    for match in save_matches:
        entity_name = match.group(1)
        req_type = match.group(2)
        
        # We need to find the record creation part inside this function:
        # record := db.DTAFNDFundFeeEdit{}
        # ...
        # if err := s.db.WithContext(ctx).Create(&record).Error;
        
        record_init_pattern = rf'record := db\.D{entity_name}Edit{{}}'
        record_match = re.search(record_init_pattern, content)
        if record_match and 'record.MakerID' not in content[record_match.start():record_match.start()+500]:
            injection = f'''
	// Auth Context & 4-Eyes Principle
	record.MakerID = shared.GetUsernameFromCtx(ctx)
	record.Status = string(model.StatusPendingApproval)
'''
            # inject after record initialization
            insert_pos = record_match.end()
            content = content[:insert_pos] + injection + content[insert_pos:]
            
    # Add Approve endpoints if they don't exist
    rpc_name_match = re.search(r'func \(s \*Service\) Save([A-Za-z0-9_]+)', content)
    if rpc_name_match:
        entity_name = rpc_name_match.group(1)
        
        if f'func (s *Service) Approve{entity_name}' not in content:
            approve_func = f'''
// Approve{entity_name} approves or rejects a pending record for the 4-Eyes Principle.
func (s *Service) Approve{entity_name}(ctx context.Context, req *fndv1.Approve{entity_name}Request) (*fndv1.SaveResponse, error) {{
	s.log.Info("Approve{entity_name} called", slog.String("data_id", req.GetDataId()))

	var record db.D{entity_name}Edit
	result := s.db.WithContext(ctx).Where( + "\"DataID\" = ?" + , req.GetDataId()).First(&record)
	
	if result.Error != nil {{
		if result.Error == gorm.ErrRecordNotFound {{
			return &fndv1.SaveResponse{{Success: false, Message: "Record not found", ReturnCode: "NOT_FOUND"}}, nil
		}}
		s.log.Error("Approve{entity_name} DB error", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}}, nil
	}}

	// Only allow approval if status is PENDING_APPROVAL
	if record.Status != string(model.StatusPendingApproval) {{
		return &fndv1.SaveResponse{{Success: false, Message: "Record is not in PENDING_APPROVAL status", ReturnCode: "INVALID_STATUS"}}, nil
	}}

	// 4-Eyes Principle constraint: Maker != Checker
	checkerID := req.GetCheckerId()
	if checkerID == "" {{
		checkerID = shared.GetUsernameFromCtx(ctx)
	}}

	if record.MakerID != "" && checkerID != "" && record.MakerID == checkerID {{
		return &fndv1.SaveResponse{{Success: false, Message: "4-Eyes Principle Violation: Maker cannot be the Checker", ReturnCode: "FOUR_EYES_VIOLATION"}}, nil
	}}

	// Process Approval/Rejection
	if req.GetIsApproved() {{
		record.Status = string(model.StatusApproved)
	}} else {{
		record.Status = string(model.StatusRejected)
	}}
	
	record.CheckerID = checkerID
	record.Remark = req.GetRemark()

	if err := s.db.WithContext(ctx).Save(&record).Error; err != nil {{
		s.log.Error("Failed to update status", slog.Any("error", err))
		return &fndv1.SaveResponse{{Success: false, Message: "Failed to update record", ReturnCode: "DB_ERROR"}}, nil
	}}

	return &fndv1.SaveResponse{{Success: true, Message: "Record reviewed successfully", ReturnCode: "0000"}}, nil
}}
'''
            content += approve_func
            
    if content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Updated {filepath}")

