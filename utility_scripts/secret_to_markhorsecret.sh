#! /usr/bin/env bash

MP_FIELD="markhorParams"
ORDER_FIELD="order"
DEFAULT_SEPARATOR="/"
SEPARATOR_FIELD="hierarchySeparator"

checkCommands() {
  local commands=("yq" "mktemp")

  for cmd in "${commands[@]}"; do
    command -v "$cmd" >/dev/null || {
      echo "Missing required command: $cmd"
      exit 1
    }
  done
}
checkCommands

while [[ $# -gt 0 ]]; do
  case "$1" in
  -i | --input)
    INPUT_FILE="$2"
    shift 2
    ;;
  -o | --output)
    OUTPUT_FILE="$2"
    shift 2
    ;;
  -s | --separator)
    SEPARATOR="$2"
    shift 2
    ;;
  -v | --version)
    echo "v1.0.0"
    exit 0
    ;;
  -h | --help)
    echo "
Usage:
  -i, --input        the input file (mandatory)
  -o, --output       the output file (default stdout)
  -s, --separator    custom hierarchy separator (default '/')
  -h, --help         display this message
  -v, --version      version of this script
  "
    exit 0
    ;;
  *)
    echo "Unknown option: $1"
    exit 1
    ;;
  esac
done

if [ -z "$INPUT_FILE" ]; then
  echo "Error: Input YAML file is required"
  exit 1
fi
if [ $(yq '.kind' "$INPUT_FILE") != "Secret" ]; then
  echo "Input file must have 'kind: Secret'"
  exit 1
fi
if [ $(yq '.apiVersion' "$INPUT_FILE") != "v1" ]; then
  echo "Input file must have 'apiVersion: v1'"
  exit 1
fi

if [ -z "$SEPARATOR" ]; then
  SEPARATOR="$DEFAULT_SEPARATOR"
fi

ALL_KEYS=$(yq '[.. | select(type != "!!map") | path | join("")] | join("")' "$INPUT_FILE")
if [[ "$ALL_KEYS" == *"${SEPARATOR}"* ]]; then
  echo "Error: one key in the YAML contains the separator character '$SEPARATOR'. Please choose another one (--separator)."
  exit 1
fi

F=$(mktemp)
cp "$INPUT_FILE" $F
yq -i ".$MP_FIELD.$ORDER_FIELD = []" "$F"
if [[ "$SEPARATOR" != "$DEFAULT_SEPARATOR" ]]; then
  yq -i ".$MP_FIELD.$SEPARATOR_FIELD = \"${SEPARATOR}\"" "$F"
fi

yq -i ".$MP_FIELD.$ORDER_FIELD = [.. | select(type != \"!!map\") | path | join(\"$SEPARATOR\")]" "$F"
yq -i '.kind = "MarkhorSecret"' "$F"
yq -i '.apiVersion = "v1"' "$F"

if [ -z "$OUTPUT_FILE" ]; then
  cat "$F"
else
  cat $F >"$OUTPUT_FILE"
fi

rm $F
