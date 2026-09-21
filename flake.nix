{
  description = "CoreHub Go Backend Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          # Go Toolchain
          go
          gopls              # Go Language Server (LSP)
          gotools            # Go utilities (goimports, etc.)
          go-tools           # Staticcheck analyzer

          # Live-reloading & Debugging
          air                # Live-reloading daemon for Go apps
          delve              # Go debugger

          # Database & Tooling
          sqlite             # Local database CLI
          golangci-lint      # Comprehensive linter
          direnv
        ];

        shellHook = ''
          export GOPATH="$PWD/.go"
          export PATH="$GOPATH/bin:$PATH"
          echo "🐹 Go development environment loaded."
        '';
      };
    };
}
