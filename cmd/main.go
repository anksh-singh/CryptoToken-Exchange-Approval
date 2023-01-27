package main

import (
	"bridge-allowance/config"
	"bridge-allowance/internal/boot"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	conf := config.LoadConfig("", "")
	tracer.Start(
		tracer.WithEnv(conf.Datadog.Env),
		tracer.WithService(conf.DATADOG_SERVICE),
		tracer.WithServiceVersion(conf.Datadog.Version),
	)
	// When the tracer is stopped, it will flush everything it has to the Datadog Agent before quitting.
	defer tracer.Stop()
	// cobra root command execute
	boot.Execute()
}
