package configuration

// Configuration is the rambler configuration type as loaded from the configuration
// file and extended by the command-line.
type Configuration struct {
	Driver       string
	Protocol     string
	Host         string
	Port         uint64
	User         string
	Password     string
	Database     string
	Directory    string
	Environments map[string]RawEnvironment
}
