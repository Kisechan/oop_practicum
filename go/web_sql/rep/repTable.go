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

type Category struct {
	ID       int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"type:varchar(100);not null" json:"name"`
	ParentID int    `gorm:"default:0" json:"parentid"`
}

type Coupon struct {
	ID             int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Code           string     `gorm:"type:varchar(32);not null" json:"code"`
	Type           string     `gorm:"type:enum('percent','minus');default:'minus';not null" json:"type"`
	Discount       float64    `gorm:"type:decimal(5,2);default:0.00;not null" json:"discount"`
	Minimum        float64    `gorm:"type:decimal(10,2);default:0.00;not null" json:"minimum"`
	UserID         int        `gorm:"not null" json:"user_id"`
	ExpirationDate *time.Time `gorm:"default:null" json:"expiration_date"`
	Status         string     `gorm:"type:enum('available','used','expired');default:'available';not null" json:"status"`

	User    User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product *Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// type DeliveryAddress struct {
// 	ID      int    `gorm:"primaryKey;autoIncrement" json:"id"`
// 	UserID  int    `gorm:"not null" json:"user_id"`
// 	Phone   string `gorm:"type:varchar(15);not null" json:"phone"`
// 	Address string `gorm:"type:text;not null" json:"address"`
// 	Name    string `gorm:"type:varchar(40);not null" json:"name"`

// 	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
// }

// type OrderItem struct {
// 	ID         int     `gorm:"primaryKey;autoIncrement" json:"id"`
// 	OrderID    int     `gorm:"not null" json:"order_id"`
// 	ProductID  int     `gorm:"not null" json:"product_id"`
// 	Quantity   int     `gorm:"not null" json:"quantity"`
// 	UnitPrice  float64 `gorm:"type:decimal(10,2);not null" json:"unit_price"`
// 	TotalPrice float64 `gorm:"type:decimal(10,2);not null" json:"total_price"`

// 	Order   Order   `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
// 	Product Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
// }

type Order struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"not null" json:"user_id"`
	Total       *float64  `gorm:"type:decimal(10,2);default:null" json:"total"`
	Status      string    `gorm:"type:enum('pending','paid','shipping','completed');default:'pending';not null" json:"status"`
	CreatedTime time.Time `gorm:"not null;autoUpdateTime" json:"created_time"`
	UpdateTime  time.Time `gorm:"not null;autoUpdateTime" json:"update_time"`
	ProductID   int       `gorm:"not null" json:"product_id"`
	Quantity    int       `gorm:"not null" json:"quantity"`

	User    User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
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
	Icon        *string `gorm:"defualt:null" json:"icon"`

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

// type Shipping struct {
// 	ID                     int        `gorm:"primaryKey;autoIncrement" json:"id"`
// 	OrderItemID            int        `gorm:"not null" json:"order_itemid"`
// 	TrackingNumber         *string    `gorm:"type:varchar(255);default:null" json:"tracking_number"`
// 	Carrier                *string    `gorm:"type:varchar(255);default:null" json:"carrier"`
// 	Status                 string     `gorm:"type:enum('pending_shipment','pending_collect','delivering','pending_pickup','pickedup','error');default:'pending_shipment';not null" json:"status"`
// 	EstimatedDeliveredTime *time.Time `gorm:"default:null" json:"estimated_delivered_time"`
// 	CreateTime             *time.Time `gorm:"default:null" json:"create_time"`
// 	ShippedTime            *time.Time `gorm:"default:null" json:"shipped_time"`
// 	CompletedTime          *time.Time `gorm:"default:null" json:"completed_time"`

// 	OrderItem OrderItem `gorm:"foreignKey:OrderItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
// }

type User struct {
	ID       int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string  `gorm:"type:varchar(50);not null" json:"username"`
	Password string  `gorm:"type:varchar(255);not null" json:"password"`
	Email    *string `gorm:"type:varchar(255)" json:"email"`
	Phone    string  `gorm:"type:varchar(32);not null" json:"phone"`
	Address  *string `gorm:"type:varchar(255)" json:"address"`

	Carts []Cart `gorm:"foreignKey:UserID"`
	// DeliveryAddresses []DeliveryAddress `gorm:"foreignKey:UserID"`
	Orders  []Order  `gorm:"foreignKey:UserID"`
	Reviews []Review `gorm:"foreignKey:UserID"`
}
