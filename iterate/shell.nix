let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-24.05";
  pkgs = import nixpkgs {
    config = { };
    overlays = [ ];
  };
in
pkgs.mkShellNoCC {
  shellHook = ''
    ./bky-as inspect || ln -s $HOME/repos/blocky/delphi/dist/bky-as-demo/bky-as_linux_amd64 bky-as
  '';

  packages = with pkgs; [
    tinygo
    jq
    gnumake
  ];
}
