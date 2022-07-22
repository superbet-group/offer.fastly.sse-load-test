package main

import (
	"github.com/superbet-group/logging.clients"
	"github.com/superbet-group/trading.commons/pkg/guard"
	"github.com/superbet-group/trading.commons/pkg/tasks"

	"github.com/superbet-group/offer.fastly.sse-load-test/internal/scheduler"
)

func main() {
	defer guard.CapturePanic(log, true, nil)
	log.Info("starting app", logging.Data{
		"environment": conf.AppEnvironment,
		"version":     conf.Version,
	})

	httpClient := buildHTTPClient()

	offerClient := buildOfferClient(httpClient)

	sseClient := buildSSEClient(httpClient)

	sch := scheduler.NewScheduler(log, offerClient, sseClient, int(conf.NumberOfMatches), int(conf.ConnectionsPerMatch))

	tasks.Run(log, tasks.NewSignalHandler(), sch)
}
