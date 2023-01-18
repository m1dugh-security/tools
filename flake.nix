{
    description = "Some tools for bug bounty hunting";

    inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    };

    outputs = {
        self,
        nixpkgs
    }: 
    let
        system = "x86_64-linux";
        inherit (nixpkgs) lib;
        supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-linux" "aarch64-linux" ];
        forAllSystems = lib.genAttrs supportedSystems;
        nixpkgsFor = forAllSystems(system: import nixpkgs { inherit system; });
    in {

        packages = forAllSystems(system:
            let
                pkgs = nixpkgsFor.${system};
            in {
                takesubs =
                let
                    src = "./shell/subdomainTakeover.sh";
                in pkgs.stdenv.mkDerivation {
                    name = "takesubs";
                    src = ./.; 

                    propagetedBuildInputs = with pkgs; [
                        bind
                    ];

                    configurePhase = ''
                        mkdir -p $out/bin
                        '';

                    installPhase = ''
                        install -D -m 0555 ${src} $out/bin/takesubs
                        '';
                };

                discordlog = pkgs.buildGoModule rec {
                    pname = "discordlog";
                    src = ./. + "/go/discordlog";
                    version = "0.0.1";

                    vendorHash = "sha256-Pdz3EpIZSxTHhf5tZ34iZnVRp2lKTWh61QAvRyrTLJg=";
                };
            }
        );

        apps = forAllSystems(system: 
        let
            pkgs = nixpkgsFor.${system};
            mypkgs = self.packages.${system};
            inherit (mypkgs)
            takesubs
            discordlog;
        in {
            takesubs = {
                type = "app";
                program = "${takesubs}/bin/takesubs";
            };

            discordlog = {
                type = "app";
                program = "${discordlog}/bin/discordlog";
            };
        });

        devShells = forAllSystems(system: 
        let
            pkgs = nixpkgsFor.${system};
        in {
            go = pkgs.mkShell {
                nativeBuildInputs = with pkgs; [
                    gnumake
                    go
                ];
            };
        });
    };
}
