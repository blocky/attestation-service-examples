let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-24.11";
  pkgs = import nixpkgs {
    config = { };
    overlays = [ ];
  };

  bky-as = pkgs.stdenv.mkDerivation rec {
    pname = "bky-as";
    version = "v0.1.0-beta.5";

    tmpDir = "/tmp";

    src = builtins.fetchurl {
      url = "https://github.com/blocky/attestation-service-demo/releases/download/${version}/bky-as_linux_amd64";
    };

    unpackPhase = ":";

    installPhase = ''
      install -D -m 555 $src $out/bin/bky-as
    '';
  };
in
pkgs.mkShellNoCC {
  packages = [
    pkgs.git
    pkgs.go
    pkgs.gotools # for tools like goimports
    pkgs.golangci-lint # for linting go files
    pkgs.tinygo # for building wasm
    pkgs.nixfmt-rfc-style # for tools like nix fmt
    pkgs.gnumake # for project management
    pkgs.nodejs_18 # for on chain examples

    bky-as
  ];
}
