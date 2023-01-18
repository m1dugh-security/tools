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
        pkgs = nixpkgs.legacyPackages.${system};
        inherit (nixpkgs) lib;
    in {

        packages.${system} = {
            subdomainTakeover =
            let
                src = "./shell/subdomainTakeover.sh";
            in pkgs.stdenv.mkDerivation {
                name = "takesubs";
                src = ./.; 

                propagetedBuildInputs = with pkgs; [
                    bind
                    curl
                    jq
                ];

                configurePhase = ''
                    mkdir -p $out/bin
                    '';

                installPhase = ''
                    install -D -m 0555 ${src} $out/bin/takesubs
                    '';
            };
        };

        apps.${system} = 
        let
            mypkgs = self.packages.${system};
            inherit (mypkgs) subdomainTakeover;
        in {
            takesubs = {
                type = "app";
                program = "${subdomainTakeover}/bin/takesubs";
            };
        };
    };
}
