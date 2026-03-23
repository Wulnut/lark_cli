package cmddeps

import (
	"io"
	"lark_cli/internal/auth"
	"lark_cli/internal/config"
	"lark_cli/internal/session"
)

// Deps holds dependencies injected into commands.
type Deps struct {
	Config config.Config
	Store  session.Store

	PluginTokenProvider auth.PluginTokenProvider
	HeaderProvider      auth.HeaderProvider

	Stdout io.Writer
	Stderr io.Writer
}
