package shell

type Shell interface {
	Name() string
	Init() string
	Uninit() string
}
