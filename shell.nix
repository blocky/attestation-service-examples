{
  githubPAT ? builtins.getEnv "GITHUB_PAT",
}:
let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-24.11";
  pkgs = import nixpkgs {
    config = { };
    overlays = [ ];
  };

  bky-as = pkgs.stdenv.mkDerivation rec {
    pname = "bky-as";
    version = "v0.1.0-beta.4";

    src = pkgs.fetchurl {
      url = "https://github.com/blocky/attestation-service-demo/releases/download/v0.1.0-beta.4/bky-as_linux_amd64";
      sha256 = "sha256-Gm/zEiP1jJaloxN6TdiCq112zk6hHlTWu65VkUlNceU=";
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
