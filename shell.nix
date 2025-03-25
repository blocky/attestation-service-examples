{
  githubPAT ? builtins.getEnv "GITHUB_PAT",
}:
let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-24.11";
  pkgs = import nixpkgs {
    config = { };
    overlays = [ ];
  };

  bky-as = pkgs.buildGoModule rec {
    pname = "bky-as";
    version = "v0.1.0-beta.4";

    tmpDir = "/tmp";

    env = {
      GOPRIVATE = "github.com/blocky/*";
      HOME = tmpDir;
    };

    src = builtins.fetchGit {
      rev = "f6a2c0f965cfca9583ab29e10b9e2e5acf006046";
      url = "git@github.com:blocky/delphi.git";
    };

    doCheck = false;

    preBuild = ''
      echo machine github.com login doesNotMatter password ${githubPAT} > ${tmpDir}/.netrc
    '';

    postBuild = ''
      alias bky-as=cli
    '';

    vendorHash = "sha256-GXlZz3L5vd1v9NHlaagKw6aY3LEyt9E10reh6EvZ4Bw=";

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
