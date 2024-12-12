package rep

import (
	"fmt"

	"gorm.io/gorm"
)

func Create[T any](db *gorm.DB, record *T) error {
	return db.Create(record).Error
}

func GetID[T any](db *gorm.DB, id int, preloads ...string) (*T, error) {
	var record T
	query := db
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.First(&record, id).Error
	return &record, err
}

func GetField[T any](db *gorm.DB, column string, value any, preloads ...string) (*T, error) {
	var record T
	query := db
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.Where(fmt.Sprintf("%s = ?", column), value).First(&record).Error
	return &record, err
}

func GetAll[T any](db *gorm.DB, preloads ...string) ([]T, error) {
	var records []T
	query := db
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.Find(&records).Error
	return records, err
}

func UpdateOne[T any](db *gorm.DB, id int, field string, value any) error {
	return db.Model(new(T)).Where("id = ?", id).Update(field, value).Error
}

func UpdateMany[T any](db *gorm.DB, id int, updates map[string]any) error {
	return db.Model(new(T)).Where("id = ?", id).Updates(updates).Error
}

func DeleteID[T any](db *gorm.DB, id int) error {
	return db.Delete(new(T), id).Error
}
