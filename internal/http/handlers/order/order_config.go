package order

type Options struct {
	ErrorDir string
	UseINI   bool // true = ignora ip/local_db vindos do app
}

var opts Options

func Configure(o Options) {
	opts = o
}
