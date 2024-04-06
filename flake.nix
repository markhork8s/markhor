{
  description = "k8s env";

  inputs = {
    nixpkgs.url = "github:nixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config = {
            # allowUnfree = true;
          };
        };
      in
      {
        devShell = with pkgs; mkShell rec {
          #ENV_VARIABLE_1 = "test";
          buildInputs = [
            go
            gopls
            (vscode-with-extensions.override {
              vscode = vscodium;
              vscodeExtensions = with vscode-extensions; [
                jnoortheen.nix-ide
                golang.go
              ]
                #  ++ pkgs.vscode-utils.extensionsFromVscodeMarketplace [
                #   {
                #     name = "codegeex";
                #     publisher = "aminer";
                #     version = "1.0.6";
                #     sha256 = "sha256-q8HSFZRhwZv5zApHsVoyKGqZsDDyUqjxv/qwGAuOE0c=";
                #   }
                # ]
              ;
            })
          ];
        };
      });
}
