package version

import "fmt"

// See http://semver.org/ for more information on Semantic Versioning
var (
	Major      = 1
	Minor      = 1
	Patch      = 0
	PreRelease = "" // Set this via makefile. See docs for more info on release process
)

var Version = fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)

func init() {
	if PreRelease != "" {
		Version += "-" + PreRelease
	}
}
