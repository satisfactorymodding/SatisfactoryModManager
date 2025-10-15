let
  # This is the commit smm was added to nixpkgs, somewhen after this, something broke packaging.
  # TODO: fix the build failure on more recent versions and then update the hash to the latest tag
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/cff40eaf67b2a24b92296894608bcff305c675fc";
  pkgs = import nixpkgs { config = {}; overlays = []; };

in
	pkgs.mkShellNoCC {
		nativeBuildInputs = with pkgs.buildPackages; [
			wails
			pnpm
			vite
		];
		buildInputs = with pkgs.buildPackages; [
			glib-networking
		];
	}
