import json
import os
import shutil

swagger_path = r'services\fnd\internal\adapter\gateway\fnd.swagger.json'
bruno_dir = r'tests\bruno\fnd-api'

# Load the swagger
with open(swagger_path, 'r', encoding='utf-8') as f:
    swagger = json.load(f)

# Clear existing .bru files in the root of the collection, and module folders
for item in os.listdir(bruno_dir):
    item_path = os.path.join(bruno_dir, item)
    if os.path.isfile(item_path) and item.endswith('.bru'):
        os.remove(item_path)
    elif os.path.isdir(item_path) and item.startswith('FNDM'):
        shutil.rmtree(item_path)

# Extract operations
for path, methods in swagger.get('paths', {}).items():
    for method, op in methods.items():
        # Get the tag for folder organization
        tags = op.get('tags', [])
        tag = tags[0] if tags else 'Default'
        
        # Example tag: APIFNDM001：TAFNDFundInfo -> module: FNDM001
        module_name = tag.split('：')[0].replace('API', '') if '：' in tag else tag
        
        # Determine sequence and folder
        folder_path = os.path.join(bruno_dir, module_name)
        os.makedirs(folder_path, exist_ok=True)
        
        op_id = op.get('operationId', 'Unknown')
        
        # Build bru attributes
        bru_method = method.lower()
        
        # Construct url with base_url variable
        bru_url = f"{{{{base_url}}}}{path}"
        
        # Determine if body is JSON
        body_type = "none"
        if bru_method in ['post', 'put', 'patch']:
            body_type = "json"
        
        # Extract query parameters to append to URL
        query_params = []
        for param in op.get('parameters', []):
            if param.get('in') == 'query':
                query_params.append(f"{param['name']}=")
        
        if query_params:
            bru_url += "?" + "&".join(query_params)
            
        bru_content = f"""meta {{
  name: {op_id}
  type: http
  seq: 1
}}

{bru_method} {{
  url: {bru_url}
  body: {body_type}
  auth: none
}}

headers {{
  Authorization: Bearer {{{{token}}}}
}}
"""
        # Add basic empty JSON body if applicable
        if body_type == "json":
            bru_content += """
body:json {
  {}
}
"""
            
        # Write to file
        file_path = os.path.join(folder_path, f"{op_id}.bru")
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(bru_content)

print("Bruno collection successfully generated!")
