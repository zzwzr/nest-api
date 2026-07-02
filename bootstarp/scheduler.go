package bootstrap

import (
	"time"

	"nest-api/pkg/logger"

	"github.com/robfig/cron/v3"
)

func InitScheduler() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logger.Error("load timezone failed", err)
		loc = time.FixedZone("CST", 8*3600)
	}

	c := cron.New(cron.WithLocation(loc))

	// user.RegisterJobs(c)

	c.Start()
}
