{
  description = "AS tests";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=24.11";
    flake-utils.url = "github:numtide/flake-utils";

    # this tool allows us use nix-shell and nix shell
    # and is used for our shell.nix
    flake-compat.url = "https://flakehub.com/f/edolstra/flake-compat/1.tar.gz";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.git
            pkgs.go
            pkgs.gotools # for tools like goimports
            pkgs.golangci-lint # for linting go files
            pkgs.tinygo # for building wasm
            pkgs.nixfmt-rfc-style # for tools like nix fmt
            pkgs.gnumake # for project management
          ];
        };
      }
    );
}
