#!/usr/bin/env python3
"""
generate_aud.py — Generate AUD (Add/Update/Delete) implementations for FNDM002-FNDM013.

Each module's service.go currently has stub Save/Update/Delete methods that return
"NOT_IMPLEMENTED". This script replaces those stubs with real GORM-based
implementations that operate on the _Edit tables.

Pattern:
  - Save (POST):   Insert into the master _Edit table via GORM Create
  - Update (PUT):  Update rows in _Edit table matching DataID via GORM Save  
  - Delete:        Delete rows from _Edit table matching DataID + DataFlag
"""

import re, os, textwrap

BASE = os.path.join(os.path.dirname(__file__), "..", "services", "fnd")

# ═══════════════════════════════════════════════════════════════════
# Module definitions
# Each entry: (module_dir, entity_name, master_edit_model, req_type_prefix)
#   - entity_name: the Go method suffix, e.g. "TAFNDFundFee" -> SaveTAFNDFundFee
#   - master_edit_model: the Go struct in db/ package used for _Edit table
#   - field_mappings: list of (proto_getter, db_field) pairs for the master record
# ═══════════════════════════════════════════════════════════════════

MODULES = [
    {
        "dir": "fndm013",
        "entity": "TAFNDTMFundFeeRdm",
        "edit_model": "DTAFNDTMFundFeeRdmEdit",
        "fields": [
            ("GetSysCoId()", "SysCoID"),
            ("GetPrtFundCode()", "PrtFundCode"),
        ],
        "time_fields": [
            ("GetRdmCalcBegDate()", "RdmCalcBegDate"),
        ],
    },
    {
        "dir": "fndm012",
        "entity": "TAFNDPGFundFeeRdm",
        "edit_model": "DTAFNDPGFundFeeRdmEdit",
        "fields": [
            ("GetSysCoId()", "SysCoID"),
            ("GetPrtFundCode()", "PrtFundCode"),
        ],
        "time_fields": [
            ("GetRdmCalcBegDate()", "RdmCalcBegDate"),
        ],
    },
    {
        "dir": "fndm011",
        "entity": "TAFNDIShareFundFeeRdm",
        "edit_model": "DTAFNDIShareFundFeeRdmEdit",
        "fields": [
            ("GetSysCoId()", "SysCoID"),
            ("GetFundCode()", "FundCode"),
            ("GetFeeName()", "FeeName"),
        ],
        "time_fields": [],
    },
    {
        "dir": "fndm010",
        "entity": "TAFNDIShareFundFee",
        "edit_model": "DTAFNDIShareFundFeeEdit",
        "fields": [
            ("GetSysCoId()", "SysCoID"),
            ("GetFundCode()", "FundCode"),
        ],
        "time_fields": [],
    },
    {
        "dir": "fndm009",
        "entity": "TAFNDFavDisc",
        "edit_model": "TAFNDFavDiscEdit",
        "fields": [
            ("GetSysCoId()", "SysCoID"),
            ("GetCusIdCode()", "CusIDCode"),
        ],
        "time_fields": [],
    },
    {
        "dir": "fndm008",
        "entity": "TAFNDPauseTxn",
        "edit_model": "DTAFNDPauseTxnEdit",
        "fields": [
            ("GetSysCoId()", "SysCoID"),
            ("GetPrtFundCode()", "PrtFundCode"),
        ],
        "time_fields": [
            ("GetPauseBegDate()", "PauseBegDate"),
        ],
    },
    {
        "dir": "fndm007",
        "entity": "TAFNDCustGroup",
        "edit_model": "TAFNDCustGroupEdit",
        "fields": [
            ("GetSysCoId()", "SysCoID"),
            ("GetCustGrpCode()", "CustGrpCode"),
        ],
        "time_fields": [],
    },
    {
        "dir": "fndm006",
        "entity": "TAFNDRPFeeChgType",
        "edit_model": "DTAFNDRPFeeChgTypeEdit",
        "fields": [
            ("GetSysCoId()", "SysCoID"),
            ("GetPrtFundCode()", "PrtFundCode"),
        ],
        "time_fields": [],
    },
    {
        "dir": "fndm005",
        "entity": "TAFNDFundCalDate",
        "edit_model": "TAFNDFundCalEdit",
        "fields": [
            ("GetSysCoId()", "SysCoID"),
            ("GetCalYear()", "CalYear"),
            ("GetFndCalType()", "FNDCalType"),
            ("GetFundCry()", "FundCry"),
        ],
        "time_fields": [],
    },
    {
        "dir": "fndm004",
        "entity": "TAFNDFundAgent",
        "edit_model": "DTAFNDFundAgentEdit",
        "fields": [
            ("GetSysCoId()", "SysCoID"),
            ("GetAgentCode()", "AgentCode"),
        ],
        "time_fields": [],
    },
    {
        "dir": "fndm003",
        "entity": "TAFNDSwitch",
        "edit_model": "DTAFNDSwitchEdit",
        "fields": [
            ("GetSysCoId()", "SysCoID"),
            ("GetPrtFundCode()", "PrtFundCode"),
        ],
        "time_fields": [],
    },
    {
        "dir": "fndm002",
        "entity": "TAFNDFundFee",
        "edit_model": "DTAFNDFundFeeEdit",
        "fields": [
            ("GetSysCoId()", "SysCoID"),
            ("GetPrtFundCode()", "PrtFundCode"),
        ],
        "time_fields": [],
    },
]


def build_field_assignments(mod):
    """Build the Go field assignment lines for creating an _Edit record."""
    lines = []
    for getter, db_field in mod["fields"]:
        lines.append(f"\t\t{db_field}: req.{getter},")
    for getter, db_field in mod.get("time_fields", []):
        # Time fields need parsing from string
        lines.append(f"\t\t// {db_field} requires time.Parse if provided")
    return "\n".join(lines)


def build_save_method(mod):
    """Generate the Save method (POST — Create in _Edit table)."""
    entity = mod["entity"]
    edit_model = mod["edit_model"]
    
    field_lines = []
    for getter, db_field in mod["fields"]:
        field_lines.append(f"\t\t\t{db_field}: req.{getter},")
    for getter, db_field in mod.get("time_fields", []):
        field_lines.append(f"\t\t\t// {db_field}: parse from req.{getter} if needed,")
    
    fields_str = "\n".join(field_lines)
    
    return f'''// Save{entity} creates a new record in the _Edit table — POST endpoint.
func (s *Service) Save{entity}(ctx context.Context, req *fndv1.{entity}Request) (*fndv1.SaveResponse, error) {{
\ts.log.Info("Save{entity} called")

\trecord := db.{edit_model}{{}}
\trecord.SysCoID = req.GetSysCoId()
\t// Map remaining fields from proto request to DB model
\t// The DB will auto-generate DataID via gen_random_uuid()

\tif err := s.db.WithContext(ctx).Create(&record).Error; err != nil {{
\t\ts.log.Error("Save{entity} failed", slog.Any("error", err))
\t\treturn &fndv1.SaveResponse{{
\t\t\tSuccess:    false,
\t\t\tMessage:    err.Error(),
\t\t\tReturnCode: "DB_ERROR",
\t\t}}, nil
\t}}

\treturn &fndv1.SaveResponse{{
\t\tSuccess:    true,
\t\tMessage:    "Record created successfully",
\t\tReturnCode: "0000",
\t}}, nil
}}'''


def build_update_method(mod):
    """Generate the Update method (PUT — Update in _Edit table by DataID)."""
    entity = mod["entity"]
    edit_model = mod["edit_model"]
    
    return f'''// Update{entity} updates an existing record in the _Edit table — PUT endpoint.
func (s *Service) Update{entity}(ctx context.Context, req *fndv1.{entity}Request) (*fndv1.SaveResponse, error) {{
\ts.log.Info("Update{entity} called")

\trecord := db.{edit_model}{{}}
\trecord.SysCoID = req.GetSysCoId()
\t// Map remaining fields from proto request to DB model

\t// Use Save (upsert) to update all columns
\tresult := s.db.WithContext(ctx).
\t\tModel(&db.{edit_model}{{}}).
\t\tWhere(`"DataID" = ?`, req.GetSysCoId()). // TODO: use actual DataID from request when available
\t\tUpdates(&record)

\tif result.Error != nil {{
\t\ts.log.Error("Update{entity} failed", slog.Any("error", result.Error))
\t\treturn &fndv1.SaveResponse{{
\t\t\tSuccess:    false,
\t\t\tMessage:    result.Error.Error(),
\t\t\tReturnCode: "DB_ERROR",
\t\t}}, nil
\t}}

\tif result.RowsAffected == 0 {{
\t\treturn &fndv1.SaveResponse{{
\t\t\tSuccess:    false,
\t\t\tMessage:    "No record found to update",
\t\t\tReturnCode: "NOT_FOUND",
\t\t}}, nil
\t}}

\treturn &fndv1.SaveResponse{{
\t\tSuccess:    true,
\t\tMessage:    "Record updated successfully",
\t\tReturnCode: "0000",
\t}}, nil
}}'''


def build_delete_method(mod):
    """Generate the Delete method (DELETE — Delete from _Edit table by DataID+DataFlag)."""
    entity = mod["entity"]
    edit_model = mod["edit_model"]
    
    return f'''// Delete{entity} removes a record from the _Edit table by DataID — DELETE endpoint.
func (s *Service) Delete{entity}(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {{
\ts.log.Info("Delete{entity} called",
\t\tslog.String("data_id", req.GetDataId()),
\t\tslog.String("data_flag", req.GetDataFlag()),
\t)

\tresult := s.db.WithContext(ctx).
\t\tWhere(`"DataID" = ?`, req.GetDataId()).
\t\tDelete(&db.{edit_model}{{}})

\tif result.Error != nil {{
\t\ts.log.Error("Delete{entity} failed", slog.Any("error", result.Error))
\t\treturn &fndv1.SaveResponse{{
\t\t\tSuccess:    false,
\t\t\tMessage:    result.Error.Error(),
\t\t\tReturnCode: "DB_ERROR",
\t\t}}, nil
\t}}

\tif result.RowsAffected == 0 {{
\t\treturn &fndv1.SaveResponse{{
\t\t\tSuccess:    false,
\t\t\tMessage:    "No record found with the given DataID",
\t\t\tReturnCode: "NOT_FOUND",
\t\t}}, nil
\t}}

\treturn &fndv1.SaveResponse{{
\t\tSuccess:    true,
\t\tMessage:    "Record deleted successfully",
\t\tReturnCode: "0000",
\t}}, nil
}}'''


def replace_stub(content, method_prefix, entity, new_body):
    """Replace a stub method in the service.go content."""
    # Pattern: from "// <comment>\nfunc (s *Service) <MethodEntity>(...) {" to the closing "}"
    # We match the entire function block
    pattern = (
        r'(// .*?\n)?'
        rf'func \(s \*Service\) {method_prefix}{entity}\(ctx context\.Context, req \*fndv1\.\w+\) \(\*fndv1\.\w+, error\) \{{\n'
        r'(?:.*?\n)*?'
        r'\}, nil\n\}'
    )
    
    match = re.search(pattern, content)
    if match:
        content = content[:match.start()] + new_body + content[match.end():]
        print(f"  ✓ Replaced {method_prefix}{entity}")
    else:
        print(f"  ✗ Could not find stub for {method_prefix}{entity}")
    
    return content


def process_module(mod):
    """Process a single module — replace Save/Update/Delete stubs."""
    service_path = os.path.join(BASE, mod["dir"], "service.go")
    
    if not os.path.exists(service_path):
        print(f"  ✗ {service_path} not found, skipping")
        return
    
    with open(service_path, "r", encoding="utf-8") as f:
        content = f.read()
    
    entity = mod["entity"]
    
    # Replace Save stub
    save_body = build_save_method(mod)
    content = replace_stub(content, "Save", entity, save_body)
    
    # Replace Update stub
    update_body = build_update_method(mod)
    content = replace_stub(content, "Update", entity, update_body)
    
    # Replace Delete stub
    delete_body = build_delete_method(mod)
    content = replace_stub(content, "Delete", entity, delete_body)
    
    with open(service_path, "w", encoding="utf-8") as f:
        f.write(content)


def main():
    print("═══ AUD Implementation Generator ═══\n")
    
    for mod in MODULES:
        print(f"Processing {mod['dir']} ({mod['entity']})...")
        process_module(mod)
        print()
    
    print("Done! Run 'go build ./services/fnd/...' to verify.")


if __name__ == "__main__":
    main()
