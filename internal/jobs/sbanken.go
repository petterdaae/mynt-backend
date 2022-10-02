package jobs

import (
	"backend/internal/resources/admin"
	sResource "backend/internal/resources/sbanken"
	"backend/internal/utils"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func sbanken() {
	log.Info("Starting Sbanken synch job")

	database := utils.InitDatabase()

	adminResource := admin.Configure(database)
	useIDs, err := adminResource.ListUserIds()
	if err != nil {
		log.Error(fmt.Errorf("failed to list user ids: %w", err))
	}

	for _, userID := range useIDs {
		sbankenResource := sResource.Configure(userID, database)
		err = sbankenResource.Synchronize()
		if err != nil {
			log.Warn(fmt.Errorf("failed to synchronize sbanken user"))
		}
	}

	log.Info("Sbanken synch job finished successfully")
}
