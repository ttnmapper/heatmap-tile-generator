package types

import (
	"database/sql/driver"
	"time"
)

type MysqlAggGridcell struct {
	Id int `db:"id"`

	GtwId       string   `db:"gtw_id"`
	FirstSample NullTime `db:"first_sample"`
	LastSample  NullTime `db:"last_sample"`
	UpdatedAt   NullTime `db:"updated_at"`

	X int `db:"x"`
	Y int `db:"y"`

	BucketHigh     int64 `db:"bucket_high"`
	Bucket100      int64 `db:"bucket_100"`
	Bucket105      int64 `db:"bucket_105"`
	Bucket110      int64 `db:"bucket_110"`
	Bucket115      int64 `db:"bucket_115"`
	Bucket120      int64 `db:"bucket_120"`
	Bucket125      int64 `db:"bucket_125"`
	Bucket130      int64 `db:"bucket_130"`
	Bucket135      int64 `db:"bucket_135"`
	Bucket140      int64 `db:"bucket_140"`
	Bucket145      int64 `db:"bucket_145"`
	BucketLow      int64 `db:"bucket_low"`
	BucketNoSignal int64 `db:"bucket_no_signal"`
}

type MysqlTileToRedraw struct {
	Id         int       `db:"id"`
	X          int       `db:"x"`
	Y          int       `db:"y"`
	Z          int       `db:"z"`
	LastQueued time.Time `db:"last_queued"`
}

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
