#!/bin/bash

APPNAME="SatisfactoryModManager"
export ARCH="x86_64" # Export because linuxdeploy gtk plugin copies i386 libraries too, so linuxdeploy can't decide on architecture
SCRIPT_DIR=$(dirname "$0")
BUILD_DIR="$SCRIPT_DIR/.."

BINARY=$(realpath "$1")
OUTPUT=$2

TMPDIR=$(mktemp -d)
APPDIR="$SCRIPT_DIR/../bin/$APPNAME.AppDir"

if [ -d "$APPDIR" ]; then
rm -rf "$APPDIR"
fi

mkdir -p "$APPDIR"

mkdir -p "$APPDIR/usr/bin"
mkdir -p "$APPDIR/usr/lib"
mkdir -p "$APPDIR/usr/lib64"

(
cd "$APPDIR" || exit
cp "$BINARY" "usr/bin/$APPNAME"
cp "$BUILD_DIR/appicon.png" "$APPNAME.png"
cp "$BUILD_DIR/appicon.png" ".DirIcon"

icons=(16 32 64 128 256 512)
for i in "${icons[@]}"; do
    mkdir -p "usr/share/icons/hicolor/${i}x${i}/apps"
    cp "$BUILD_DIR/icons/${i}x${i}.png" "usr/share/icons/hicolor/${i}x${i}/apps/$APPNAME.png"
done

mkdir -p "usr/share/applications"
cp "$SCRIPT_DIR/$APPNAME.desktop" "usr/share/applications/$APPNAME.desktop"
ln -sf "usr/share/applications/$APPNAME.desktop" "$APPNAME.desktop"
)

(
cd "$APPDIR" || exit

# Copy webkit2gtk libraries
find -L /usr/lib* -name WebKitNetworkProcess -exec mkdir -p "$(dirname '{}')" \; -exec cp --parents '{}' "." \; || true
find -L /usr/lib* -name WebKitWebProcess -exec mkdir -p "$(dirname '{}')" \; -exec cp --parents '{}' "." \; || true
find -L /usr/lib* -name libwebkit2gtkinjectedbundle.so -exec mkdir -p "$(dirname '{}')" \; -exec cp --parents '{}' "." \; || true

# Download AppRun
wget -O AppRun https://github.com/AppImage/AppImageKit/releases/download/continuous/AppRun-${ARCH}
chmod +x AppRun
)

(
cd "$TMPDIR" || exit

wget https://raw.githubusercontent.com/tauri-apps/linuxdeploy-plugin-gtk/master/linuxdeploy-plugin-gtk.sh
chmod +x linuxdeploy-plugin-gtk.sh

wget https://raw.githubusercontent.com/tauri-apps/linuxdeploy-plugin-gstreamer/master/linuxdeploy-plugin-gstreamer.sh
chmod +x linuxdeploy-plugin-gstreamer.sh

wget -O linuxdeploy.AppImage https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous/linuxdeploy-${ARCH}.AppImage
chmod +x linuxdeploy.AppImage
)

mkdir -p "$(dirname "$OUTPUT")"

LDAI_OUTPUT="$OUTPUT" DEPLOY_GTK_VERSION="3" "$TMPDIR/linuxdeploy.AppImage" --appimage-extract-and-run --appdir "$APPDIR" --plugin gtk --plugin gstreamer --output appimage

rm -rf "$TMPDIR"
rm -rf "$APPDIR"
