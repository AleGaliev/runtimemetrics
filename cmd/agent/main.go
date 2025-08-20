package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/AleGaliev/kubercontroller/internal/agent"
	"github.com/AleGaliev/kubercontroller/internal/logger"
	"github.com/AleGaliev/kubercontroller/internal/repository"
)

func main() {
	logServer, err := logger.CreateLogger()
	if err != nil {
		panic(err)
	}
	clientCfg := repository.NewClientConfig(logServer)
	agentCfg, err := agent.NewAgentConfig(clientCfg)
	if err != nil {
		log.Fatalf("error parsing agent config: %v", err)
	}

	flag.Parse()

	for {
		if err := agentCfg.Run(); err != nil {
			fmt.Println(err)
		}
		time.Sleep(1 * time.Second)
	}
}
