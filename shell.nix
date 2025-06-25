{
  # This value controls the default version of the bky-as cli that is setup in
  # the development shell. On release branches, such as "release/v0.1.0-beta.4"
  # this value should be "v0.1.0-beta.4".  On main, it should be set to
  # "latest".
  #
  # This default value can be overwritten from the command line by using a valid
  # semver tag to grab a stable version, a git commit to grab a specific unstable
  # version, or "latest" to grab the latest unstable version, e.g.
  #   `nix-shell --argstr version v0.1.0-beta.5`
  #   `nix-shell --argstr version <full git commit sha>`
  #   `nix-shell --argstr version latest`
  # or use the default value by omitting the argument, e.g.
  #   `nix-shell`
  version ? "latest",
}:
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

  version = version;

  devDependencies = [
    pkgs.git # for project management
    pkgs.gnumake # for project management
    pkgs.go # for prepping for building wasm
    pkgs.golangci-lint # for linting go files
    pkgs.gotools # for tools like goimports
    pkgs.jq # for processing data in examples
    pkgs.nixfmt-rfc-style # for formatting nix files
    pkgs.nodejs_22 # for on chain examples
    pkgs.docker # for building wasm
    pkgs.mo # for stamping version
    pkgs.ca-certificates # for npm install
  ];
}
