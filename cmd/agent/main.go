package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/AleGaliev/kubercontroller/internal/agent"
	"github.com/AleGaliev/kubercontroller/internal/repository"
)

func main() {
	clientCfg := repository.NewClientConfig()
	agentCfg := agent.NewAgentConfig(clientCfg)

	flag.Parse()

	for {
		if err := agentCfg.Run(); err != nil {
			fmt.Println(err)
		}
		time.Sleep(1 * time.Second)
	}
}
