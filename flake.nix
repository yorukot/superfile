{
  description = "A fancy, pretty terminal file manager";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-utils.url = "github:numtide/flake-utils";

    flake-compat.url = "github:edolstra/flake-compat";
    flake-compat.flake = false;

    gomod2nix.url = "github:nix-community/gomod2nix";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
    gomod2nix.inputs.flake-utils.follows = "flake-utils";
  };

  outputs = inputs @ {...}:
    inputs.flake-utils.lib.eachDefaultSystem
    (
      system: let
        overlays = [
          inputs.gomod2nix.overlays.default
        ];
        pkgs = import inputs.nixpkgs {
          inherit system overlays;
        };
      in rec {
        packages = rec {
          superfile = pkgs.buildGoApplication {
            pname = "superfile";
            version = "1.1.5";
            src = ./.;
            modules = ./gomod2nix.toml;
          };
          default = superfile;
        };

        apps = rec {
          superfile = {
            type = "app";
            program = "${packages.superfile}/bin/superfile";
          };
          default = superfile;
        };

        devShells = {
          default = pkgs.mkShell {
            packages = with pkgs; [
              ## golang
              delve
              go-outline
              go
              golangci-lint
              gopkgs
              gopls
              gotools
              nix
              gomod2nix
              nixpkgs-fmt
            ];
          };
        };
      }
    );
}
