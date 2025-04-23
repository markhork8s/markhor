# #! /usr/bin/env nix-shell
# #! nix-shell -i python3 -p python3
 
go get -u && go mod tidy

echo "--------------- Modules updated"
