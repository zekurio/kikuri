package embedded

import (
	"embed"
)

var (

	//go:embed Version.txt
	AppVersion string
	//go:embed Commit.txt
	AppCommit string
	//go:embed migrations
	Migrations embed.FS
	//go:embed webdist
	FrontendFiles embed.FS
)
