{
  description = "A Nix flake based development environment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
          overlays = [ ];
        };

      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            nixd

            nodejs_24

            cope # provides make
            mockgen

            go
            gopls
            golangci-lint
          ];
        };
      }
    );
}
