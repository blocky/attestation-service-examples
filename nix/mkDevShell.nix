{
  pkgs,
  asVersion,
  cVersion,
  devDependencies,
}:
let
  system = import ./system.nix { pkgs = pkgs; };
  goos = system.goos;
  goarch = system.goarch;

  isCommit = x: builtins.match "^[0-9a-f]{40}$" x != null;

  bky-as-stable = pkgs.stdenv.mkDerivation {
    pname = "bky-as";
    version = asVersion;
    src = builtins.fetchurl {
      url = "https://github.com/blocky/attestation-service-cli/releases/download/${asVersion}/bky-as_${goos}_${goarch}";
    };
    unpackPhase = ":";
    installPhase = ''
      install -D -m 555 $src $out/bin/bky-as
    '';
  };

  bky-c-stable = pkgs.stdenv.mkDerivation {
    pname = "bky-c";
    version = cVersion;
    src = builtins.fetchurl {
      # todo: switch bky-as to bky-c
      # url = "https://github.com/blocky/compiler/releases/download/${cVersion}/bky-c_v${cVersion}_${goos}_${goarch}";
      # todo: until bky-c is released, download bky-as binary instead (to test the nix setup)
      url = "https://github.com/blocky/attestation-service-cli/releases/download/${cVersion}/bky-as_${goos}_${goarch}";
    };
    unpackPhase = ":";
    installPhase = ''
      install -D -m 555 $src $out/bin/bky-c
    '';
  };

  stableShell = pkgs.mkShellNoCC {
    packages = devDependencies ++ [ bky-as-stable bky-c-stable ];
    shellHook = ''
      set -e
      export AS_VERSION=${asVersion}
      export C_VERSION=${cVersion}

      render-md() {
        for file in $(git ls-files '*.md'); do
          echo "Processing $file"
          mo --open="{{{"  --close="}}}" "$file" > "$file.tmp" && mv "$file.tmp" "$file"
        done
      }

      upgrade-basm() {
        find . -type f -name go.mod -execdir bash -c 'pwd && go get -u github.com/blocky/basm-go-sdk && go mod tidy' \;
      }

      echo "Stable bky-as version: $AS_VERSION"
      echo "Stable bky-c version: $C_VERSION"
      set +e
    '';
  };

  bky-as-unstable = pkgs.stdenv.mkDerivation {
    pname = "bky-as";
    asVersion = asVersion;
    src = ./fetch-bky-as.sh;
    unpackPhase = ":";
    installPhase = ''
      install -D -m 555 $src $out/bin/fetch-bky-as.sh
    '';
  };

  unstableShell = pkgs.mkShellNoCC {
    packages =
      devDependencies
      ++ [ bky-as-unstable bky-c-stable ]
      #dependencies required by the fetch script
      ++ [
        pkgs.gh
        pkgs.awscli2
        pkgs.jq
      ];
    shellHook = ''
      set -e

      bin=$(pwd)/tmp/bin
      fetch-bky-as.sh $bin ${asVersion} ${goos} ${goarch}
      export PATH=$bin:$PATH
      export AS_VERSION=${asVersion};

      echo "Unstable bky-as version: $AS_VERSION"
      echo "Unstable bky-c version: $C_VERSION"
      set +e
    '';
  };
in
if isCommit asVersion || asVersion == "latest" then
  unstableShell
else
  # If the asVersion is not a commit hash or "latest", we assume it is a stable
  # release asVersion.
  stableShell
