package entity

import "time"

// Storage is one stored object. Key is the only identifier exposed upward;
// Path is internal to the object store and must never leak through an API.
type Storage struct {
	ID           uint64    `gorm:"column:id;primaryKey"`
	UserID       uint64    `gorm:"column:user_id;not null"`
	Key          string    `gorm:"column:key;size:32;not null;uniqueIndex"`
	Path         string    `gorm:"column:path;type:text;not null"`
	Bucket       string    `gorm:"column:bucket;size:255;not null"`
	OriginalName string    `gorm:"column:original_name;size:255;not null"`
	MIMEType     string    `gorm:"column:mime_type;size:255;not null"`
	SizeBytes    int64     `gorm:"column:size_bytes;not null"`
	Meta         *string   `gorm:"column:meta;type:jsonb"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (Storage) TableName() string {
	return "storages"
}
