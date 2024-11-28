package rep

import (
	"fmt"

	"gorm.io/gorm"
)

func Create[T any](db *gorm.DB, record *T) error {
	return db.Create(record).Error
}

func GetID[T any](db *gorm.DB, id int) (*T, error) {
	var record T
	err := db.First(&record, id).Error
	return &record, err
}

func GetField[T any](db *gorm.DB, column string, value any) (*T, error) {
	var record T
	err := db.Where(fmt.Sprintf("%s = ?", column), value).First(&record).Error
	return &record, err
}

func GetAll[T any](db *gorm.DB) ([]T, error) {
	var records []T
	err := db.Find(&records).Error
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
