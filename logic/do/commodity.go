package do

import (
	"time"
)

type CommodityCategory struct {
	ID        int64     `json:"id"`
	Level     int       `json:"level"`
	ParentId  int64     `json:"parent_id"`
	Name      string    `json:"name"`
	IconImg   string    `json:"icon_img"`
	Rank      int       `json:"rank"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HierarchicCommodityCategory 按等级划分的商品分类信息
type HierarchicCommodityCategory struct {
	ID            int64                          `json:"id"`
	Level         int                            `json:"level"`
	ParentId      int64                          `json:"parent_id"`
	Name          string                         `json:"name"`
	IconImg       string                         `json:"icon_img"`
	Rank          int                            `json:"rank"`
	SubCategories []*HierarchicCommodityCategory `json:"sub_categories"` // 分类的子分类
	CreatedAt     time.Time                      `json:"created_at"`
	UpdatedAt     time.Time                      `json:"updated_at"`
}

type Commodity struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Intro         string    `json:"intro"`
	CategoryId    int64     `json:"category_id"`
	CoverImg      string    `json:"cover_img"`
	Images        string    `json:"images"`
	DetailContent string    `json:"detail_content"`
	OriginalPrice int       `json:"original_price"`
	SellingPrice  int       `json:"selling_price"`
	StockNum      int       `json:"stock_num"`
	Tag           string    `json:"tag"`
	SellStatus    int       `json:"sell_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CommodityListElem  Element of commodity list
// 删掉，领域对象不应该跟接口的响应对象保存一致，而是应该能更通用地表示业务领域里的概念
// 领域服务的商品列表方法改为返回 []*do.Commodity，在App服务中在把[]*doCommodity  转换为 []*reply.CommodityListElem
//type CommodityListElem struct {
//	ID            int64     `json:"id"`
//	Name          string    `json:"name"`
//	Intro         string    `json:"intro"`
//	CategoryId    int64     `json:"category_id"`
//	CoverImg      string    `json:"cover_img"`
//	OriginalPrice int       `json:"original_price"`
//	SellingPrice  int       `json:"selling_price"`
//	Tag           string    `json:"tag"`
//	SellStatus    int       `json:"sell_status"`
//	CreatedAt     time.Time `json:"created_at"`
//}
