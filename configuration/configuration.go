package configuration

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
