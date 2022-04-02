package version

import "fmt"

var (
	Commit    string
	BuildTime string
)

func BuildInfo() string {
	return fmt.Sprintf(`
Build Info:
===========
commit:      %s
build time:  %s
	`, Commit, BuildTime)
}
