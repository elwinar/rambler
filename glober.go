package rambler

type Glober interface {
	Glob(string) ([]string, error)
}
