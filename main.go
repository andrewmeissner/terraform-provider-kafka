package main

import (
	"context"
	"flag"
	"log"

	"github.com/andrewmeissner/terraform-provider-kafka/kafka"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: kafka.Provider}

	if debugMode {
		if err := plugin.Debug(context.Background(), "registry.terraform.io/andrewmeissner/kafka", opts); err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
