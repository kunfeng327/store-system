package routes

import (
	mrand "math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"store/backend/models"
)

// Setup 注册全部路由
func Setup(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")

	// ---- 用户 ----
	api.POST("/register", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required,min=2,max=50"`
			Password string `json:"password" binding:"required,min=6"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "用户名至少2位、密码至少6位"})
			return
		}
		var exist models.User
		if err := db.Where("username = ?", req.Username).First(&exist).Error; err == nil {
			c.JSON(400, gin.H{"error": "用户名已存在"})
			return
		}
		u := models.User{Username: req.Username, Coins: 10000}
		if err := u.HashPassword(req.Password); err != nil {
			c.JSON(500, gin.H{"error": "服务器错误"})
			return
		}
		if err := db.Create(&u).Error; err != nil {
			c.JSON(500, gin.H{"error": "注册失败"})
			return
		}
		token, _ := u.GenerateToken()
		c.JSON(200, gin.H{"token": token, "user": u})
	})

	api.POST("/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "参数错误"})
			return
		}
		var u models.User
		if err := db.Where("username = ?", req.Username).First(&u).Error; err != nil || !u.CheckPassword(req.Password) {
			c.JSON(401, gin.H{"error": "用户名或密码错误"})
			return
		}
		token, _ := u.GenerateToken()
		c.JSON(200, gin.H{"token": token, "user": u})
	})

	// ---- 商品 ----
	api.GET("/products", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
		if page < 1 {
			page = 1
		}
		if size < 1 || size > 50 {
			size = 20
		}
		q := db.Model(&models.Product{})
		if cat := c.Query("category"); cat != "" {
			q = q.Where("category = ?", cat)
		}
		if kw := c.Query("keyword"); kw != "" {
			q = q.Where("name LIKE ? OR character LIKE ? OR series LIKE ?", "%"+kw+"%", "%"+kw+"%", "%"+kw+"%")
		}
		var total int64
		q.Count(&total)
		var list []models.Product
		q.Order("id").Offset((page - 1) * size).Limit(size).Find(&list)
		c.JSON(200, gin.H{"total": total, "list": list})
	})

	api.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		var p models.Product
		if err := db.First(&p, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "商品不存在"})
			return
		}
		// 评分统计
		type stat struct {
			Rating int `json:"rating"`
			Cnt    int `json:"cnt"`
		}
		var stats []stat
		db.Model(&models.Review{}).Select("rating, count(*) as cnt").Where("product_id = ?", p.ID).Group("rating").Scan(&stats)
		dist := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}
		sum, cnt := 0, 0
		for _, s := range stats {
			dist[s.Rating] = s.Cnt
			sum += s.Rating * s.Cnt
			cnt += s.Cnt
		}
		avg := 0.0
		if cnt > 0 {
			avg = float64(sum) / float64(cnt)
		}
		c.JSON(200, gin.H{"product": p, "rating_avg": avg, "rating_count": cnt, "rating_dist": dist})
	})

	// ---- 评价（淘宝风格）----
	api.GET("/products/:id/reviews", func(c *gin.Context) {
		id := c.Param("id")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
		rating := c.Query("rating") // 星级筛选
		q := db.Model(&models.Review{}).Where("product_id = ?", id)
		if rating != "" {
			q = q.Where("rating = ?", rating)
		}
		if hasFollow := c.Query("has_follow"); hasFollow == "1" {
			q = q.Where("has_follow = ?", true)
		}
		var total int64
		q.Count(&total)
		var list []models.Review
		q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&list)
		c.JSON(200, gin.H{"total": total, "list": list})
	})

	auth := api.Group("", AuthMiddleware())
	{
		auth.POST("/products/:id/reviews", func(c *gin.Context) {
			var req struct {
				Rating   int    `json:"rating" binding:"required,min=1,max=5"`
				Content  string `json:"content" binding:"required,min=1,max=1000"`
				FollowUp string `json:"follow_up"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "请填写评分（1-5星）和评价内容"})
				return
			}
			uid := c.GetUint("user_id")
			username := c.GetString("username")
			rv := models.Review{
				ProductID: parseUint(c.Param("id")),
				UserID:    uid,
				Username:  username,
				Rating:    req.Rating,
				Content:   req.Content,
				FollowUp:  req.FollowUp,
				HasFollow: strings.TrimSpace(req.FollowUp) != "",
			}
			if err := db.Create(&rv).Error; err != nil {
				c.JSON(500, gin.H{"error": "发布失败"})
				return
			}
			c.JSON(200, rv)
		})

		// ---- 购买 ----
		auth.POST("/orders", func(c *gin.Context) {
			var req struct {
				ProductID     uint   `json:"product_id" binding:"required"`
				Quantity      int    `json:"quantity" binding:"required,min=1"`
				PaymentMethod string `json:"payment_method"` // wechat/alipay/union
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "参数错误"})
				return
			}
			payMethod, ok := validPayMethod(req.PaymentMethod)
			if !ok {
				c.JSON(400, gin.H{"error": "支付方式仅支持 wechat（微信）/ alipay（支付宝）/ union（银联）"})
				return
			}
			uid := c.GetUint("user_id")
			var user models.User
			var p models.Product
			if err := db.First(&user, uid).Error; err != nil {
				c.JSON(401, gin.H{"error": "用户不存在"})
				return
			}
			if err := db.First(&p, req.ProductID).Error; err != nil {
				c.JSON(404, gin.H{"error": "商品不存在"})
				return
			}
			if p.Stock < req.Quantity {
				c.JSON(400, gin.H{"error": "库存不足"})
				return
			}
			total := p.Price * float64(req.Quantity)
			if user.Coins < int64(total) {
				c.JSON(400, gin.H{"error": "金币不足，去抽盒碰碰运气吧"})
				return
			}
			err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Model(&user).UpdateColumn("coins", gorm.Expr("coins - ?", int64(total))).Error; err != nil {
					return err
				}
				if err := tx.Model(&p).Updates(map[string]interface{}{
					"stock": gorm.Expr("stock - ?", req.Quantity),
					"sales": gorm.Expr("sales + ?", req.Quantity),
				}).Error; err != nil {
					return err
				}
				order := models.Order{UserID: uid, Total: total, Type: "buy", PaymentMethod: payMethod, Items: []models.OrderItem{{
					ProductID: p.ID, Name: p.Name, Image: p.Image, Rarity: p.Rarity, Price: p.Price, Quantity: req.Quantity,
				}}}
				return tx.Create(&order).Error
			})
			if err != nil {
				c.JSON(500, gin.H{"error": "下单失败"})
				return
			}
			c.JSON(200, gin.H{"message": "购买成功", "total": total, "payment_method": payMethod})
		})

		auth.GET("/orders", func(c *gin.Context) {
			uid := c.GetUint("user_id")
			var orders []models.Order
			db.Preload("Items").Where("user_id = ?", uid).Order("created_at DESC").Find(&orders)
			c.JSON(200, orders)
		})

		// ---- 盲盒抽卡 ----
		auth.POST("/box/draw", func(c *gin.Context) {
			var req struct {
				PaymentMethod string `json:"payment_method"`
			}
			_ = c.ShouldBindJSON(&req) // 允许空 body
			payMethod, _ := validPayMethod(req.PaymentMethod)
			uid := c.GetUint("user_id")
			var user models.User
			if err := db.First(&user, uid).Error; err != nil {
				c.JSON(401, gin.H{"error": "用户不存在"})
				return
			}
			const boxPrice = 60
			if user.Coins < boxPrice {
				c.JSON(400, gin.H{"error": "金币不足"})
				return
			}
			// 稀有度加权随机：N 60% / R 25% / SR 12% / SSR 3%
			rarity := weightedRarity()
			var pool []models.Product
			db.Where("rarity = ?", rarity).Find(&pool)
			if len(pool) == 0 {
				db.Where("rarity = ?", "N").Find(&pool)
				rarity = "N"
			}
			p := pool[randIntn(len(pool))]

			err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Model(&user).UpdateColumn("coins", gorm.Expr("coins - ?", boxPrice)).Error; err != nil {
					return err
				}
				if err := tx.Model(&p).Updates(map[string]interface{}{
					"stock": gorm.Expr("stock - 1"),
					"sales": gorm.Expr("sales + 1"),
				}).Error; err != nil {
					return err
				}
				order := models.Order{UserID: uid, Total: boxPrice, Type: "box", PaymentMethod: payMethod, Items: []models.OrderItem{{
					ProductID: p.ID, Name: p.Name, Image: p.Image, Rarity: p.Rarity, Price: p.Price, Quantity: 1,
				}}}
				if err := tx.Create(&order).Error; err != nil {
					return err
				}
				return tx.Create(&models.BlindBoxRecord{UserID: uid, ProductID: p.ID}).Error
			})
			if err != nil {
				c.JSON(500, gin.H{"error": "抽取失败"})
				return
			}
			c.JSON(200, gin.H{"product": p, "coins_left": user.Coins - boxPrice, "drawn_at": time.Now(), "payment_method": payMethod})
		})

		auth.GET("/box/records", func(c *gin.Context) {
			uid := c.GetUint("user_id")
			var records []models.BlindBoxRecord
			db.Preload("Product").Where("user_id = ?", uid).Order("created_at DESC").Limit(50).Find(&records)
			c.JSON(200, records)
		})

		auth.GET("/me", func(c *gin.Context) {
			var user models.User
			if err := db.First(&user, c.GetUint("user_id")).Error; err != nil {
				c.JSON(401, gin.H{"error": "用户不存在"})
				return
			}
			c.JSON(200, user)
		})
	}
}

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "请先登录"})
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(models.JWTSsecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "登录已过期，请重新登录"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("username", claims["username"].(string))
		c.Next()
	}
}

// validPayMethod 校验支付方式，默认微信
func validPayMethod(m string) (string, bool) {
	switch m {
	case "", "wechat":
		return "wechat", true
	case "alipay", "union":
		return m, true
	default:
		return "", false
	}
}

// weightedRarity 稀有度加权随机
func weightedRarity() string {
	r := mrand.Intn(100)
	switch {
	case r < 60:
		return "N"
	case r < 85:
		return "R"
	case r < 97:
		return "SR"
	default:
		return "SSR"
	}
}

func parseUint(s string) uint {
	n, _ := strconv.ParseUint(s, 10, 64)
	return uint(n)
}

func randIntn(n int) int {
	if n <= 0 {
		return 0
	}
	return mrand.Intn(n)
}
