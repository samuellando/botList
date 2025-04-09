package cron

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

type CronJob struct {
	Crontab   string
	Message   []byte
}

type CronService struct {messageQueue chan []byte}

func Create(messageQueue chan []byte) CronService {
    return CronService{messageQueue: messageQueue}
}

func (s *CronService) AddJob(j CronJob) {
    go func() {
        parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional)
        sched, err := parser.Parse(j.Crontab)
        if err != nil {
            panic(err)
        }
        for {
            next := sched.Next(time.Now())
            until := next.Sub(time.Now())
            fmt.Println("Sleeping", until)
            if until > 0 {
                time.Sleep(until)
            } 
            s.messageQueue <- j.Message
        }
    }()
}
