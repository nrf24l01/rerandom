package modes

import (
	"log"
	"os"
	"time"

	"github.com/nrf24l01/rerandom/gs_sync/config"
	"github.com/nrf24l01/rerandom/gs_sync/redis"
	"github.com/nrf24l01/rerandom/gs_sync/task"
)

func RunSync() {
	// Load configuration from environment variables or a config file
	cfg := config.BuildConfigFromEnv()
	var gscfg *config.GSConfig
	if os.Getenv("PRODUCTION_ENV") != "true" {
		gscfg = config.BuildFromFile("service_account.json")
	} else {
		gscfg = config.GSBuildFromEnv()
	}
	redis := redis.InitRedisFromCFG(cfg)

	for {
		log.Printf("Starting sync...")
		sheet := task.DownloadSheet(cfg.SpreadsheetID, cfg.SheetName, *gscfg)

		users := sheet.ClearUsers

		// Save users to Redis
		err := redis.SaveStruct(users, 0)
		if err != nil {
			log.Fatalf("failed to save users to redis: %v", err)
		}

		log.Printf("Successfully synced %d users to Redis", len(users))
		time.Sleep(20 * time.Second)
	}
}
