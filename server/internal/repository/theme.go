package repository

import (
	"flec_blog/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ThemeRepository 主题配置仓储
type ThemeRepository struct {
	db *gorm.DB
}

// NewThemeRepository 创建主题配置仓储
func NewThemeRepository(db *gorm.DB) *ThemeRepository {
	return &ThemeRepository{db: db}
}

// List 获取全部主题配置项
func (r *ThemeRepository) List() ([]model.ThemeConfig, error) {
	var items []model.ThemeConfig
	err := r.db.Order("key ASC").Find(&items).Error
	return items, err
}

// Get 根据键获取配置项
func (r *ThemeRepository) Get(key string) (*model.ThemeConfig, error) {
	var item model.ThemeConfig
	if err := r.db.Where("key = ?", key).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// UpsertMany 批量写入配置项
func (r *ThemeRepository) UpsertMany(items []model.ThemeConfig) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&items).Error
}

// Upsert 写入单个配置项
func (r *ThemeRepository) Upsert(key string, value string) error {
	return r.UpsertMany([]model.ThemeConfig{{Key: key, Value: value}})
}

// ExistsByFileURL 检查主题配置或菜单是否引用该文件
func (r *ThemeRepository) ExistsByFileURL(url string) (bool, error) {
	var count int64
	err := r.db.Model(&model.ThemeConfig{}).
		Where("value LIKE ?", "%"+url+"%").
		Count(&count).Error
	return count > 0, err
}
