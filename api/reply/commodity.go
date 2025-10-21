package reply

type HierarchicCommodityCategory struct {
	ID            int64                          `json:"id"`
	Level         int                            `json:"level"`
	ParentId      int64                          `json:"parent_id"`
	Name          string                         `json:"name"`
	IconImg       string                         `json:"icon_img"`
	Rank          int                            `json:"rank"`
	SubCategories []*HierarchicCommodityCategory `json:"sub_categories"` // 分类的子分类
}

type CommodityCategory struct {
	ID       int64  `json:"id"`
	Level    int    `json:"level"`
	ParentId int64  `json:"parent_id"`
	Name     string `json:"name"`
	IconImg  string `json:"icon_img"`
	Rank     int    `json:"rank"`
}
