package graph

import (
	"github.com/jordanknott/monitor/internal/config"
	"github.com/jordanknott/monitor/internal/db"
)

type Resolver struct {
	Data   db.Data
	Config config.AppConfig
}
