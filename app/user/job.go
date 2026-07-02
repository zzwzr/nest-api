package user

import (
	"context"

	"github.com/robfig/cron/v3"
)

type Job struct{}

func RegisterJobs(c *cron.Cron) {
	svc := UserService{}
	c.AddFunc("* * * * *", func() {
		_ = svc.LogActiveUsers(context.Background())
	})
}
