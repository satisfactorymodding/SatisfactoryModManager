#! /bin/bash

# abort on all errors
set -e

SCRIPT="$(basename "$(readlink -f "$0")")"

show_usage() {
    echo "Usage: $SCRIPT --appdir <path to AppDir>"
    echo
    echo "Bundles resources for applications that use WebKitGTK into an AppDir"
    echo
    echo "Required variables:"
    echo "  LINUXDEPLOY=\".../linuxdeploy\" path to linuxdeploy (e.g., AppImage); set automatically when plugin is run directly by linuxdeploy"
    echo
    echo "Optional variables:"
    echo "  DEPLOY_WEBKITGTK_VERSION (complete version name of WebKitGTK to deploy, e.g., webkit2gtk-4.1; auto-detect by default)"
}

variable_is_true() {
    local var="$1"

    if [ -n "$var" ] && { [ "$var" == "true" ] || [ "$var" -gt 0 ]; } 2> /dev/null; then
        return 0 # true
    else
        return 1 # false
    fi
}

get_pkgconf_variable() {
    local variable="$1"
    local library="$2"
    local default_value="$3"

    pkgconfig_ret="$("$PKG_CONFIG" --variable="$variable" "$library")"
    if [ -n "$pkgconfig_ret" ]; then
        echo "$pkgconfig_ret"
    elif [ -n "$default_value" ]; then
        echo "$default_value"
    else
        echo "$0: there is no '$variable' variable for '$library' library." > /dev/stderr
        echo "Please check the '$library.pc' file is present in \$PKG_CONFIG_PATH (you may need to install the appropriate -dev/-devel package)." > /dev/stderr
        exit 1
    fi
}

copy_tree() {
    local src=("${@:1:$#-1}")
    local dst="${*:$#}"

    for elem in "${src[@]}"; do
        mkdir -p "${dst::-1}$elem"
        cp "$elem" --archive --parents --target-directory="$dst" $verbose
    done
}

copy_lib_tree() {
    # The source lib directory could be /usr/lib, /usr/lib64, or /usr/lib/x86_64-linux-gnu
    # Therefore, when copying lib directories, we need to transform that target path
    # to a consistent /usr/lib
    local src=("${@:1:$#-1}")
    local dst="${*:$#}"

    for elem in "${src[@]}"; do
        mkdir -p "${dst::-1}${elem/$LD_GTK_LIBRARY_PATH//usr/lib}"
        pushd "$LD_GTK_LIBRARY_PATH"
        cp "$(realpath --relative-to="$LD_GTK_LIBRARY_PATH" "$elem")" --archive --parents --target-directory="$dst/usr/lib" $verbose
        popd
    done
}

get_triplet_path() {
    if command -v dpkg-architecture > /dev/null; then
        echo "/usr/lib/$(dpkg-architecture -qDEB_HOST_MULTIARCH)"
    fi
}



search_library_path() {
    PATH_ARRAY=(
        "$(get_triplet_path)"
        "/usr/lib64"
        "/usr/lib"
    )

    for path in "${PATH_ARRAY[@]}"; do
        if [ -d "$path" ]; then
            echo "$path"
            return 0
        fi
    done
}

search_tool() {
    local tool="$1"
    local directory="$2"

    if command -v "$tool"; then
        return 0
    fi

    PATH_ARRAY=(
        "$(get_triplet_path)/$directory/$tool"
        "/usr/lib64/$directory/$tool"
        "/usr/lib/$directory/$tool"
        "/usr/bin/$tool"
        "/usr/bin/$tool-64"
        "/usr/bin/$tool-32"
    )

    for path in "${PATH_ARRAY[@]}"; do
        if [ -x "$path" ]; then
            echo "$path"
            return 0
        fi
    done
}

find_all() {
    local -n find_array="$1"
    shift
    local -n result_array="$1"
    shift
    
    for (( i=0; i<${#find_array[@]}; i+=2 )); do
        directory=${find_array[i]}
        file=${find_array[i+1]}
        while IFS= read -r -d '' file; do
            result_array+=( "$file" )
        done < <(find "$directory" \( -type l -o -type f \) -name "$file" -print0)
    done
}

DEPLOY_WEBKITGTK_VERSION="${DEPLOY_WEBKITGTK_VERSION:-0}" # When not set by user, this variable use the integer '0' as a sentinel value
APPDIR=""

while [ "$1" != "" ]; do
    case "$1" in
        --plugin-api-version)
            echo "0"
            exit 0
            ;;
        --appdir)
            APPDIR="$2"
            shift
            shift
            ;;
        *)
            echo "Invalid argument: $1"
            echo
            show_usage
            exit 2
    esac
done

if [[ "$APPDIR" == "" ]]; then
    show_usage
    exit 2
fi

APPDIR="$(realpath "$APPDIR")"
mkdir -p "$APPDIR"

if command -v pkgconf > /dev/null; then
    PKG_CONFIG="pkgconf"
elif command -v pkg-config > /dev/null; then
    PKG_CONFIG="pkg-config"
else
    echo "$0: pkg-config/pkgconf not found in PATH, aborting"
    exit 1
fi

if ! command -v find &>/dev/null && ! type find &>/dev/null; then
    echo -e "$0: find not found.\nInstall findutils then re-run the plugin."
    exit 1
fi

if [ -z "$LINUXDEPLOY" ]; then
    echo -e "$0: LINUXDEPLOY environment variable is not set.\nDownload a suitable linuxdeploy AppImage, set the environment variable and re-run the plugin."
    exit 1
fi

webkitgtk_versions=0 # Count versions of WebKitGTK when auto-detecting WebKitGTK version
if [ "$DEPLOY_WEBKITGTK_VERSION" = "0" ]; then
    echo "Determining which WebKitGTK version to deploy"
    while IFS= read -r -d '' file; do
        if [ "$DEPLOY_WEBKITGTK_VERSION" != "webkit2gtk-4.0" ] && ldd "$file" | grep -q "libwebkit2gtk-4.0.so"; then
            DEPLOY_WEBKITGTK_VERSION="webkit2gtk-4.0"
            webkitgtk_versions="$((webkitgtk_versions+1))"
        fi
        if [ "$DEPLOY_WEBKITGTK_VERSION" != "webkit2gtk-4.1" ] && ldd "$file" | grep -q "libwebkit2gtk-4.1.so"; then
            DEPLOY_WEBKITGTK_VERSION="webkit2gtk-4.1"
            webkitgtk_versions="$((webkitgtk_versions+1))"
        fi
        if [ "$DEPLOY_WEBKITGTK_VERSION" != "webkitgtk-6.0" ] && ldd "$file" | grep -q "libwebkitgtk-6.0.so"; then
            DEPLOY_WEBKITGTK_VERSION="webkitgtk-6.0"
            webkitgtk_versions="$((webkitgtk_versions+1))"
        fi
    done < <(find "$APPDIR/usr/bin" -executable -type f -print0)
fi

if [ "$webkitgtk_versions" -gt 1 ]; then
    echo "$0: can not deploy multiple WebKitGTK versions at the same time."
    echo "Please set DEPLOY_WEBKITGTK_VERSION to {webkit2gtk-4.0, webkit2gtk-4.1, webkitgtk-6.0}."
    exit 1
elif [ "$DEPLOY_WEBKITGTK_VERSION" = "0" ]; then
    echo "$0: failed to auto-detect WebKitGTK version."
    echo "Please set DEPLOY_WEBKITGTK_VERSION to {webkit2gtk-4.0, webkit2gtk-4.1, webkitgtk-6.0}."
    exit 1
fi



echo "Installing AppRun hook"
HOOKSDIR="$APPDIR/apprun-hooks"
HOOKFILE="$HOOKSDIR/linuxdeploy-plugin-webkitgtk.sh"
mkdir -p "$HOOKSDIR"
cat > "$HOOKFILE" <<EOF
#! /usr/bin/env bash

EOF


echo "Deploying extra WebKitGTK executables and libraries for $DEPLOY_WEBKITGTK_VERSION"

LD_GTK_LIBRARY_PATH="$(realpath "${LD_GTK_LIBRARY_PATH:-$(search_library_path)}")"
gio_libdir="$(get_pkgconf_variable "libdir" "gio-2.0" "$LD_GTK_LIBRARY_PATH")"
webkit2gtk_libexecdir="$(get_pkgconf_variable "libdir" "$DEPLOY_WEBKITGTK_VERSION" "/usr/lib")/$DEPLOY_WEBKITGTK_VERSION"

WEBKIT_FIND_ARRAY=(
    "$webkit2gtk_libexecdir" "WebKitNetworkProcess"
    "$webkit2gtk_libexecdir" "WebKitWebProcess"
    "$webkit2gtk_libexecdir" "lib${DEPLOY_WEBKITGTK_VERSION%-*}injectedbundle.so"
)
WEBKIT_LIBRARIES=()
find_all WEBKIT_FIND_ARRAY WEBKIT_LIBRARIES


GIO_FIND_ARRAY=(
    "$gio_libdir"            "libgiognutls.so*"
)
GIO_MODULES=()
find_all GIO_FIND_ARRAY GIO_MODULES

LIBRARIES=(
    "${WEBKIT_LIBRARIES[@]}"
    "${GIO_MODULES[@]}"
)

LIBRARIES_ARGS=()
for lib in "${LIBRARIES[@]}"; do
    LIBRARIES_ARGS+=( "--library=$lib" )
done

env LINUXDEPLOY_PLUGIN_MODE=1 "$LINUXDEPLOY" --appdir="$APPDIR" "${LIBRARIES_ARGS[@]}"



echo "Creating symlinks for hardcoded WebKitGTK executable paths"

# The location to these files is hardcoded in WebKitGTK, so we need to create symlinks from the expected location to the actual location of these files in the AppDir
for lib in "${LIBRARIES[@]}"; do
    lib_basename="$(basename "$lib")"
    lib_dir="$(dirname "$APPDIR/$lib")"
    mkdir -p "$lib_dir"
    relative_path="$(realpath --relative-to="$lib_dir" "$APPDIR/usr/lib/$lib_basename")"
    ln -sf "$relative_path" "$APPDIR/$lib"
done



echo "Patching absolute paths in WebKitGTK libraries"

# binary patch absolute paths in libwebkit files (from tauri's linuxdeploy-plugin-gtk)
find "$APPDIR"/usr/lib* -name 'libwebkit*' -exec sed -i -e "s|/usr|././|g" '{}' \;



echo "Setting GIO_EXTRA_MODULES"

GIO_EXTRA_MODULES=()
for lib in "${GIO_MODULES[@]}"; do
    dir="$(dirname "$lib")"
    entry="\$APPDIR$dir"
    if [[ ! " ${GIO_EXTRA_MODULES[*]} " =~ " ${entry} " ]]; then
        GIO_EXTRA_MODULES+=( "$entry" )
    fi
done

GIO_EXTRA_MODULES_STR="$(IFS=:; echo "${GIO_EXTRA_MODULES[*]}")"

cat >> "$HOOKFILE" <<EOF
export GIO_EXTRA_MODULES="$GIO_EXTRA_MODULES_STR"
EOF
