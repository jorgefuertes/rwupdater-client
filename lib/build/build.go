package build

import (
	"fmt"
	"runtime"
)

var version string = "undefined"
var user string = "undefined"
var time string = "undefined"
var number string = "undefined"

// Version - complete version string
func Version() string {
	return fmt.Sprintf(`
 ______         __               ________ __ __     __
|   __ \.-----.|  |_.----.-----.|  |  |  |__|  |--.|__|
|      <|  -__||   _|   _|  _  ||  |  |  |  |    < |  |
|___|__||_____||____|__| |_____||________|__|__|__||__|

FPGA TEAM UPDATER %s %s/%s

`, version, runtime.GOOS, runtime.GOARCH)
}

// VersionShort - short version string
func VersionShort() string {
	return version
}

// BinVersion - Bin version for AT
func BinVersion() string {
	return fmt.Sprintf("Bin version (%s %s): %s", user, number, version)
}

// CompileTime - Compile time string
func CompileTime() string {
	return fmt.Sprintf("compile time:%s", time)
}
