{
  description = "A simple shell implementation built for educational purposes.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};

      fresh = pkgs.buildGoModule rec {
        name = "fresh";
        src = pkgs.lib.fileset.toSource {
          root = ./.;
          fileset = pkgs.lib.fileset.unions [
            ./cmd/fresh/main.go
            ./internal
            ./go.mod
            ./go.sum
          ];
        };
        version = "0.1.0";
        vendorHash = "sha256-fxxp7ECuMMmLw7L5/lPi7lfmJM2p+Te9AqXm5Xed3G0=";

        ldflags = [
          "-s -w -X github.com/vncsmyrnk/fresh/cmd.version=${version}"
        ];
      };

      lintScript = pkgs.writeShellApplication {
        name = "run-linters";
        runtimeInputs = with pkgs; [
          golangci-lint
        ];

        text = ''
          golangci-lint run
        '';
      };

      runScript = pkgs.writeShellApplication {
        name = "run-fresh";
        runtimeInputs = with pkgs; [
          go
        ];

        text = ''
          go run ./cmd/fresh/main.go
        '';
      };

      debugScript = pkgs.writeShellApplication {
        name = "run-dlv";
        runtimeInputs = with pkgs; [
          go
          delve
        ];

        text = ''
          dlv debug --headless --listen=:2345 --api-version=2 ./cmd/fresh/main.go
        '';
      };
    in
    {
      packages.${system}.default = fresh;
      apps.${system} = {
        lint = {
          type = "app";
          program = "${pkgs.lib.getExe lintScript}";
        };
        debug = {
          type = "app";
          program = "${pkgs.lib.getExe debugScript}";
        };
        fresh = {
          type = "app";
          program = "${pkgs.lib.getExe runScript}";
        };
      };
    };
}
