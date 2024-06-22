package utils

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func StartCronJobs() {
	c := cron.New()
	c.AddFunc("0 * * * *", runEngine)
	c.Start()
	cronCount := len(c.Entries())
	fmt.Printf("Setup %d cron jobs: \n", cronCount)

}

func runEngine() {
	fmt.Println("Starting REngine")
}
