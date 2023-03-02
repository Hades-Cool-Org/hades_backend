package environment

import "hades_backend/app/config"

var Cfg = config.Cfg

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
