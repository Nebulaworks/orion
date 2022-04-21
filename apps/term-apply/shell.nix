# TODO: Pin this to a specific version of nixpkgs
{ pkgs ? import <nixpkgs> {} }:

with pkgs;

mkShell {

  buildInputs = [
    go_1_17
    awscli2
  ];

  shellHook = ''
  '';
}
