package repository

import (
	"encoding/json"
	"errors"
	"fmt"

	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// SettingRepository 配置仓库
type SettingRepository struct {
	db *gorm.DB
}

// NewSettingRepository 创建配置仓库
func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

// GetByGroup 获取某个分组的所有配置
func (r *SettingRepository) GetByGroup(group string, isPublicOnly ...bool) (map[string]string, error) {
	var settings []model.Setting
	query := r.db.Where("\"group\" = ?", group)

	// 如果指定只返回公开配置，则添加过滤条件
	if len(isPublicOnly) > 0 && isPublicOnly[0] {
		query = query.Where("is_public = ?", true)
	}

	err := query.Order("id ASC").Find(&settings).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}

// UpdateGroup 更新配置
func (r *SettingRepository) UpdateGroup(group string, updates map[string]string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for key, value := range updates {
			var setting model.Setting
			err := tx.Where("\"group\" = ? AND \"key\" = ?", group, key).First(&setting).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("配置项不存在: %s", key)
			}
			if err != nil {
				return err
			}
			if err := tx.Model(&setting).Update("value", value).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ExistsByValueAndKeys 检查指定配置键中是否引用该文件
func (r *SettingRepository) ExistsByValueAndKeys(value string, keys []string) (bool, error) {
	var count int64
	query := r.db.Model(&model.Setting{}).Where("value = ?", value)
	if len(keys) > 0 {
		query = query.Where("\"key\" IN ?", keys)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// ExistsByJSONArrayField 检查 JSON 数组配置中是否包含指定 URL
func (r *SettingRepository) ExistsByJSONArrayField(value string, keyFieldMap map[string]string) (bool, error) {
	for key, field := range keyFieldMap {
		var setting model.Setting
		if err := r.db.Where("\"key\" = ?", key).First(&setting).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return false, err
		}
		// 解析 JSON 数组
		var items []map[string]interface{}
		if err := json.Unmarshal([]byte(setting.Value), &items); err != nil {
			continue
		}
		// 检查每个项的指定字段
		for _, item := range items {
			if v, ok := item[field].(string); ok && v == value {
				return true, nil
			}
		}
	}
	return false, nil
}
