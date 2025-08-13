package main

import (
	"time"

	"github.com/AleGaliev/kubercontroller/internal/agent"
)

func main() {
	counter := 0
	for {
		agent.Run(counter)
		time.Sleep(1 * time.Second)
	}
}
