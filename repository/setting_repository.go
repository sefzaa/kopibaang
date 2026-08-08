package repository

import (
	"context"
	"gorm.io/gorm"
	"kopibang/domain"
)

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) domain.SettingRepository {
	return &settingRepository{db}
}

// Struct private untuk mapping tabel system_settings
type systemSetting struct {
	SettingKey   string `gorm:"primaryKey;column:setting_key"`
	SettingValue string `gorm:"column:setting_value"`
}

func (r *settingRepository) GetSetting(ctx context.Context, key string) (string, error) {
	var setting systemSetting
	err := r.db.WithContext(ctx).Where("setting_key = ?", key).First(&setting).Error
	if err != nil {
		return "", err
	}
	return setting.SettingValue, nil
}

func (r *settingRepository) UpdateSetting(ctx context.Context, key string, value string) error {
	setting := systemSetting{
		SettingKey:   key,
		SettingValue: value,
	}
	// Save akan melakukan insert atau update berdasarkan Primary Key
	return r.db.WithContext(ctx).Save(&setting).Error
}