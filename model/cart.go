// model/cart.go
package model

import ( // import โมเดล Product ที่อยู่ในไฟล์ product.go
	"time"
)

// Cart คือตารางสำหรับเก็บข้อมูลรถเข็นของลูกค้า
type Cart struct {
	CartID     int       `gorm:"column:cart_id;primary_key;AUTO_INCREMENT"`
	CustomerID int       `gorm:"column:customer_id;NOT NULL"`
	CartName   string    `gorm:"column:cart_name"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`

	Items []CartItem `gorm:"foreignKey:CartID"` // ความสัมพันธ์กับ CartItem โดยใช้ CartID
}

func (m *Cart) TableName() string {
	return "cart"
}

// CartItem คือตารางสำหรับเก็บข้อมูลสินค้าที่อยู่ในรถเข็น
type CartItem struct {
	CartItemID int       `gorm:"column:cart_item_id;primary_key;AUTO_INCREMENT"`
	CartID     int       `gorm:"column:cart_id;NOT NULL"`
	ProductID  int       `gorm:"column:product_id;NOT NULL"`
	Quantity   int       `gorm:"column:quantity;NOT NULL"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`

	Product Product `gorm:"foreignKey:ProductID"` // ความสัมพันธ์กับ Product โดยใช้ ProductID
}

func (m *CartItem) TableName() string {
	return "cart_item"
}
