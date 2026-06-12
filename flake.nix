{
  description = "Ngent development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          name = "ngent-dev";

          packages = with pkgs; [
            go            # Go 1.26 (supports 1.24 modules)
            nodejs_22     # Node.js ≥ 20 (CI requirement)
            gnumake       # GNU Make
            gofumpt       # Go formatter (used in CI)
          ];

          shellHook = ''
            echo "🎯 ngent dev shell"
            echo "   go    $(go version)"
            echo "   node  $(node --version)"
            echo "   npm   $(npm --version)"
            echo ""
            echo "   make build  — build frontend + Go binary"
            echo "   make run    — build and start server"
            echo "   make test   — run Go tests"
          '';
        };
      }
    );
}
