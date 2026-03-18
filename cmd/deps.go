package cmd

import "io"

// Deps holds dependencies injected into commands.
// Concrete types will be added as internal packages are created.
type Deps struct {
	SessionPath string
	Stdout      io.Writer
	Stderr      io.Writer
}
