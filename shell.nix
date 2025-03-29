let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-24.11";
  pkgs = import nixpkgs {
    config = { };
    overlays = [ ];
  };

  mkDevShell = import ./nix/mkDevShell.nix;
in
mkDevShell {
  pkgs = pkgs;

  # this value controls the version of the bky-as cli that is setup in the
  # development shell. On release branches, such as "release/v0.1.0-beta.4"
  # this value should be "v0.1.0-beta.4".  On main, it should be set to
  # "unstable"
  version = "unstable";

  devDependencies = [
    pkgs.git # for project management
    pkgs.gnumake # for project management
    pkgs.go # for prepping for building wasm
    pkgs.golangci-lint # for linting go files
    pkgs.gotools # for tools like goimports
    pkgs.jq # for processing data in examples
    pkgs.nixfmt-rfc-style # for formatting nix files
    pkgs.nodejs_18 # for on chain examples
    pkgs.tinygo # for building wasm
  ];
}
