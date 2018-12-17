package main

import (
	"fmt"
	"net/http"
	"os"

	"./broker"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

func main() {
	brokerLogger := lager.NewLogger("broker")
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	brokerLogger.Info("Starting broker")
	brokerCredentials := brokerapi.BrokerCredentials{
		Username: "scott",
		Password: "yolo",
	}

	var serviceBroker brokerapi.ServiceBroker
	var err error

	serviceBroker, err = broker.New()
	if err != nil {
		fmt.Println("oh no")
	}

	fmt.Println("oh ya")
	brokerAPI := brokerapi.New(serviceBroker, brokerLogger, brokerCredentials)

	http.Handle("/", brokerAPI)

	brokerLogger.Fatal("http-listen", http.ListenAndServe(":3000", nil))
}
