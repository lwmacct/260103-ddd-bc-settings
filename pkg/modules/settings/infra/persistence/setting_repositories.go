package persistence

import (
	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
	"gorm.io/gorm"
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
