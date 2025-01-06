package pages

import "time"

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Type        string  `json:"type"`
	Category    string  `json:"category"`
	Seller      string  `json:"seller"`
	IsActive    string  `json:"is_active"`
	Icon        *string `json:"icon"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
}

type Review struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	Rating    string    `json:"rating"`
	Comment   string    `json:"comment"`
	Time      time.Time `json:"time"`

	User User
}

type ReviewsResponse struct {
	Reviews []Review `json:"reviews"`
}

// User 结构体
type User struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Email    *string `json:"email"`
	Phone    string  `json:"phone"`
	Address  *string `json:"address"`

	Carts   []Cart   `json:"Carts"`
	Orders  []Order  `json:"Orders"`
	Coupons []Coupon `json:"Coupons"`
}

type Cart struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"not null" json:"user_id"`
	ProductID int       `gorm:"not null" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	AddTime   time.Time `gorm:"not null;autoUpdateTime" json:"add_time"`

	User    User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Coupon struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Code           string    `gorm:"type:varchar(32);not null" json:"code"`
	Discount       float64   `gorm:"type:decimal(5,2);default:0.00;not null" json:"discount"`
	Minimum        float64   `gorm:"type:decimal(10,2);default:0.00;not null" json:"minimum"`
	UserID         int       `gorm:"not null" json:"user_id"`
	ExpirationDate time.Time `gorm:"default:null" json:"expiration_date"`
	Status         string    `gorm:"type:enum('available','used','expired');default:'available';not null" json:"status"`

	User User
}

type Order struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"not null" json:"user_id"`
	Total       float64   `gorm:"type:decimal(10,2);default:null" json:"total"`
	Status      string    `gorm:"type:enum('pending','paid','shipping','completed');default:'pending';not null" json:"status"`
	CreatedTime time.Time `gorm:"not null" json:"created_time"`
	UpdateTime  time.Time `gorm:"not null;autoUpdateTime" json:"update_time"`
	ProductID   int       `gorm:"not null" json:"product_id"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	Discount    float64   `gorm:"type:decimal(10,2);not null" json:"discount"`
	Payable     float64   `gorm:"type:decimal(10,2);not null" json:"payable"`
	CouponCode  string    `gorm:"type:varchar(32);not null" json:"coupon_code"`
	OrderNumber string    `gorm:"type:varchar(255);not null" json:"order_number"`

	User    User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

var OrderStatus = map[string]string{
	"pending":   "待支付",
	"paid":      "已支付",
	"shipping":  "配送中",
	"completed": "已完成",
}

var Message = map[string]string{
	"Insufficient inventory": "库存不足",
	"Invalid coupon":         "无效优惠券",
	"Order completed":        "下单成功",
	"failed":                 "下单失败",
	"success":                "成功",
}
