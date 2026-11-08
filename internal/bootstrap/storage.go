package bootstrap

import (
	"context"
	"log/slog"

	"gorm.io/gorm"

	storagerepo "go-fiber-vue/internal/app/application/core/repository/storage"
	"go-fiber-vue/internal/config"
	platformstorage "go-fiber-vue/internal/platform/storage"
)

// newStorageRepository builds the optional storage layer. It returns nil (and
// logs one warning) when storage is not configured or the database is
// unavailable, so the API still starts without an object store.
func newStorageRepository(cfg *config.Config, db *gorm.DB) storagerepo.Repository {
	if cfg == nil || !cfg.IsStorageConfigured() {
		slog.Warn("storage disabled: set STORAGE_ENDPOINT, STORAGE_BUCKET, STORAGE_ACCESS_KEY and STORAGE_SECRET_KEY to enable it")
		return nil
	}
	if db == nil {
		slog.Warn("storage disabled: database unavailable")
		return nil
	}

	objects, err := platformstorage.NewS3(context.Background(), platformstorage.S3Config{
		Endpoint:  cfg.StorageEndpoint,
		Region:    cfg.StorageRegion,
		Bucket:    cfg.StorageBucket,
		AccessKey: cfg.StorageAccessKey,
		SecretKey: cfg.StorageSecretKey,
		PathStyle: cfg.StoragePathStyle,
	})
	if err != nil {
		slog.Warn("storage disabled: client init failed", "error", err)
		return nil
	}

	return storagerepo.NewRepository(objects, storagerepo.NewMetaRepository(db), cfg.StorageBucket)
}
