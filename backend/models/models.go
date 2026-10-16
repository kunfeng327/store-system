package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 全局 JWT 配置（main 中初始化）
var JWTSsecret string
var JWTExpire time.Duration

// User 用户
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:50" json:"username"`
	PasswordHash string    `json:"-"`
	Coins        int64     `gorm:"default:10000" json:"coins"` // 金币余额，注册送 10000
	CreatedAt    time.Time `json:"created_at"`
}

// HashPassword 加密密码
func (u *User) HashPassword(pwd string) error {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	u.PasswordHash = string(b)
	return err
}

// CheckPassword 校验密码
func (u *User) CheckPassword(pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pwd)) == nil
}

// GenerateToken 生成 JWT
func (u *User) GenerateToken() (string, error) {
	claims := jwt.MapClaims{
		"user_id":  u.ID,
		"username": u.Username,
		"exp":      time.Now().Add(JWTExpire).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSsecret))
}

// Product 谷子商品
type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100" json:"name"`   // 商品名
	Character string    `gorm:"size:100" json:"character"` // 角色
	Series    string    `gorm:"size:100" json:"series"`    // 作品
	Category  string    `gorm:"size:20;index" json:"category"` // 吧唧/立牌/挂件/色纸/手办
	Rarity    string    `gorm:"size:10" json:"rarity"`    // N/R/SR/SSR
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	Sales     int       `gorm:"default:0" json:"sales"`
	Image     string    `gorm:"size:255" json:"image"` // emoji 或图片 URL
	CreatedAt time.Time `json:"created_at"`
}

// Review 商品评价
type Review struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProductID  uint      `gorm:"index" json:"product_id"`
	UserID     uint      `json:"user_id"`
	Username   string    `gorm:"size:50" json:"username"`
	Rating     int       `json:"rating"` // 1-5 星
	Content    string    `gorm:"type:text" json:"content"`
	FollowUp   string    `gorm:"type:text" json:"follow_up"` // 追评
	HasFollow  bool      `json:"has_follow"`
	CreatedAt  time.Time `json:"created_at"`
}

// Order 订单
type Order struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	UserID        uint        `gorm:"index" json:"user_id"`
	Total         float64     `json:"total"`
	Type          string      `gorm:"size:10" json:"type"`                 // buy 购买 / box 盲盒
	PaymentMethod string      `gorm:"size:20;default:wechat" json:"payment_method"` // wechat 微信 / alipay 支付宝 / union 银联
	Items         []OrderItem `json:"items"`
	CreatedAt     time.Time   `json:"created_at"`
}

// OrderItem 订单商品
type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"index" json:"order_id"`
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Image     string  `json:"image"`
	Rarity    string  `json:"rarity"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

// BlindBoxRecord 盲盒抽取记录
type BlindBoxRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	ProductID uint      `json:"product_id"`
	Product   Product   `json:"product"`
	CreatedAt time.Time `json:"created_at"`
}
