#!/bin/bash

# abort on all errors
set -e

APPNAME="SatisfactoryModManager"
export ARCH="x86_64" # Export because linuxdeploy gtk plugin copies i386 libraries too, so linuxdeploy can't decide on architecture
SCRIPT_DIR=$(dirname "$0")
BUILD_DIR="$SCRIPT_DIR/.."

BINARY=$(realpath "$1")
OUTPUT=$2

TMPDIR=$(mktemp -d)
APPDIR="$SCRIPT_DIR/../bin/$APPNAME.AppDir"

# The input binary name is appimage-bin-tmp, so we need to fix that
cp "$BINARY" "$TMPDIR/$APPNAME"
BINARY="$TMPDIR/$APPNAME"

if [ -d "$APPDIR" ]; then
rm -rf "$APPDIR"
fi

mkdir -p "$APPDIR"

# We still copy icons manually instead of using linuxdeploy 
# because the icons are not square, and linuxdeploy checks that
(
cd "$APPDIR" || exit

icons=(16 32 64 128 256 512)
for i in "${icons[@]}"; do
    mkdir -p "usr/share/icons/hicolor/${i}x${i}/apps"
    cp "$BUILD_DIR/icons/${i}x${i}.png" "usr/share/icons/hicolor/${i}x${i}/apps/$APPNAME.png"
done
)

(
cd "$TMPDIR" || exit

wget https://raw.githubusercontent.com/linuxdeploy/linuxdeploy-plugin-gtk/master/linuxdeploy-plugin-gtk.sh
chmod +x linuxdeploy-plugin-gtk.sh

wget https://raw.githubusercontent.com/linuxdeploy/linuxdeploy-plugin-gstreamer/master/linuxdeploy-plugin-gstreamer.sh
chmod +x linuxdeploy-plugin-gstreamer.sh

cp "$SCRIPT_DIR/linuxdeploy-plugin-webkitgtk.sh" .
chmod +x linuxdeploy-plugin-webkitgtk.sh

wget -O linuxdeploy.AppImage https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous/linuxdeploy-${ARCH}.AppImage
chmod +x linuxdeploy.AppImage
)

mkdir -p "$(dirname "$OUTPUT")"

LDAI_OUTPUT="$OUTPUT" DEPLOY_GTK_VERSION="3" "$TMPDIR/linuxdeploy.AppImage" --appimage-extract-and-run \
    --executable "$BINARY" \
    --desktop-file "$SCRIPT_DIR/$APPNAME.desktop" \
    "${ICON_FILES[@]}" \
    --icon-filename "$APPNAME" \
    --custom-apprun "$SCRIPT_DIR/AppRun" \
    --appdir "$APPDIR" \
    --plugin gtk \
    --plugin webkitgtk \
    --plugin gstreamer \
    --output appimage

rm -rf "$TMPDIR"
rm -rf "$APPDIR"
