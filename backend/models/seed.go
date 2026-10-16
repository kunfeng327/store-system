package models

import "gorm.io/gorm"

// Seed 种子数据：首次运行时插入示例谷子
func Seed(db *gorm.DB) {
	var count int64
	db.Model(&Product{}).Count(&count)
	if count > 0 {
		return
	}

	products := []Product{
		// 盲盒池商品（各种稀有度）
		{Name: "初音未来 雪初音吧唧", Character: "初音未来", Series: "VOCALOID", Category: "吧唧", Rarity: "SSR", Price: 129, Stock: 5, Sales: 1200, Image: "🎵"},
		{Name: "雷电将军 鸣神岛立牌", Character: "雷电将军", Series: "原神", Category: "立牌", Rarity: "SSR", Price: 259, Stock: 8, Sales: 980, Image: "⚡"},
		{Name: "芙莉莲 精灵色纸", Character: "芙莉莲", Series: "葬送的芙莉莲", Category: "色纸", Rarity: "SR", Price: 89, Stock: 15, Sales: 2100, Image: "✨"},
		{Name: "甘城猫猫 像素挂件", Character: "猫猫", Series: "砂糖苹果 Fairy 帽子店", Category: "挂件", Rarity: "SR", Price: 69, Stock: 20, Sales: 1800, Image: "🐱"},
		{Name: "五条悟 无下限吧唧", Character: "五条悟", Series: "咒术回战", Category: "吧唧", Rarity: "SR", Price: 79, Stock: 12, Sales: 3200, Image: "🌀"},
		{Name: "菜狗 憨憨挂件", Character: "菜狗", Series: "菜狗", Category: "挂件", Rarity: "R", Price: 39, Stock: 50, Sales: 5600, Image: "🥬"},
		{Name: "帕恰狗 春日立牌", Character: "帕恰狗", Series: "三丽鸥", Category: "立牌", Rarity: "R", Price: 49, Stock: 30, Sales: 4100, Image: "🐶"},
		{Name: "玉桂狗 云朵色纸", Character: "玉桂狗", Series: "三丽鸥", Category: "色纸", Rarity: "R", Price: 45, Stock: 25, Sales: 3900, Image: "☁️"},
		{Name: "数码宝贝 亚古兽吧唧", Character: "亚古兽", Series: "数码宝贝", Category: "吧唧", Rarity: "N", Price: 19, Stock: 100, Sales: 8200, Image: "🦖"},
		{Name: "皮卡丘 闪电挂件", Character: "皮卡丘", Series: "宝可梦", Category: "挂件", Rarity: "N", Price: 25, Stock: 80, Sales: 7600, Image: "⚡"},
		{Name: "哆啦A梦 口袋立牌", Character: "哆啦A梦", Series: "哆啦A梦", Category: "立牌", Rarity: "N", Price: 29, Stock: 60, Sales: 6500, Image: "🔔"},
		{Name: "蜡笔小新 动感色纸", Character: "小新", Series: "蜡笔小新", Category: "色纸", Rarity: "N", Price: 22, Stock: 70, Sales: 5900, Image: "🌈"},
	}
	db.Create(&products)

	// 演示用户 + 示例评价
	demo := User{Username: "demo", Coins: 10000}
	if err := demo.HashPassword("demo123"); err == nil {
		db.Create(&demo)
	}

	if len(products) >= 4 {
		reviews := []Review{
			{ProductID: products[0].ID, UserID: demo.ID, Username: "demo", Rating: 5, Content: "雪初音太美了！做工精致，别针很结实，无限回购！", FollowUp: "半年了依旧没有褪色，官方质量 yyds", HasFollow: true},
			{ProductID: products[0].ID, UserID: demo.ID, Username: "吃谷老炮", Rating: 4, Content: "整体不错，就是包装有点简陋，建议加气泡膜。"},
			{ProductID: products[1].ID, UserID: demo.ID, Username: "鸣神岛岛主", Rating: 5, Content: "雷神立牌颜值爆表！底座稳，立绘细节拉满，桌面排面！"},
			{ProductID: products[2].ID, UserID: demo.ID, Username: "精灵旅人", Rating: 5, Content: "芙莉莲色纸印刷很清晰，色彩通透，值得收藏。"},
			{ProductID: products[3].ID, UserID: demo.ID, Username: "猫猫推", Rating: 4, Content: "像素风很可爱，挂包上回头率高，就是有点小。"},
		}
		db.Create(&reviews)
	}
}
