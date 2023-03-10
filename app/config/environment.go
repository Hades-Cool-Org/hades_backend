package config

var cfg = Cfg

const (
	EnvProd = "prod"
)

func IsProd() bool {
	return cfg.Env == EnvProd
}

func ExecIfProd(fn func()) {
	if IsProd() {
		fn()
	}
}
