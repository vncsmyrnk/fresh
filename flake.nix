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
            ./internal
            ./go.mod
          ];
        };
        version = "0.1.0";
        vendorHash = null;

        ldflags = [
          "-s -w -X github.com/vncsmyrnk/fresh/cmd.version=${version}"
        ];
      };

      devShell = pkgs.mkShell {
        packages = with pkgs; [
          go
          delve
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
    in
    {
      packages.${system}.default = fresh;
      devShells.${system}.default = devShell;
      apps.${system} = {
        lint = {
          type = "app";
          program = "${pkgs.lib.getExe lintScript}";
        };
        fresh = {
          type = "app";
          program = "${pkgs.lib.getExe runScript}";
        };
      };
    };
}
