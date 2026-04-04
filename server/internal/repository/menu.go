package repository

import (
	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// MenuRepository 菜单仓储
type MenuRepository struct {
	db *gorm.DB
}

// NewMenuRepository 创建菜单仓储
func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

// Create 创建菜单
func (r *MenuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

// Update 更新菜单
func (r *MenuRepository) Update(menu *model.Menu) error {
	return r.db.Save(menu).Error
}

// Delete 删除菜单
func (r *MenuRepository) Delete(id uint) error {
	return r.db.Delete(&model.Menu{}, id).Error
}

// Get 根据ID获取菜单
func (r *MenuRepository) Get(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// GetWithParent 根据ID获取菜单及其父菜单
func (r *MenuRepository) GetWithParent(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.Preload("Parent").First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// GetTree 获取菜单树
func (r *MenuRepository) GetTree(menuType string, isEnabled *bool) ([]model.Menu, error) {
	var menus []model.Menu
	query := r.db.Model(&model.Menu{}).Where("parent_id IS NULL")

	if menuType != "" {
		query = query.Where("type = ?", menuType)
	}
	if isEnabled != nil {
		query = query.Where("is_enabled = ?", *isEnabled)
	}

	// 加载主菜单
	err := query.Order("sort ASC, id ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 为每个主菜单加载子菜单
	for i := range menus {
		if err := r.loadChildren(&menus[i], isEnabled); err != nil {
			return nil, err
		}
	}

	return menus, nil
}

// loadChildren 加载子菜单
func (r *MenuRepository) loadChildren(menu *model.Menu, isEnabled *bool) error {
	query := r.db.Model(&model.Menu{}).Where("parent_id = ?", menu.ID)

	if isEnabled != nil {
		query = query.Where("is_enabled = ?", *isEnabled)
	}

	// 只加载直接子菜单，不再递归
	err := query.Order("sort ASC, id ASC").Find(&menu.Children).Error
	if err != nil {
		return err
	}

	return nil
}


// ExistsByIcon 检查是否有菜单图标引用该文件
func (r *MenuRepository) ExistsByIcon(url string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Menu{}).Where("icon = ?", url).Count(&count).Error
	return count > 0, err
}

// HasChildren 检查菜单是否有子菜单
func (r *MenuRepository) HasChildren(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Menu{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

// GetChildrenCount 获取子菜单数量
func (r *MenuRepository) GetChildrenCount(id uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Menu{}).Where("parent_id = ?", id).Count(&count).Error
	return count, err
}

// CheckParentExists 检查父菜单是否存在
func (r *MenuRepository) CheckParentExists(parentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Menu{}).Where("id = ?", parentID).Count(&count).Error
	return count > 0, err
}

// DeleteChildren 删除所有子菜单
func (r *MenuRepository) DeleteChildren(parentID uint) error {
	return r.db.Where("parent_id = ?", parentID).Delete(&model.Menu{}).Error
}

// UpgradeChildren 将子菜单升级为主菜单（设置parent_id为NULL）
func (r *MenuRepository) UpgradeChildren(parentID uint) error {
	return r.db.Model(&model.Menu{}).Where("parent_id = ?", parentID).Update("parent_id", nil).Error
}

// GetChildren 获取子菜单列表
func (r *MenuRepository) GetChildren(parentID uint) ([]model.Menu, error) {
	var children []model.Menu
	err := r.db.Where("parent_id = ?", parentID).Find(&children).Error
	return children, err
}
