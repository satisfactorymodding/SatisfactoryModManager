#!/usr/bin/env bash

# Wails uses does not convert slashes to backslashes in fileAssociations.iconName, as it does everywhere else
# So we keep the backslashes in the JSON for development and convert them to slashes here
# Writing to a temporary file and then moving it is necessary because jq does not support in-place editing
jq '.info.fileAssociations[].iconName |= gsub("\\\\"; "/")' wails.json > wails.json.tmp && mv wails.json.tmp wails.json