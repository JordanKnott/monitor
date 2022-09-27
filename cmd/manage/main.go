package main

import (
	"github.com/jordanknott/monitor/internal/commands/manage"
	_ "github.com/lib/pq"
)

func main() {
	manage.Execute()
}
