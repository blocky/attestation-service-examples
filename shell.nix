{
  # This value controls the default version of the bky-as cli that is setup in
  # the development shell. On release branches, such as "release/v0.1.0-beta.4"
  # this value should be "v0.1.0-beta.4".  On main, it should be set to
  # "latest".
  #
  # This default value can be overwritten from the command line by using a valid
  # semver tag to grab a stable version, a git commit to grab a specific unstable
  # version, or "latest" to grab the latest unstable version, e.g.
  #   `nix-shell --argstr asVersion v0.1.0-beta.5`
  #   `nix-shell --argstr asVersion <full git commit sha>`
  #   `nix-shell --argstr asVersion latest`
  # or use the default value by omitting the argument, e.g.
  #   `nix-shell`
  asVersion ? "latest",

  # This value controls the default version of the bky-c cli that is setup in
  # the development shell. On all branches this value should be a valid semver
  # tag pointing to a stable version - like: "v0.1.0-alpha.1".
  #
  # This default value can be overwritten from the command line by using a valid
  # semver tag pointing to a stable version
  #   `nix-shell --argstr cVersion v0.1.0-alpha.1`
  # or use the default value by omitting the argument, e.g.
  #   `nix-shell`
  cVersion ? "latest",
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

  asVersion = asVersion;
  cVersion = cVersion;

  devDependencies = [
    pkgs.git # for project management
    pkgs.gnumake # for project management
    pkgs.go # for prepping for building wasm
    pkgs.golangci-lint # for linting go files
    pkgs.gotools # for tools like goimports
    pkgs.jq # for processing data in examples
    pkgs.nixfmt-rfc-style # for formatting nix files
    pkgs.nodejs_22 # for on chain examples
    pkgs.docker # for building wasm via bky-c
    pkgs.mo # for stamping version
    pkgs.cacert # for npm install
  ];
}
