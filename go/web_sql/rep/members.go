package rep

var Members = map[string][]string{
	"Cart": {
		"ID",
		"UserID",
		"ProductID",
		"Quantity",
		"AddTime",
	},

	"Category": {
		"ID",
		"Name",
		"ParentID",
	},

	"Coupon": {
		"ID",
		"Code",
		"Type",
		"Discount",
		"Minimum",
		"UserID",
		"ProductID",
		"ExpirationDate",
		"Status",
	},

	"DeliveryAddress": {
		"ID",
		"UserID",
		"Phone",
		"Address",
		"Name",
	},

	"OrderItem": {
		"ID",
		"OrderID",
		"ProductID",
		"Quantity",
		"UnitPrice",
		"TotalPrice",
	},

	"Order": {
		"ID",
		"UserID",
		"Total",
		"Status",
		"CreatedTime",
		"UpdateTime",
	},

	"Product": {
		"ID",
		"Name",
		"Description",
		"Price",
		"Stock",
		"Type",
		"Category",
		"Seller",
		"IsActive",
		"Icon",
	},

	"Review": {
		"ID",
		"UserID",
		"ProductID",
		"Rating",
		"Comment",
		"Time",
	},

	"Shipping": {
		"ID",
		"OrderItemID",
		"TrackingNumber",
		"Carrier",
		"Status",
		"EstimatedDeliveredTime",
		"CreateTime",
		"ShippedTime",
		"CompletedTime",
	},

	"User": {
		"ID",
		"Username",
		"Password",
		"Email",
		"Phone",
	},
}
