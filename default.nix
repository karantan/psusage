{ buildGoModule
, nix-gitignore
}:

let
  pkgs = import (builtins.fetchTarball {
    url = "https://github.com/NixOS/nixpkgs/archive/refs/tags/22.05.tar.gz";
  }) { };
in buildGoModule.override { go = pkgs.go_1_18; }  rec {
  pname = "psusage";
  version = "0.1.1";

  src = nix-gitignore.gitignoreSource [ ] ./.;

  # The checksum of the Go module dependencies. `vendorSha256` will change if go.mod changes.
  # If you don't know the hash, the first time, set:
  # sha256 = "0000000000000000000000000000000000000000000000000000";
  # then nix will fail the build with such an error message:
  # hash mismatch in fixed-output derivation '/nix/store/m1ga09c0z1a6n7rj8ky3s31dpgalsn0n-source':
  # wanted: sha256:0000000000000000000000000000000000000000000000000000
  # got:    sha256:173gxk0ymiw94glyjzjizp8bv8g72gwkjhacigd1an09jshdrjb4
  vendorSha256 = "1slzmfl36pmkn2a681874mkx2c2hkyalpqaz89lrbv4xkmb71sx6";
  ldflags = "-X cdp/version.Version=${version}";
}
