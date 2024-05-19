#! /usr/bin/env nix-shell
#! nix-shell -i python3 -p python3 python312Packages.pyyaml

import yaml
import argparse
import sys

markhor_params_field = "markhorParams"
default_separator = "/"

def add_markhor_field(input_file, output_file, separator, reorder):
    with open(input_file, "r") as file:
        data = yaml.safe_load(file)

    custom_separator = separator != default_separator
    kind_field = "kind"
    if kind_field not in data or data[kind_field] != "Secret":
        print(f'Error: Input YAML file must have a field {kind_field} with value "Secret"')
        sys.exit(1)
    else:
        data[kind_field] = "MarkhorSecret"

    api_version_field = "apiVersion"
    if api_version_field not in data or data[api_version_field] != "v1":
        print(f'Error: Input YAML file must have a field {api_version_field} with value "v1"')
        sys.exit(1)
    else:
        data[api_version_field] = "markhork8s.github.io/v1"

    ordering = []
    def add_keys_to_ordering(data, parent_key=""):
        for key, value in data.items():
            if isinstance(value, dict):
                add_keys_to_ordering(value, f"{parent_key}{separator}{key}" if parent_key else key)
            else:
                if separator in key:
                    print(f'Error: Key "{key}" contains the separator character "{separator}". Please use another separator (--separator).')
                    sys.exit(1)
                ordering.append(f"{parent_key}{separator}{key}" if parent_key else key)

    add_keys_to_ordering(data)

    ordering.append(f"{markhor_params_field}{separator}ordering")
    if custom_separator:
        ordering.append(f"{markhor_params_field}{separator}customSeparator")

    if reorder:
        ordering.sort()

    data[markhor_params_field] = {"order": ordering}
    if custom_separator:
        data[markhor_params_field]["customSeparator"] = separator

    if output_file:
        with open(output_file, "w") as file:
            yaml.dump(data, file, default_flow_style=False, sort_keys=reorder)
    else:
        print(yaml.dump(data, default_flow_style=False, sort_keys=reorder))

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=f'Utility script to convert a Kubernetes Secret to a MarkhorSecret')
    parser.add_argument("-a", "--order-alphabetically", dest='reorder', action='store_true', default=False, help="If specified, the fields in the YAML are reordered alphabetically (else, their order is preserved)")
    parser.add_argument("-i", "--input", type=str, required=True, help="Input YAML file")
    parser.add_argument("-o", "--output", type=str, help="Output YAML file path (default is stdout)")
    parser.add_argument("-s", "--separator", type=str, default=default_separator, help=f'Separator for nesting keys (default is "{default_separator}")')

    args = parser.parse_args()

    add_markhor_field(args.input, args.output, args.separator, args.reorder)
