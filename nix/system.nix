{ pkgs }:
let
  system = pkgs.lib.strings.splitString "-" pkgs.stdenv.hostPlatform.system;
  arch = builtins.elemAt (system) 0;
  os = builtins.elemAt (system) 1;
in
{
  goos = os;
  goarch =
    if arch == "x86_64" then
      "amd64"
    else if arch == "aarch64" then
      "arm64"
    else
      throw "unknown arch '${arch}', supported arches are 'x86_64' and 'aarch64'";
}
