package main

import (
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/terrycain/actions-cache-server/pkg/database"
	"github.com/terrycain/actions-cache-server/pkg/storage"
	"github.com/terrycain/actions-cache-server/pkg/utils/logging"
	"github.com/terrycain/actions-cache-server/pkg/web"
)

var cli struct {
	// Database backends
	DBSqlite   string `env:"DB_SQLITE" required:"" xor:"db" help:"SQLite filepath e.g. /tmp/db.sqlite"`
	DBPostgres string `env:"DB_POSTGRES" required:"" xor:"db" help:"Postgres URI e.g. postgresql://blah"`

	// Storage backends
	StorageDisk      string `env:"STORAGE_DISK" required:"" xor:"storage" help:"Use disk storage for cache data e.g. /tmp/cache"`
	StorageS3        string `env:"STORAGE_S3" required:"" xor:"storage" name:"storage-s3" help:"Use S3 storage for cache data e.g. s3://bucket"`
	StorageAzureBlob string `env:"STORAGE_AZUREBLOB" required:"" xor:"storage" name:"storage-azureblob" help:"Use Azure Blob Storage for cache data e.g. connectionstring with ;Container=blah on the end"`

	// Misc
	LogLevel             string `env:"LOG_LEVEL" default:"info" enum:"debug,info,warn,error"`
	ListenAddress        string `env:"LISTEN_ADDR" default:"0.0.0.0:8080" help:"Listen address e.g. 0.0.0.0:8080"`
	MetricsListenAddress string `env:"METRICS_LISTEN_ADDR" default:"0.0.0.0:9102" help:"Listen address for prometheus metrics e.g. 0.0.0.0:9102"`
	Debug                bool   `env:"DEBUG" help:"Enable debug mode"`
}

func main() {
	kong.Parse(&cli)

	logging.SetupLogging(cli.LogLevel)

	var databaseBackendName, dbConnectionString string
	if cli.DBSqlite != "" {
		databaseBackendName = "sqlite"
		dbConnectionString = cli.DBSqlite
	}

	var storageBackendName, storageConnectionString string
	if cli.StorageDisk != "" {
		storageBackendName = "disk"
		storageConnectionString = cli.StorageDisk
	}
	if cli.StorageS3 != "" {
		storageBackendName = "s3"
		storageConnectionString = cli.StorageS3
	}
	if cli.StorageAzureBlob != "" {
		storageBackendName = "azureblob"
		storageConnectionString = cli.StorageAzureBlob
	}

	dbBackend, err := database.GetBackend(databaseBackendName, dbConnectionString)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initiate database backend")
	}

	storageBackend, err := storage.GetStorageBackend(storageBackendName, storageConnectionString)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initiate storage backend")
	}

	handlers := web.Handlers{
		Database: dbBackend,
		Storage:  storageBackend,
		Debug:    cli.Debug,
	}

	router := web.GetRouter(cli.MetricsListenAddress, handlers, true)

	log.Info().Msgf("Listening on %s", cli.ListenAddress)
	if err = router.Run(cli.ListenAddress); err != nil {
		log.Fatal().Err(err).Msg("Failed HTTP server loop")
	}
}
