package main

import (
	"github.com/AleGaliev/kubercontroller/internal/agent"
	"time"
)

func main() {
	for {
		agent.Run()
		time.Sleep(2 * time.Second)
	}
}
