package monitoring

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/integrations/nrzap"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"os"
)

var (
	NewRelicApp *newrelic.Application
)

func New(l *zap.Logger) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("hades"),
		newrelic.ConfigLicense("46e40bb3b9ffca244d9ff0edee46bdb4ae8cNRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
		nrzap.ConfigLogger(l.Named("hades")),
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	NewRelicApp = app
}
