package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			lock.Lock()
			q := queue
			lock.Unlock()

			queueSize.Set(float64(q))
			time.Sleep(time.Second)
		}
	}()
}

var (
	queueSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "queue_size",
	})
	lock  = &sync.Mutex{}
	queue int
)

func main() {
	recordMetrics()
	go consume()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/consume", addToQueue)
	http.ListenAndServe(":2112", nil)
}

func consume() {
	for {
		// Sleep a random duration between 0.5s and 1.5s between consuming from the queue
		x := rand.Float64()
		minD := time.Millisecond * 500
		maxD := time.Millisecond * 1500
		duration := time.Duration(float64(minD) + (float64(maxD)-float64(minD))*x)
		time.Sleep(duration)

		// Consume from the queue
		lock.Lock()
		didWork := false
		if queue > 0 {
			queue -= 1
			didWork = true
		}
		q := queue
		lock.Unlock()

		if didWork {
			fmt.Println("Queue size", q)
		}
	}
}

func addToQueue(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	queue++
	lock.Unlock()
}
