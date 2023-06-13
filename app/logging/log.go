package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// Initialize creates a new JSON zap logger.
func Initialize() *zap.Logger {
	logEnc := zap.NewProductionEncoderConfig()
	logEnc.EncodeTime = zapcore.ISO8601TimeEncoder
	logEnc.TimeKey = "timestamp"

	log := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(logEnc),
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevel(),
	))

	log = log.With(
		zap.String("app_name", "hades_backend"),
		zap.Time("start_time", time.Now().Local()),
	)

	return log
}
