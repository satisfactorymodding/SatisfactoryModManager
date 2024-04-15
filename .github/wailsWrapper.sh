#!/usr/bin/env bash

# This script is a wrapper for wails such that it can be used as the "gobinary" of goreleaser build tasks
# Wails only takes the filename as the -o argument, placing the binary in build/bin,
# while goreleaser provides a full path
#
# So this script will extract the filename, pass that to wails along with the rest of the arguments,
# then copy the binary to the location expected by goreleaser, creating the directories if necessary

ARGS=()

while [[ $# -gt 0 ]]; do
  case $1 in
    -o)
      OUTPUT_FULL_PATH="$2"
      shift # past argument
      shift # past value
      ;;
    .)
      # wails does not like the module path as an argument
      # it ignores everything after the first non-flag argument
      shift # past argument
      ;;
    *)
      ARGS+=("$1") # save arg
      shift # past argument
      ;;
  esac
done

if [ -n "$OUTPUT_FULL_PATH" ]; then
  OUTPUT_DIR=$(dirname "$OUTPUT_FULL_PATH")
  mkdir -p "$OUTPUT_DIR"
  OUTPUT_FILENAME=$(basename "$OUTPUT_FULL_PATH")
  ARGS+=("-o")
  ARGS+=("$OUTPUT_FILENAME")
fi

wails "${ARGS[@]}"

if [ -n "$OUTPUT_FULL_PATH" ]; then
  if [ "$GOOS" == "darwin" ]; then
    pushd "build/bin" || exit 1
    zip -r "$OUTPUT_FILENAME.zip" "$OUTPUT_FILENAME.app" # zip would add the .zip extension anyway
    popd || exit 1
    OUTPUT_FILENAME="$OUTPUT_FILENAME.zip"
  fi
  cp -r "build/bin/$OUTPUT_FILENAME" "$OUTPUT_FULL_PATH"
fi