package config

const (
	EnvProd = "prod"
)

func IsProd() bool {
	return Cfg.Env == EnvProd
}

func ExecIfProd(fn func()) {
	if IsProd() {
		fn()
	}
}
