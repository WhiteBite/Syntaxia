package version

// Version information - set via ldflags during build
// go build -ldflags "-X syntaxia/infrastructure/version.Version=v1.0.0 -X syntaxia/infrastructure/version.GitCommit=abc123"
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// Info returns version information
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`
}

// GetInfo returns current version info
func GetInfo() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
	}
}
