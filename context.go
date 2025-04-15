package summer

import (
	"context"
)

var (
	// ContextAuthorizationValuePrefix 授权令牌值前缀
	ContextAuthorizationValuePrefix = "bearer"

	// ContextTraceId 日志追踪ID
	ContextTraceId = "trace.id"

	// ContextAppId 应用ID
	ContextAppId = "app.id"

	// ContextUserId 应用用户ID
	ContextUserId = "app.session.user.id"

	// ContextAppAdmin 管理员标记
	ContextAppAdmin = "app.session.user.admin"

	// ContextAppAdminYes 是否管理员
	ContextAppAdminYes = "admin"
)

// ContextFieldRule 上下文字段传输规则
type ContextFieldRule struct {
	// 字段名
	Name string
	// 是否必有值
	Require bool
	// 单值
	One bool
	// 透传
	Through bool
}

var (
	// ContextFieldList 需要额外处理的上下文
	ContextFieldList = []ContextFieldRule{
		{
			Name:    ContextAppId,
			Require: true,
			One:     true,
			Through: true,
		},
		{
			Name:    ContextAppAdmin,
			Require: false,
			One:     true,
			Through: true,
		},
		{
			Name:    ContextUserId,
			Require: false,
			One:     true,
			Through: true,
		},
	}
)

func ContextGetAppID(ctx context.Context) string {
	return GetOutgoingContextOne(ctx, ContextAppId)
}

func ContextGetUserId(ctx context.Context) string {
	return GetOutgoingContextOne(ctx, ContextUserId)
}

func ContextIsAdmin(ctx context.Context) bool {
	return GetOutgoingContextOne(ctx, ContextAppAdmin) == ContextAppAdminYes
}
