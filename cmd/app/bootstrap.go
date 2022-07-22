package main

import (
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/superbet-group/logging.clients"
	"github.com/superbet-group/logging.clients/zap"
	"go.uber.org/zap/zapcore"

	"github.com/superbet-group/offer.fastly.sse-load-test/cmd/app/config"
	"github.com/superbet-group/offer.fastly.sse-load-test/internal/offer"
	"github.com/superbet-group/offer.fastly.sse-load-test/internal/sse"
)

var log logging.Logger
var conf config.Config

func init() {
	conf = config.Load()

	logLvl, err := logging.LevelFromString(conf.LogLevel)
	if err != nil {
		panic(err)
	}

	consoleCore, err := zap.NewConsoleCore(zap.ConsoleConfig{
		ThresholdLevel: logLvl,
	})
	if err != nil {
		panic(err)
	}

	cores := []zapcore.Core{consoleCore}

	log, err = zap.NewLogger(cores...)
	if err != nil {
		panic(err)
	}
}

func buildOfferClient(hc *http.Client) *offer.Client {
	loc, err := time.LoadLocation(conf.Location)
	if err != nil {
		panic(err)
	}

	return offer.NewClient(conf.OfferHost, hc, loc)
}

func buildSSEClient(hc *http.Client) *sse.Client {
	return sse.NewClient(conf.OfferHost, hc)
}

func buildHTTPClient() *http.Client {

	hc := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10000,
			IdleConnTimeout:     35 * time.Second,
			MaxIdleConns:        10000,
			DisableCompression:  true,
		},

		Timeout: 35 * time.Second,
	}

	return hc
}
