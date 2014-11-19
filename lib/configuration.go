package lib

type Configuration struct {
	Environment
	Environments map[string]RawEnvironment
}
