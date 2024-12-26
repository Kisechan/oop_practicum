package rep

import (
	"time"
)

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

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Order struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"not null" json:"user_id"`
	Total       *float64  `gorm:"type:decimal(10,2);default:null" json:"total"`
	Status      string    `gorm:"type:enum('pending','paid','shipping','completed');default:'pending';not null" json:"status"`
	CreatedTime time.Time `gorm:"not null" json:"created_time"`
	UpdateTime  time.Time `gorm:"not null;autoUpdateTime" json:"update_time"`
	ProductID   int       `gorm:"not null" json:"product_id"`
	Quantity    int       `gorm:"not null" json:"quantity"`

	User    User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Product struct {
	ID          int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string  `gorm:"type:varchar(100);not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int     `gorm:"not null" json:"stock"`
	Type        string  `gorm:"type:enum('presale','normal');default:'normal';not null" json:"type"`
	Category    string  `gorm:"default:null" json:"category"`
	Seller      string  `gorm:"type:varchar(100);not null" json:"seller"`
	IsActive    string  `gorm:"type:enum('true','false');default:'true';not null" json:"is_active"`
	Icon        *string `gorm:"default:null" json:"icon"`

	Reviews []Review `gorm:"foreignKey:ProductID"`
}

type Review struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"not null" json:"user_id"`
	ProductID int       `gorm:"not null" json:"product_id"`
	Rating    string    `gorm:"type:enum('1','2','3','4','5');default:'5';not null" json:"rating"`
	Comment   string    `gorm:"type:text" json:"comment"`
	Time      time.Time `gorm:"not null;autoUpdateTime" json:"time"`

	User    User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type User struct {
	ID       int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string  `gorm:"type:varchar(50);not null" json:"username"`
	Password string  `gorm:"type:varchar(255);not null" json:"password"`
	Email    *string `gorm:"type:varchar(255)" json:"email"`
	Phone    string  `gorm:"type:varchar(32);not null" json:"phone"`
	Address  *string `gorm:"type:varchar(255)" json:"address"`

	Carts   []Cart   `gorm:"foreignKey:UserID"`
	Orders  []Order  `gorm:"foreignKey:UserID"`
	Reviews []Review `gorm:"foreignKey:UserID"`
	Coupons []Coupon `gorm:"foreignKey:UserID"`
}
