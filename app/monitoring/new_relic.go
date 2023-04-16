package monitoring

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"os"
)

var (
	NewRelicApp *newrelic.Application
)

func init() {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("hades"),
		newrelic.ConfigLicense("46e40bb3b9ffca244d9ff0edee46bdb4ae8cNRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	NewRelicApp = app
}
