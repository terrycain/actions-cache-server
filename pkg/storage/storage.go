package storage

import (
	"errors"
	"io"

	"github.com/terrycain/actions-cache-server/pkg/s"

	s3 "github.com/terrycain/actions-cache-server/pkg/storage/aws-s3"
	"github.com/terrycain/actions-cache-server/pkg/storage/azureblob"
	"github.com/terrycain/actions-cache-server/pkg/storage/disk"
)

type Backend interface {
	Setup() error
	Type() string
	Write(repoKey string, cacheID int, r io.Reader, start, end int, size int64) (string, int64, error)
	Delete(repoKey string, partData string) error

	// Finalise Takes a list of upload parts, and somehow concatenates them and returns a path which can be passed to GenerateArchiveURL
	Finalise(repoKey string, cacheID int, parts []s.CachePart) (string, error)
	GenerateArchiveURL(scheme, host, repoKey, path string) (string, error)
	GetFilePath(key string) (string, error)
}

func GetStorageBackend(backend, connectionString string) (Backend, error) {
	var b Backend
	var err error

	switch backend {
	case "disk":
		b, err = disk.New(connectionString)
	case "s3":
		b, err = s3.New(connectionString)
	case "azureblob":
		b, err = azureblob.New(connectionString)
	default:
		return nil, errors.New("invalid storage backend")
	}

	if err != nil {
		return nil, err
	}

	if err = b.Setup(); err != nil {
		return nil, err
	}

	return b, nil
}
