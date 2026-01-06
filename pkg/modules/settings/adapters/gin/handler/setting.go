package handler

import (
	"errors"
	"log/slog"
	"runtime/debug"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"
	settingDomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/response"
)

// SettingHandler handles setting management operations (DDD+CQRS Use Case Pattern)
type SettingHandler struct {
	// Setting Command Handlers
	createHandler      *setting.CreateHandler
	updateHandler      *setting.UpdateHandler
	deleteHandler      *setting.DeleteHandler
	batchUpdateHandler *setting.BatchUpdateHandler

	// Setting Query Handlers
	getHandler            *setting.GetHandler
	listHandler           *setting.ListHandler
	listSchemaHandler     *setting.ListSettingsHandler
	publicSettingsHandler *setting.PublicSettingsHandler

	// Category Command Handlers
	createCategoryHandler *setting.CreateCategoryHandler
	updateCategoryHandler *setting.UpdateCategoryHandler
	deleteCategoryHandler *setting.DeleteCategoryHandler

	// Category Query Handlers
	getCategoryHandler    *setting.GetCategoryHandler
	listCategoriesHandler *setting.ListCategoriesHandler
}

// NewSettingHandler creates a new SettingHandler instance from SettingUseCases.
func NewSettingHandler(useCases *setting.SettingUseCases) *SettingHandler {
	return &SettingHandler{
		createHandler:         useCases.Create,
		updateHandler:         useCases.Update,
		deleteHandler:         useCases.Delete,
		batchUpdateHandler:    useCases.BatchUpdate,
		getHandler:            useCases.Get,
		listHandler:           useCases.List,
		listSchemaHandler:     useCases.ListSettings,
		publicSettingsHandler: useCases.PublicSettings,
		createCategoryHandler: useCases.CreateCategory,
		updateCategoryHandler: useCases.UpdateCategory,
		deleteCategoryHandler: useCases.DeleteCategory,
		getCategoryHandler:    useCases.GetCategory,
		listCategoriesHandler: useCases.ListCategories,
	}
}

// GetSettings 获取系统配置（扁平结构）
//
//	@Summary		配置列表
//	@Description	获取配置数据，扁平结构（每个 item 包含 category 和 group 字段供前端分组）。支持按分类过滤（懒加载）。
//	@Tags			settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query	handler.ListSettingsQuery	false	"查询参数"
//	@Success		200		{object}	response.DataResponse[[]setting.SettingsItemDTO]	"配置列表（扁平结构）"
//	@Failure		401		{object}	response.ErrorResponse									"未授权"
//	@Failure		403		{object}	response.ErrorResponse									"权限不足"
//	@Failure		404		{object}	response.ErrorResponse									"分类不存在"
//	@Failure		500		{object}	response.ErrorResponse									"服务器内部错误"
//	@Router			/api/admin/settings [get]
func (h *SettingHandler) GetSettings(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("GetSettings panic", "error", err, "stack", string(debug.Stack()))
			response.InternalError(c, "Internal server error")
			c.Abort()
		}
	}()

	var query ListSettingsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	categoryKey := query.Category

	// 检查 Handler 是否为 nil
	if h.listSchemaHandler == nil {
		response.InternalError(c, "List settings handler not initialized")
		return
	}

	// 调用 Schema Handler（返回层级结构）
	schema, err := h.listSchemaHandler.Handle(c.Request.Context(), setting.ListSettingsQuery{
		CategoryKey: categoryKey,
	})
	if err != nil {
		// 检查是否为分类不存在错误
		if categoryKey != "" && err.Error() == "category not found: "+categoryKey {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, schema)
}

// ListSettingsQuery 查询配置列表参数
type ListSettingsQuery struct {
	Category string `form:"category"`
}

// GetSetting 获取单个配置
//
//	@Summary		配置详情
//	@Description	根据配置键获取配置详情
//	@Tags			settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			key	path		string										true	"配置键"	example:"site_name"
//	@Success		200	{object}	response.DataResponse[setting.SettingDTO]	"配置详情"
//	@Failure		401	{object}	response.ErrorResponse						"未授权"
//	@Failure		403	{object}	response.ErrorResponse						"权限不足"
//	@Failure		404	{object}	response.ErrorResponse						"配置不存在"
//	@Router			/api/admin/settings/{key} [get]
func (h *SettingHandler) GetSetting(c *gin.Context) {
	key := c.Param("key")

	// 调用 Use Case Handler
	setting, err := h.getHandler.Handle(c.Request.Context(), setting.GetQuery{
		Key: key,
	})

	if err != nil {
		response.NotFoundMessage(c, err.Error())
		return
	}

	response.OK(c, setting)
}

// CreateSettingRequest 创建配置请求
type CreateSettingRequest struct {
	Key          string `json:"key" binding:"required" example:"site_name"`
	DefaultValue any    `json:"default_value" binding:"required"`
	CategoryID   uint   `json:"category_id" binding:"required" example:"1"`
	Group        string `json:"group" example:"basic"`
	ValueType    string `json:"value_type" example:"string"`
	Label        string `json:"label" example:"网站名称"`
	UIConfig     string `json:"ui_config" example:"{}"`
	Order        int    `json:"order" example:"0"`
}

// CreateSetting 创建配置
//
//	@Summary		创建配置
//	@Description	管理员创建新的系统配置项
//	@Tags			settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateSettingRequest						true	"配置信息"
//	@Success		201		{object}	response.DataResponse[setting.SettingDTO]	"配置创建成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/admin/settings [post]
func (h *SettingHandler) CreateSetting(c *gin.Context) {
	var req CreateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用 Use Case Handler 创建配置
	_, err := h.createHandler.Handle(c.Request.Context(), setting.CreateCommand{
		Key:          req.Key,
		DefaultValue: req.DefaultValue,
		CategoryID:   req.CategoryID,
		Group:        req.Group,
		ValueType:    req.ValueType,
		Label:        req.Label,
		UIConfig:     req.UIConfig,
		Order:        req.Order,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 查询完整的配置信息返回
	settingDTO, err := h.getHandler.Handle(c.Request.Context(), setting.GetQuery{
		Key: req.Key,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, settingDTO)
}

// UpdateSettingRequest 更新配置请求
type UpdateSettingRequest struct {
	DefaultValue any    `json:"default_value"`
	Label        string `json:"label" example:"更新后的标签"`
	UIConfig     string `json:"ui_config"`
	Order        int    `json:"order"`
}

// UpdateSetting 更新配置
//
//	@Summary		更新配置
//	@Description	管理员更新指定配置项的值和标签
//	@Tags			settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			key		path		string										true	"配置键"	example:"site_name"
//	@Param			request	body		UpdateSettingRequest						true	"更新信息"
//	@Success		200		{object}	response.DataResponse[setting.SettingDTO]	"配置更新成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		404		{object}	response.ErrorResponse						"配置不存在"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/admin/settings/{key} [put]
func (h *SettingHandler) UpdateSetting(c *gin.Context) {
	key := c.Param("key")

	var req UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用 Use Case Handler
	settingDTO, err := h.updateHandler.Handle(c.Request.Context(), setting.UpdateCommand{
		Key:          key,
		DefaultValue: req.DefaultValue,
		Label:        req.Label,
		UIConfig:     req.UIConfig,
		Order:        req.Order,
	})

	if err != nil {
		if errors.Is(err, settingDomain.ErrValidationFailed) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, settingDTO)
}

// DeleteSetting 删除配置
//
//	@Summary		删除配置
//	@Description	管理员删除指定的系统配置项
//	@Tags			settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			key	path	string	true	"配置键"	example:"site_name"
//	@Success		204	"配置删除成功"
//	@Failure		401	{object}	response.ErrorResponse	"未授权"
//	@Failure		403	{object}	response.ErrorResponse	"权限不足"
//	@Failure		404	{object}	response.ErrorResponse	"配置不存在"
//	@Failure		500	{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/admin/settings/{key} [delete]
func (h *SettingHandler) DeleteSetting(c *gin.Context) {
	key := c.Param("key")

	// 调用 Use Case Handler
	err := h.deleteHandler.Handle(c.Request.Context(), setting.DeleteCommand{
		Key: key,
	})

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.NoContent(c)
}

// BatchUpdateSettingsRequest 批量更新配置请求
type BatchUpdateSettingsRequest struct {
	Settings []struct {
		Key   string `json:"key" binding:"required"`
		Value any    `json:"value"` // JSONB 原生值
	} `json:"settings" binding:"required,min=1"` // 至少需要一个设置项
}

// BatchUpdateSettings 批量更新配置
//
//	@Summary		批量更新配置
//	@Description	管理员批量更新多个系统配置项的值
//	@Tags			settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		BatchUpdateSettingsRequest	true	"配置列表"
//	@Success		200		{object}	response.EmptyResponse	"批量更新成功"
//	@Failure		400		{object}	response.ErrorResponse		"参数错误"
//	@Failure		401		{object}	response.ErrorResponse		"未授权"
//	@Failure		403		{object}	response.ErrorResponse		"权限不足"
//	@Failure		500		{object}	response.ErrorResponse		"服务器内部错误"
//	@Router			/api/admin/settings/batch [post]
func (h *SettingHandler) BatchUpdateSettings(c *gin.Context) {
	var req BatchUpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 转换为 Command
	items := make([]setting.SettingItemCommand, len(req.Settings))
	for i, s := range req.Settings {
		items[i] = setting.SettingItemCommand{
			Key:   s.Key,
			Value: s.Value,
		}
	}

	// 调用 Use Case Handler
	err := h.batchUpdateHandler.Handle(c.Request.Context(), setting.BatchUpdateCommand{
		Settings: items,
	})

	if err != nil {
		if errors.Is(err, settingDomain.ErrValidationFailed) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}

// =============================================================================
// Category Admin API
// =============================================================================

// GetCategories 获取配置分类列表
//
//	@Summary		配置分类列表
//	@Description	获取所有配置分类，按排序权重升序排列
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.DataResponse[[]setting.CategoryDTO]	"分类列表"
//	@Failure		401	{object}	response.ErrorResponse							"未授权"
//	@Failure		403	{object}	response.ErrorResponse							"权限不足"
//	@Failure		500	{object}	response.ErrorResponse							"服务器内部错误"
//	@Router			/api/admin/settings/categories [get]
func (h *SettingHandler) GetCategories(c *gin.Context) {
	categories, err := h.listCategoriesHandler.Handle(c.Request.Context(), setting.ListCategoriesQuery{})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, categories)
}

// GetCategory 获取单个配置分类
//
//	@Summary		配置分类详情
//	@Description	根据 ID 获取配置分类详情
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int											true	"分类 ID"
//	@Success		200	{object}	response.DataResponse[setting.CategoryDTO]	"分类详情"
//	@Failure		401	{object}	response.ErrorResponse						"未授权"
//	@Failure		403	{object}	response.ErrorResponse						"权限不足"
//	@Failure		404	{object}	response.ErrorResponse						"分类不存在"
//	@Router			/api/admin/settings/categories/{id} [get]
func (h *SettingHandler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, settingDomain.ErrInvalidCategoryID.Error())
		return
	}

	category, err := h.getCategoryHandler.Handle(c.Request.Context(), setting.GetCategoryQuery{
		ID: uint(id),
	})
	if err != nil {
		if errors.Is(err, settingDomain.ErrCategoryNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, category)
}

// CreateCategoryRequest 创建配置分类请求
type CreateCategoryRequest struct {
	Key   string `json:"key" binding:"required" example:"custom"`
	Label string `json:"label" binding:"required" example:"自定义配置"`
	Icon  string `json:"icon" example:"mdi-cog"`
	Order int    `json:"order" example:"100"`
}

// CreateCategory 创建配置分类
//
//	@Summary		创建配置分类
//	@Description	管理员创建新的配置分类
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateCategoryRequest						true	"分类信息"
//	@Success		201		{object}	response.DataResponse[setting.CategoryDTO]	"分类创建成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/admin/settings/categories [post]
func (h *SettingHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.createCategoryHandler.Handle(c.Request.Context(), setting.CreateCategoryCommand{
		Key:   req.Key,
		Label: req.Label,
		Icon:  req.Icon,
		Order: req.Order,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 查询完整的分类信息返回
	category, err := h.getCategoryHandler.Handle(c.Request.Context(), setting.GetCategoryQuery{
		ID: result.ID,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, category)
}

// UpdateCategoryRequest 更新配置分类请求
type UpdateCategoryRequest struct {
	Label string `json:"label" example:"更新后的标签"`
	Icon  string `json:"icon" example:"mdi-settings"`
	Order int    `json:"order" example:"50"`
}

// UpdateCategory 更新配置分类
//
//	@Summary		更新配置分类
//	@Description	管理员更新指定配置分类的信息（Key 不可修改）
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int											true	"分类 ID"
//	@Param			request	body		UpdateCategoryRequest						true	"更新信息"
//	@Success		200		{object}	response.DataResponse[setting.CategoryDTO]	"分类更新成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		404		{object}	response.ErrorResponse						"分类不存在"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/admin/settings/categories/{id} [put]
func (h *SettingHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, settingDomain.ErrInvalidCategoryID.Error())
		return
	}

	var req UpdateCategoryRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		response.BadRequest(c, bindErr.Error())
		return
	}

	category, err := h.updateCategoryHandler.Handle(c.Request.Context(), setting.UpdateCategoryCommand{
		ID:    uint(id),
		Label: req.Label,
		Icon:  req.Icon,
		Order: req.Order,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, category)
}

// DeleteCategory 删除配置分类
//
//	@Summary		删除配置分类
//	@Description	管理员删除指定的配置分类（如有关联配置项则拒绝删除）
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"分类 ID"
//	@Success		204	"分类删除成功"
//	@Failure		400	{object}	response.ErrorResponse	"存在关联配置项"
//	@Failure		401	{object}	response.ErrorResponse	"未授权"
//	@Failure		403	{object}	response.ErrorResponse	"权限不足"
//	@Failure		404	{object}	response.ErrorResponse	"分类不存在"
//	@Failure		500	{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/admin/settings/categories/{id} [delete]
func (h *SettingHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, settingDomain.ErrInvalidCategoryID.Error())
		return
	}

	err = h.deleteCategoryHandler.Handle(c.Request.Context(), setting.DeleteCategoryCommand{
		ID: uint(id),
	})
	if err != nil {
		// 检查是否为关联错误
		if err.Error() != "" && err.Error()[:6] == "cannot" {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.NoContent(c)
}

// =============================================================================
// Public API
// =============================================================================

// PublicSettingsQuery 公开配置查询参数
type PublicSettingsQuery struct {
	Category string `form:"category"`
}

// GetPublicSettings 获取公开配置（无需认证）
//
//	@Summary		公开配置列表
//	@Description	获取公开可见的配置数据（VisibleAt="public"），用于前端展示站点信息等。无需认证。
//	@Tags			public
//	@Accept			json
//	@Produce		json
//	@Param			params	query		handler.PublicSettingsQuery	false	"查询参数"
//	@Success		200		{object}	response.DataResponse[[]setting.PublicSettingItemDTO]	"公开配置列表"
//	@Failure		404		{object}	response.ErrorResponse	"分类不存在"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/public/settings [get]
func (h *SettingHandler) GetPublicSettings(c *gin.Context) {
	var query PublicSettingsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 检查 Handler 是否为 nil
	if h.publicSettingsHandler == nil {
		response.InternalError(c, "Public settings handler not initialized")
		return
	}

	// 调用 PublicSettingsHandler
	result, err := h.publicSettingsHandler.Handle(c.Request.Context(), setting.PublicSettingsQuery{
		CategoryKey: query.Category,
	})
	if err != nil {
		// 检查是否为分类不存在错误
		if query.Category != "" && err.Error() == "category not found: "+query.Category {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}
