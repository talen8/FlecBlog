package repository

import (
	"flec_blog/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ThemeRepository 主题仓储
type ThemeRepository struct {
	db *gorm.DB
}

// NewThemeRepository 创建主题仓储
func NewThemeRepository(db *gorm.DB) *ThemeRepository {
	return &ThemeRepository{db: db}
}

// Get 根据 slug 获取主题
func (r *ThemeRepository) Get(slug string) (*model.ThemeInstance, error) {
	var theme model.ThemeInstance
	if err := r.db.Where("slug = ?", slug).First(&theme).Error; err != nil {
		return nil, err
	}
	return &theme, nil
}

// GetActive 获取当前激活主题
func (r *ThemeRepository) GetActive() (*model.ThemeInstance, error) {
	var theme model.ThemeInstance
	if err := r.db.Where("is_active = ?", true).First(&theme).Error; err != nil {
		return nil, err
	}
	return &theme, nil
}

// List 获取全部主题
func (r *ThemeRepository) List() ([]model.ThemeInstance, error) {
	var themes []model.ThemeInstance
	err := r.db.Order("is_active DESC, created_at DESC").Find(&themes).Error
	return themes, err
}

// DeactivateAll 取消所有主题激活状态
func (r *ThemeRepository) DeactivateAll() error {
	return r.db.Model(&model.ThemeInstance{}).Where("is_active = ?", true).Update("is_active", false).Error
}

// SyncTheme 同步主题全部数据，冲突时完整覆盖
func (r *ThemeRepository) SyncTheme(theme *model.ThemeInstance) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "slug"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"version",
			"author",
			"description",
			"license",
			"repo",
			"schema",
			"config",
			"menus",
			"is_active",
			"updated_at",
		}),
	}).Create(theme).Error
}

// UpdateConfig 精确更新主题配置
func (r *ThemeRepository) UpdateConfig(slug string, config string) error {
	return r.db.Model(&model.ThemeInstance{}).
		Where("slug = ?", slug).
		Select("config", "updated_at").
		Updates(map[string]interface{}{"config": config, "updated_at": gorm.Expr("CURRENT_TIMESTAMP")}).
		Error
}

// UpdateMenus 精确更新主题菜单
func (r *ThemeRepository) UpdateMenus(slug string, menus string) error {
	return r.db.Model(&model.ThemeInstance{}).
		Where("slug = ?", slug).
		Select("menus", "updated_at").
		Updates(map[string]interface{}{"menus": menus, "updated_at": gorm.Expr("CURRENT_TIMESTAMP")}).
		Error
}

// UpdateDefaultSchemaIfEmpty 当默认主题的 schema 为空时，用给定内容填充
func (r *ThemeRepository) UpdateDefaultSchemaIfEmpty(schema string) error {
	return r.db.Model(&model.ThemeInstance{}).
		Where("slug = ? AND (schema IS NULL OR schema::text = '{}')", "default").
		Select("schema", "updated_at").
		Updates(map[string]interface{}{"schema": schema, "updated_at": gorm.Expr("CURRENT_TIMESTAMP")}).
		Error
}

// ExistsByFileURL 检查主题配置或菜单是否引用该文件
func (r *ThemeRepository) ExistsByFileURL(url string) (bool, error) {
	var count int64
	err := r.db.Model(&model.ThemeInstance{}).
		Where("config::text LIKE ? OR menus::text LIKE ?", "%"+url+"%", "%"+url+"%").
		Count(&count).Error
	return count > 0, err
}
