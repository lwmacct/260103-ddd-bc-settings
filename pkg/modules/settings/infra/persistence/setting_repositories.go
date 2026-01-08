package persistence

import (
	"gorm.io/gorm"

	"github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/domain/setting"
)

// SettingRepositories 聚合配置定义读写仓储
type SettingRepositories struct {
	Command         setting.CommandRepository
	Query           setting.QueryRepository
	CategoryCommand setting.SettingCategoryCommandRepository
	CategoryQuery   setting.SettingCategoryQueryRepository
}

// NewSettingRepositories 创建配置定义仓储聚合实例
func NewSettingRepositories(db *gorm.DB) SettingRepositories {
	return SettingRepositories{
		Command:         NewSettingCommandRepository(db),
		Query:           NewSettingQueryRepository(db),
		CategoryCommand: NewSettingCategoryCommandRepository(db),
		CategoryQuery:   NewSettingCategoryQueryRepository(db),
	}
}
