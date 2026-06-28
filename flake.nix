{
  description = "Terminal Git contribution graph for local repositories";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs =
    { self, nixpkgs }:
    let
      systems = [
        "aarch64-darwin"
        "aarch64-linux"
        "x86_64-darwin"
        "x86_64-linux"
      ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = pkgs.buildGoModule {
            pname = "gitviz";
            version = "0.1.0";
            src = ./.;

            vendorHash = "sha256-BPXGllfFhyGjVQXY9g9GzrGtaWuwtx1bJqVzajephHc=";

            meta = {
              description = "Terminal Git contribution graph for local repositories";
              homepage = "https://github.com/anton-fuji/gitviz";
              license = pkgs.lib.licenses.mit;
              mainProgram = "gitviz";
            };
          };
        }
      );

      apps = forAllSystems (system: {
        default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/gitviz";
        };
      });

      devShells = forAllSystems (
        system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = pkgs.mkShell {
            packages = [
              pkgs.go
              pkgs.gotools
            ];
          };
        }
      );
    };
}
