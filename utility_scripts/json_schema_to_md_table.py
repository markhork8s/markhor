#! /usr/bin/env nix-shell
#! nix-shell -i python3 -p python3

import json
import sys

def extract_properties(data, path=''):
    result = []
    if 'properties' in data:
        for key, value in data['properties'].items():
            new_path = f"{path}.{key}" if path else key
            if 'type' in value and value['type'] != 'object':
                prop_type = value.get('type')
                if prop_type == 'array':
                    item_type = value['items']['type']
                    prop_type = f'array[{item_type}]'
                prop_default = value.get('default', '')
                if prop_type == 'string':
                    prop_default = f"\"{prop_default}\""
                prop_desc = value.get('description', '')
                result.append((new_path, prop_type, prop_default, prop_desc))
            result.extend(extract_properties(value, new_path))
    return result

def json_schema_to_markdown_table(file_path):
    with open(file_path) as file:
        json_data = json.load(file)
    
    properties = extract_properties(json_data)
    
    table_header = "|Property name|Type|Default value|Description|\n|-|-|-|-|"
    table_body = '\n'.join([f"|{prop[0]}|{prop[1]}|{prop[2]}|{prop[3]}|" for prop in properties])
    disclaimer = f"<!-- This table was generated automatically by ./utility_scripts/json_schema_to_md_table.py using the file {file_path} -->"
    markdown_table = "\n".join([disclaimer, table_header, table_body, disclaimer])
    
    return markdown_table

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Please provide the path to the JSON schema file as an argument.")
    else:
        file_path = sys.argv[1]
        markdown_output = json_schema_to_markdown_table(file_path)
        print(markdown_output)
