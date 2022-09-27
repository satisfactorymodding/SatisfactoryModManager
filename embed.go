//go:build !bindings

package main

import "embed"

//go:embed all:frontend/build
var assets embed.FS
