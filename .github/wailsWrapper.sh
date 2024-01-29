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

wails build "${ARGS[@]}"

if [ -n "$OUTPUT_FULL_PATH" ]; then
  cp "build/bin/$OUTPUT_FILENAME" "$OUTPUT_FULL_PATH"
fi