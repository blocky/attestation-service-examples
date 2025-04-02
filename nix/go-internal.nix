{
  pkgs,
}:
pkgs.buildGoModule rec {
  pname = "go-internal";
  version = "1.14.1";

  src = pkgs.fetchFromGitHub {
    owner = "rogpeppe";
    repo = "go-internal";
    rev = "v${version}";
    hash = "sha256-6NzhXCCD1Qhj05WHbCDxH5hwfNM6psoAk7uIxm7N55E=";
  };

  vendorHash = "sha256-WoRmZbYYpwVVetlxJDjUu9jGgwLXUD3/PnUF6ksUT70=";

  doCheck = false;

  meta = {
    description = "The testscript command runs testscript scripts in a fresh temporary work directory tree";
    homepage = "https://github.com/rogpeppe/go-internal/cmd/testscript";
    license = pkgs.lib.licenses.bsd3;
    maintainers = with pkgs.lib.maintainers; [ dlm ];
  };
}
