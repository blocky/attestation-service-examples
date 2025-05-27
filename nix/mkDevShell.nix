{
  pkgs,
  version,
  devDependencies,
}:
let
  system = import ./system.nix { pkgs = pkgs; };
  goos = system.goos;
  goarch = system.goarch;

  isCommit = x: builtins.match "^[0-9a-f]{40}$" x != null;

  bky-as-stable = pkgs.stdenv.mkDerivation {
    pname = "bky-as";
    version = version;
    src = builtins.fetchurl {
      url = "https://github.com/blocky/attestation-service-cli/releases/download/${version}/bky-as_${goos}_${goarch}";
    };
    unpackPhase = ":";
    installPhase = ''
      install -D -m 555 $src $out/bin/bky-as
    '';
  };

  stableShell = pkgs.mkShellNoCC {
    packages = devDependencies ++ [ bky-as-stable ];
    shellHook = ''
      set -e
      export AS_VERSION=${version}

      render-md() {
        for file in $(find . -type f -name '*.md'); do
          mo "$file" > "$file.tmp" && mv "$file.tmp" "$file"
        done
      }
      echo "Stable bky-as version: $AS_VERSION"
      set +e
    '';
  };

  bky-as-unstable = pkgs.stdenv.mkDerivation {
    pname = "bky-as";
    version = version;
    src = ./fetch-bky-as.sh;
    unpackPhase = ":";
    installPhase = ''
      install -D -m 555 $src $out/bin/fetch-bky-as.sh
    '';
  };

  unstableShell = pkgs.mkShellNoCC {
    packages =
      devDependencies
      ++ [ bky-as-unstable ]
      #dependencies required by the fetch script
      ++ [
        pkgs.gh
        pkgs.awscli2
        pkgs.jq
      ];
    shellHook = ''
      set -e

      bin=$(pwd)/tmp/bin
      fetch-bky-as.sh $bin ${version} ${goos} ${goarch}
      export PATH=$bin:$PATH
      export AS_VERSION=${version};

      echo "Unstable bky-as version: $AS_VERSION"
      set +e
    '';
  };
in
if isCommit version || version == "latest" then
  unstableShell
else
  # If the version is not a commit hash or "latest", we assume it is a stable
  # release version.
  stableShell
