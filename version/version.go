package version

import "fmt"

// See http://semver.org/ for more information on Semantic Versioning
var (
	Major      = 0
	Minor      = 6
	Patch      = 9
	PreRelease = "" // Set this via makefile. See docs for more info on release process
)

var Version = fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)

func init() {
	if PreRelease != "" {
		Version += "-" + PreRelease
	}
}
