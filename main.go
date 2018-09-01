package main

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/joho/godotenv"
	"./launchpad"
	"fmt"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		launchpad.GetLogger().Fatal("Failed to load .env!")
		os.Exit(1)
	}

	ee := launchpad.CheckEnv()

	if ee != nil {
		launchpad.GetLogger().Fatal("Missing environment variables, exiting")
		fmt.Println(ee)
		os.Exit(1)
	}

	if os.Getenv("GO_ENV") == "DEBUG" {
		launchpad.EnableDebug()
	}

	launchpad.GetDbInstance().Setup()

	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		done <- true
	}()

	launchpad.GetLogger().Info("Starting Launchpad server...")

	instance := launchpad.Instance{}

	launchpad.GetGithubInstance().Init()

	go instance.StartWebServer()
	go launchpad.StartStatsFetcher()

	<-done

	launchpad.GetLogger().Info("Stopping server...")
}
