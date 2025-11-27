{
  description = "Subscription service for Effective Mobile";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;

          config = {
            allowUnfree = false;
          };
        };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.yazi
            pkgs.lynx
            pkgs.tmux
            pkgs.htop

            pkgs.git
            pkgs.curl
            pkgs.wget
            pkgs.jq
            pkgs.yq-go

            pkgs.ripgrep
            pkgs.fd
            pkgs.gat
            pkgs.fzf

            pkgs.disko

            pkgs.go
            pkgs.golangci-lint
            pkgs.air
            pkgs.gnumake
            pkgs.docker-compose
            pkgs.git

            pkgs.postgresql

            pkgs.oapi-codegen
            pkgs.go-swag
            pkgs.go-swagger
          ];
        };
      }
    );
}
