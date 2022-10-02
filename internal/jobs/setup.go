package jobs

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

func Schedule() {
	log.Info("Scheduling jobs")
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.TagsUnique()

	_, err := scheduler.Every(30).Second().Do(sbanken)
	if err != nil {
		log.Error(fmt.Errorf("failed to schedule sbanken synch job: %w", err))
	}

	scheduler.StartAsync()
	log.Info("All jobs scheduled successfully")
}
