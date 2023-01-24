{
    description = "Some tools for bug bounty hunting";

    outputs = {
        self,
        nixpkgs,
        ...
    }: 
    let
        system = "x86_64-linux";
        inherit (nixpkgs) lib;
        supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
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

                filesetup = pkgs.buildGoModule rec {
                    pname = "filesetup";
                    src = ./. + "/go/filesetup";
                    version = "0.0.1";

                    vendorHash = "sha256-nOA56vsDcqiVaF/4ETNt/9rphwBdPt3yFogCJ7ILN3M=";
                };

                recon-engine = pkgs.buildGoModule rec {
                    pname = "recon";
                    src = ./. + "/go/recon-engine";
                    version = "0.0.1";

                    doCheck = false;

                    vendorHash = "sha256-kBbCTC8r5m4F108Ur9WMCQCK9qfVyRU8KcNHvOaxb2Q=";
                };
            }
        );

        apps = forAllSystems(system: 
        let
            pkgs = nixpkgsFor.${system};
            mypkgs = self.packages.${system};
            inherit (mypkgs)
            takesubs
            discordlog
            filesetup
            recon-engine;
        in {
            takesubs = {
                type = "app";
                program = "${takesubs}/bin/takesubs";
            };

            discordlog = {
                type = "app";
                program = "${discordlog}/bin/discordlog";
            };

            filesetup = {
                type = "app";
                program = "${filesetup}/bin/filesetup";
            };

            recon-engine = {
                type = "app";
                program = "${recon-engine}/bin/recon";
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
                    nmap
                ];
            };
        });
    };
}
