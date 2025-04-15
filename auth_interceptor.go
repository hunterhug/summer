package summer

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	log "github.com/hunterhug/golog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

var (
	// AuthAppIdMap 允许的AppId列表，外部需要定时从管理服务加载
	AuthAppIdMap = sync.Map{}

	// AuthAPIMap 接口列表
	AuthAPIMap = map[string]AuthAPI{}

	// AuthServiceInterface 授权服务
	AuthServiceInterface AuthSessionInterface
)

type AuthAPI struct {
	// API名称
	Name string
	// 是否需要认证
	NeedToken bool
	// 是否是管理员接口，管理员接口默认需要授权
	IsAdmin bool
	// 对应的方法，没有任何作用，方便跳转
	method interface{}
}

func NewAuthAPI(name string, needToken bool, isAdmin bool, method interface{}) AuthAPI {
	if isAdmin {
		needToken = true
	}
	return AuthAPI{
		Name:      name,
		NeedToken: needToken,
		IsAdmin:   isAdmin,
		method:    method,
	}
}

// AuthSessionInterface Session暴露接口
type AuthSessionInterface interface {
	// GetSessionInfoByAccessToken 获取某令牌对信息，先走JWT无状态解码数据，force则强制查询服务端
	GetSessionInfoByAccessToken(accessToken string, force bool) (tokenData SessionTokenData, err error)
}

// SessionTokenData 授权Token数据
type SessionTokenData struct {
	// 用户ID
	UserId string `json:"user_id,omitempty"`
	// 创建时间，毫秒
	CreateMSTime int64 `json:"create_ms_time,omitempty"`
	// 过期时间，毫秒
	ExpiryMSTime int64 `json:"expiry_ms_time,omitempty"`
	// 是否已过期
	IsExpire bool `json:"is_expire,omitempty"`
	// 客户端Payload数据，返回给客户端的数据，客户端解密可看
	ClientPayload map[string]interface{} `json:"client_payload,omitempty"`
	// 服务端Session数据，保存在服务器端的数据，客户端不可看
	ServerSessionData map[string]interface{} `json:"server_session_data,omitempty"`
}

func AuthGRPCInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// AppId要符合，过滤非法应用
	appId := GetIncomingContextOne(ctx, ContextAppId)
	_, ok := AuthAppIdMap.Load(appId)
	if !ok {
		log.WarnContext(ctx, "Request unauthenticated with %s not allow=%s", ContextAppId, appId)
		return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated with %s not allow=%s", ContextAppId, appId)
	}

	// 进行认证和授权
	method := info.FullMethod
	auth, ok := AuthAPIMap[method]

	// 打个日志看看找不找得到路由
	log.DebugContextWithFields(ctx, map[string]interface{}{
		"service.gateway_auth": map[string]interface{}{
			"auth_method": method,
			"auth_found":  ok,
			"auth_api":    auth,
		},
	}, "AuthGRPCInterceptor start")

	// 不存在的路由不允许
	if !ok {
		log.WarnContext(ctx, "Request unauthenticated with msg=method %s not allow", method)
		return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated with msg=method %s not allow", method)
	}

	sessionData := make(map[string]interface{}, 0)
	//payloadData := make(map[string]interface{}, 0)

	// 需要认证的
	if auth.NeedToken {
		// authorization: bearer a.b.c
		accessToken, err := grpc_auth.AuthFromMD(ctx, ContextAuthorizationValuePrefix)
		if err != nil {
			log.WarnContext(ctx, "Request AuthFromMD err msg=%s", err.Error())
			return nil, err
		}

		// 强制校验令牌
		tokenData, err := AuthServiceInterface.GetSessionInfoByAccessToken(accessToken, true)
		if err != nil {
			log.WarnContext(ctx, "Request unauthenticated with msg=%s", err.Error())
			return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated with msg=%s", err.Error())
		}

		// 打印用户信息出来看看
		log.DebugContextWithFields(ctx, map[string]interface{}{
			"service.gateway_auth_access_token": map[string]interface{}{
				"auth_access_token":      accessToken,
				"auth_access_token_data": tokenData,
				"auth_method":            method,
				"auth_api":               auth,
			},
		}, "AuthGRPCInterceptor auth accessToken")

		sessionData = tokenData.ServerSessionData

		// 令牌不属于当前的应用
		if value, ok := sessionData[ContextAppId].(string); !ok || value != appId {
			log.WarnContext(ctx, "Request unauthenticated with msg=token is not match %s=%s", ContextAppId, appId)
			return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated with msg=token is not match %s=%s", ContextAppId, appId)
		}

		// 透传用户ID
		ctx = SetOutgoingContext(ctx, []string{ContextUserId, tokenData.UserId})
		//payloadData = tokenData.ClientPayload
	}

	// 是否session中有管理员标志,透传
	isAdminValue, ok := sessionData[ContextAppAdmin].(string)
	if ok && isAdminValue == ContextAppAdminYes {
		// 管理员标志透传
		ctx = SetOutgoingContext(ctx, []string{ContextAppAdmin, ContextAppAdminYes})
	}

	// 如果需要管理员权限却没有
	if auth.IsAdmin && !(isAdminValue == ContextAppAdminYes) {
		// 拒绝
		log.WarnContext(ctx, "Request unauthenticated with msg=method %s should be admin", method)
		return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated with msg=%s should be admin", method)
	}

	return handler(ctx, req)
}
