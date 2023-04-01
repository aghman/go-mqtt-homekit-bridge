package main

import (
	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"

	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create the switch accessory.
	a := accessory.NewTemperatureSensor(accessory.Info{
		Name:         "MicroTemp",
		Manufacturer: "Austin",
	})

	// Store the data in the "./db" directory.
	fs := hap.NewFsStore("./db")

	// Create the hap server.
	server, err := hap.NewServer(fs, a.A)
	if err != nil {
		// stop if an error happens
		log.Panic(err)
	}
	server.Pin = "00201003"
	// Setup a listener for interrupts and SIGTERM signals
	// to stop the server.
	stopChannel := make(chan os.Signal)
	tempChannel := make(chan float32)
	signal.Notify(stopChannel, os.Interrupt)
	signal.Notify(stopChannel, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-stopChannel
		// Stop delivering signals.
		signal.Stop(stopChannel)
		// Cancel the context to stop the server.
		cancel()
	}()
	a.TempSensor.CurrentTemperature.SetValue(10.0)
	go checkPairing(tempChannel, stopChannel)

	// Run the server.
	server.ListenAndServe(ctx)
}

func checkPairing(tempChannel chan float32, stopChannel chan os.Signal) {
	for {
		select {
		case <-stopChannel:
			return
		default:
			a.TempSensor.CurrentTemperature.SetValue(tempChannel)
		}
	}
}
