package main

import (
	"github.com/jordanknott/monitor/internal/commands/monitor"
	_ "github.com/lib/pq"
)

func main() {
	monitor.Execute()
}
