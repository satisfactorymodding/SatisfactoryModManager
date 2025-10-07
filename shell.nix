let
  # This is the commit smm was added to nixpkgs, somewhen after this, something broke packaging.
  # TODO: fix the build failure on more recent versions and then update the hash to the latest tag
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/a7808b41c3cfd040fbe03e12a68654db3c02d5bc";
  pkgs = import nixpkgs { config = {}; overlays = []; };

in
	pkgs.mkShellNoCC {
		nativeBuildInputs = 
		let
			wails' = pkgs.wails.override { nodejs = pkgs.nodejs_20; };
		in
		with pkgs.buildPackages; [
			wails'
			pnpm_8
			vite
			wrapGAppsHook3
		];
		buildInputs = with pkgs.buildPackages; [
			glib-networking
		];

		shellHook = ''
		export GIO_MODULE_DIR=${pkgs.glib-networking}/lib/gio/modules/
		'';
	}
