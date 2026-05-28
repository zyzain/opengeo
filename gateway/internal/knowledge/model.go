package knowledge

import "time"

// Entity 知识图谱实体
// 用于管理品牌、产品、概念等实体，提升AI引用权威性
// 支持实体与内容关联，帮助AI搜索引擎更好地理解和引用内容
type Entity struct {
	// ID 实体唯一标识
	ID int64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// UserID 创建者用户ID，用于多用户隔离
	UserID int64 `gorm:"index;not null" json:"user_id"`

	// EntityName 实体名称，如品牌名、产品名、人名等
	EntityName string `gorm:"size:128;not null" json:"entity_name"`

	// EntityType 实体类型：brand(品牌), product(产品), concept(概念), person(人物), place(地点)
	EntityType string `gorm:"size:32;not null;index" json:"entity_type"`

	// EntityData 实体详细数据，JSON格式存储
	// 示例: {"description": "...", "website": "...", "logo": "..."}
	EntityData string `gorm:"type:text" json:"entity_data"`

	// AuthorityLinks 权威链接，JSON数组格式
	// 用于关联百科、官网等权威节点，提升AI信任度
	// 示例: ["https://en.wikipedia.org/wiki/xxx", "https://official-site.com"]
	AuthorityLinks string `gorm:"type:text" json:"authority_links"`

	// ContentCount 关联内容数量，统计有多少内容引用了该实体
	ContentCount int32 `gorm:"default:0" json:"content_count"`

	// CreatedAt 创建时间
	CreatedAt time.Time `gorm:"type:datetime" json:"created_at"`

	// UpdatedAt 最后更新时间
	UpdatedAt time.Time `gorm:"type:datetime" json:"updated_at"`
}

// TableName 指定表名
func (Entity) TableName() string {
	return "entities"
}
