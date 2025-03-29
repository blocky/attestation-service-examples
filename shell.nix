let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-24.11";
  pkgs = import nixpkgs {
    config = { };
    overlays = [ ];
  };

  utils =
    { pkgs }:
    let
      system = builtins.split "-" pkgs.stdenv.hostPlatform.system;
      arch = builtins.elemAt (system) 0;
      os = builtins.elemAt (system) 2;
    in
    {
      goos = os;
      goarch =
        if arch == "x86_64" then
          "amd64"
        else if arch == "aarch64" then
          "arm64"
        else
          throw "unknow arch '${arch}', supported arches are 'x86_64' and 'aarch64'";
    };

  bky-as-stable =
    {
      pkgs,
      tag,
    }:
    let
      u = utils { pkgs = pkgs; };
      goos = u.goos;
      goarch = u.goarch;
    in
    pkgs.stdenv.mkDerivation rec {
      pname = "bky-as";
      version = tag;

      src = builtins.fetchurl {
        url = "https://github.com/blocky/attestation-service-demo/releases/download/${version}/bky-as_${goos}_${goarch}";
      };

      unpackPhase = ":";

      installPhase = ''
        install -D -m 555 $src $out/bin/bky-as
      '';
    };

  bky-as-unstable = pkgs.stdenv.mkDerivation rec {
    pname = "bky-as";
    version = "unstable";

    src = ./nix/fetch-bky-as.sh;

    unpackPhase = ":";

    installPhase = ''
      mkdir -p $out/bin
      cp $src $out/bin/fetch-bky-as.sh
      chmod +x $out/bin/fetch-bky-as.sh
    '';
  };

  bky-as-dev-shell =
    {
      pkgs,
      version,
      dev-dependencies,
    }:
    let
      u = utils { pkgs = pkgs; };
      os = u.goos;
      arch = u.goarch;

      unstable-shell = pkgs.mkShellNoCC {
        packages = dev-dependencies ++ [ bky-as-unstable ];

        shellHook = ''
          echo "bky-as hook ran"
          bin=./tmp/bin
          fetch-bky-as.sh $bin ${os} ${arch}
          export PATH=$bin:$PATH
        '';
      };

      bky-as-stable-version = bky-as-stable {
        pkgs = pkgs;
        tag = version;
      };

      stable-shell = pkgs.mkShellNoCC {
        packages = dev-dependencies ++ [ bky-as-stable-version ];
      };
    in
    if version == "unstable" then unstable-shell else stable-shell;
in
bky-as-dev-shell {
  pkgs = pkgs;
  version = "unstable";
  # version = "v0.1.0-beta.5";
  dev-dependencies = [
    pkgs.awscli2 # for setting up the cli
    pkgs.gh # for setting up the cli
    pkgs.git # for project management
    pkgs.gnumake # for project management
    pkgs.go # for building examples
    pkgs.golangci-lint # for linting go files
    pkgs.gotools # for tools like goimports
    pkgs.jq # for examples and setting up the cli
    pkgs.nixfmt-rfc-style # for tools like nix fmt
    pkgs.nodejs_18 # for on chain examples
    pkgs.tinygo # for building wasm
  ];
}
