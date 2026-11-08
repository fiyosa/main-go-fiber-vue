// Package storage is the application-facing storage repository. It combines
// object I/O with the storages table so callers only ever deal with keys.
package storage

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"path"
	"strings"
	"time"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
	platformstorage "go-fiber-vue/internal/platform/storage"
)

// SaveInput describes one upload. Server-generated fields intentionally have
// no representation here: key, path, size and timestamps are owned by this
// layer.
type SaveInput struct {
	UserID       uint64
	Folder       string
	OriginalName string
	MIMEType     string
	Reader       io.Reader
	Meta         json.RawMessage
}

// File is a stored object together with its metadata. The caller owns Body
// and must close it.
type File struct {
	Storage *entity.Storage
	Body    io.ReadCloser
}

// Close releases the underlying object stream.
func (f *File) Close() error {
	if f == nil || f.Body == nil {
		return nil
	}
	return f.Body.Close()
}

// MetaRepository persists the storages table.
type MetaRepository interface {
	Create(ctx context.Context, row *entity.Storage) error
	FindByKey(ctx context.Context, key string) (*entity.Storage, error)
	DeleteByKey(ctx context.Context, key string) error
}

// Repository is the only surface consumers should depend on.
type Repository interface {
	Save(ctx context.Context, in SaveInput) (*entity.Storage, error)
	Read(ctx context.Context, key string) (*File, error)
	Delete(ctx context.Context, key string) error
}

type repository struct {
	objects platformstorage.ObjectStore
	meta    MetaRepository
	bucket  string
}

// NewRepository wires object I/O and metadata persistence. It returns nil when
// either dependency is missing so bootstrap can expose an optional layer.
func NewRepository(objects platformstorage.ObjectStore, meta MetaRepository, bucket string) Repository {
	if objects == nil || meta == nil {
		return nil
	}
	return &repository{objects: objects, meta: meta, bucket: bucket}
}

// Save uploads the object first and records its row afterwards. When the row
// insert fails the uploaded object is removed so the bucket does not keep an
// orphan.
func (r *repository) Save(ctx context.Context, in SaveInput) (*entity.Storage, error) {
	if r == nil {
		return nil, exceptions.ServiceUnavailable(nil)
	}
	if in.UserID == 0 {
		return nil, exceptions.Validation(map[string][]string{"user_id": {"user id is required"}})
	}

	folder, err := platformstorage.NormalizeFolder(in.Folder)
	if err != nil {
		return nil, exceptions.Validation(map[string][]string{"folder": {err.Error()}})
	}

	if len(in.Meta) > 0 && !json.Valid(in.Meta) {
		return nil, exceptions.Validation(map[string][]string{"meta": {"meta must be valid JSON"}})
	}

	key, err := platformstorage.NewKey()
	if err != nil {
		return nil, exceptions.Internal(err)
	}
	objectPath := path.Join(folder, key)

	mimeType := strings.TrimSpace(in.MIMEType)
	if mimeType == "" {
		mimeType = platformstorage.DefaultContentType
	}
	originalName := strings.TrimSpace(in.OriginalName)

	_, size, err := r.objects.Save(ctx, platformstorage.SaveInput{
		Path:        objectPath,
		ContentType: mimeType,
		Body:        in.Reader,
	})
	if err != nil {
		return nil, exceptions.Internal(err)
	}

	var meta *string
	if len(in.Meta) > 0 {
		value := string(in.Meta)
		meta = &value
	}

	now := time.Now().UTC()
	row := &entity.Storage{
		UserID:       in.UserID,
		Key:          key,
		Path:         objectPath,
		Bucket:       r.bucket,
		OriginalName: originalName,
		MIMEType:     mimeType,
		SizeBytes:    size,
		Meta:         meta,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := r.meta.Create(ctx, row); err != nil {
		// Best effort: the insert error is what the caller needs to see.
		if cleanupErr := r.objects.Delete(context.WithoutCancel(ctx), objectPath); cleanupErr != nil {
			slog.Warn("storage: failed to remove orphaned object after insert error",
				"path", objectPath, "error", cleanupErr)
		}
		return nil, err
	}

	return row, nil
}

// Read streams the object behind key. Keys that are not in the table, or whose
// object disappeared from the bucket, are reported as not found.
func (r *repository) Read(ctx context.Context, key string) (*File, error) {
	if r == nil {
		return nil, exceptions.ServiceUnavailable(nil)
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, exceptions.Validation(map[string][]string{"key": {"key is required"}})
	}

	row, err := r.meta.FindByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, exceptions.NotFound(nil)
	}

	body, err := r.objects.Read(ctx, row.Path)
	if err != nil {
		if errors.Is(err, platformstorage.ErrObjectNotFound) {
			return nil, exceptions.NotFound(err)
		}
		return nil, exceptions.Internal(err)
	}

	return &File{Storage: row, Body: body}, nil
}

// Delete removes the object before the row so metadata never claims an object
// that is already gone. A failed object delete leaves the row untouched.
func (r *repository) Delete(ctx context.Context, key string) error {
	if r == nil {
		return exceptions.ServiceUnavailable(nil)
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return exceptions.Validation(map[string][]string{"key": {"key is required"}})
	}

	row, err := r.meta.FindByKey(ctx, key)
	if err != nil {
		return err
	}
	if row == nil {
		return exceptions.NotFound(nil)
	}

	if err := r.objects.Delete(ctx, row.Path); err != nil {
		return exceptions.Internal(err)
	}
	return r.meta.DeleteByKey(ctx, key)
}

// ensure the gorm implementation matches the interface at compile time.
var _ MetaRepository = (*gormMetaRepository)(nil)

type gormMetaRepository struct {
	db *gorm.DB
}

// NewMetaRepository returns a GORM-backed metadata repository, or nil when db
// is nil so the storage layer can be disabled without a database.
func NewMetaRepository(db *gorm.DB) MetaRepository {
	if db == nil {
		return nil
	}
	return &gormMetaRepository{db: db}
}

func (r *gormMetaRepository) Create(ctx context.Context, row *entity.Storage) error {
	return exceptions.FromDB(r.db.WithContext(ctx).Create(row).Error)
}

func (r *gormMetaRepository) FindByKey(ctx context.Context, key string) (*entity.Storage, error) {
	var row entity.Storage
	err := r.db.WithContext(ctx).
		Where("key = ?", key).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, exceptions.FromDB(err)
	}
	return &row, nil
}

func (r *gormMetaRepository) DeleteByKey(ctx context.Context, key string) error {
	return exceptions.FromDB(r.db.WithContext(ctx).
		Where("key = ?", key).
		Delete(&entity.Storage{}).Error)
}
