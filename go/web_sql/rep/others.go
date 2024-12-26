package rep

import (
	"gorm.io/gorm"
)

func GetCartsByUserID(db *gorm.DB, userID int) ([]Cart, error) {
	var carts []Cart
	err := db.Where("user_id = ?", userID).Find(&carts).Error
	return carts, err
}

func GetCouponsByUserID(db *gorm.DB, userID int) ([]Coupon, error) {
	var coupons []Coupon
	err := db.Where("user_id = ?", userID).Find(&coupons).Error
	return coupons, err
}

func GetOrdersByUserID(db *gorm.DB, userID int) ([]Order, error) {
	var orders []Order
	err := db.Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func GetReviewsByProductID(db *gorm.DB, productID int) ([]Review, error) {
	var reviews []Review
	err := db.Where("product_id = ?", productID).Find(&reviews).Error
	return reviews, err
}

func UpdateOrderStruct(db *gorm.DB, id int, data Order) error {
	// 使用 Omit 方法，排除 CreatedTime 字段
	return db.Model(&Order{}).Where("id = ?", id).Omit("CreatedTime").Updates(data).Error

}
