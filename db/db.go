package db

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"monks.co/piano-alone/data"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/id"
	"monks.co/piano-alone/songs"
)

type DB struct {
	*gorm.DB
}

func OpenDB(path string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&game.Performance{}); err != nil {
		return nil, err
	}

	var count int64
	if err = db.Model(&game.Performance{}).
		Where("date = ?", data.PerformanceDate).
		Count(&count).
		Error; err != nil {
		return nil, err
	} else if count == 0 {
		if err := db.Create(&game.Performance{
			Configuration: &game.Configuration{
				PerformanceID: id.Random128(),
				Score:         songs.PreludeOpus3No2Bytes,
				Title:         "Prelude in C♯ Minor",
				Composer:      "Sergei Rachmaninoff",
			},
			Date:       data.PerformanceDate,
			IsFeatured: true,
		}).Error; err != nil {
			return nil, err
		}
	}

	return &DB{db}, nil
}

func (db *DB) GetPerformance(id string) (*game.Performance, error) {
	var p game.Performance
	if err := db.Model(&game.Performance{}).
		Where("id = ?", id).
		First(&p).
		Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (db *DB) GetScheduledPerformances() ([]*game.Performance, error) {
	var ps []*game.Performance
	if err := db.Model(&game.Performance{}).
		Select("id", "date", "title", "composer", "is_complete", "player_count", "is_featured").
		Where("is_complete = false").
		Find(&ps).
		Error; err != nil {
		return nil, err
	}
	return ps, nil
}

func (db *DB) DeletePerformance(id string) error {
	if err := db.
		Where("id = ?", id).
		Delete(&game.Performance{}).
		Error; err != nil {
		return err
	}
	return nil
}

func (db *DB) GetFeaturedPerformances() ([]*game.Performance, error) {
	var ps []*game.Performance
	if err := db.Model(&game.Performance{}).
		Select("id", "date", "title", "composer", "player_count", "is_complete", "is_featured").
		Where("is_featured = true").
		Order("date asc, title").
		Find(&ps).
		Error; err != nil {
		return nil, err
	}
	return ps, nil
}

func (db *DB) SchedulePerformance(p *game.Performance) error {
	if err := db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (db *DB) StartPerformance(conf *game.Configuration) error {
	p := &game.Performance{
		Configuration: conf,
		Date:          time.Now(),
	}
	if err := db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (db *DB) SaveRendition(id string, playerCount int, rendition []byte) error {
	if err := db.Model(&game.Performance{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"is_complete":  true,
			"rendition":    rendition,
			"player_count": playerCount,
		}).
		Error; err != nil {

		return err
	}
	return nil
}

func (db *DB) Feature(id string) error {
	if err := db.Model(&game.Performance{}).
		Where("id = ?", id).
		Update("is_featured", true).
		Error; err != nil {

		return err
	}
	return nil
}

func (db *DB) GetMIDIFile(id string) ([]byte, error) {
	var p game.Performance
	if err := db.First(&p, "id = ?", id).
		Error; err != nil {
		return nil, err
	}
	return p.Rendition, nil
}
