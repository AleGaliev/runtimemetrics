package main

import (
	"flag"
	"time"

	"github.com/AleGaliev/kubercontroller/internal/agent"
)

func main() {
	counter := 0
	flag.Parse()
	for {
		agent.Run(counter)
		time.Sleep(1 * time.Second)
	}
}
