package whisky

type whiskyPlist struct {
	DefaultBottleLocation string `plist:"defaultBottleLocation"`
}

type bottleVMPlist struct {
	Paths []urlPlist `plist:"paths"`
}

type urlPlist struct {
	Relative string `plist:"relative"`
}
