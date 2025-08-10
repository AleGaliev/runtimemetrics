package agent

import (
	"github.com/AleGaliev/kubercontroller/internal/repository"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
)

var (
	metRuntime  = runtime.MemStats{}
	sendMetrics = repository.SendMetrics{
		URL:    make(map[string]url.URL),
		Client: &http.Client{},
	}
	PollCount = 1
)

func Run() {
	runtime.ReadMemStats(&metRuntime)

	mericsNameValue := repository.ConvertMemStatsInNameMetrics(metRuntime)

	sendMetrics.InitMetrics(mericsNameValue)
	sendMetrics.InitMetrics(map[string]string{
		"PollCount":   strconv.Itoa(PollCount),
		"RandomValue": strconv.FormatFloat(rand.Float64(), 'f', 10, 64),
	})

	if PollCount%5 == 0 {
		sendMetrics.DoRequest()
	}

	PollCount++

}
