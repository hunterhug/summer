package summer

import (
	"context"
	"fmt"
	log "github.com/hunterhug/golog"
	"github.com/hunterhug/summer/stool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"reflect"
)

// LogInterceptor 打印日志拦截器
func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	requestType := reflect.TypeOf(req).String()
	traceRequestLogStr := "Receive rpc request: " + requestType

	log.DebugContextWithFields(ctx, map[string]interface{}{"service.grpc.request": stool.ToJsonString(req)}, traceRequestLogStr)

	resp, err = handler(ctx, req)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	traceResponseLogStr := fmt.Sprintf("Receive rpc response: %s", errMsg)

	log.DebugContextWithFields(ctx, map[string]interface{}{"service.grpc.response": stool.ToJsonString(resp)}, traceResponseLogStr)
	return
}

// TraceInterceptor 追踪ID拦截器
func TraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	traceIdList := make([]string, 0)
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// 追踪ID获取
		traceIdList = md.Get(ContextTraceId)
		// 追踪ID透传
		for _, v := range traceIdList {
			ctx = metadata.AppendToOutgoingContext(ctx, ContextTraceId, v)
		}
	}

	if len(traceIdList) == 0 {
		//  自定义事务ID
		ctx = metadata.AppendToOutgoingContext(ctx, ContextTraceId, stool.GetGUID())
	}

	return handler(ctx, req)
}

// ContextFieldInterceptor 上下文透传处理拦截器
func ContextFieldInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)

	// 批量操作上下文
	for _, field := range ContextFieldList {
		value := make([]string, 0)
		if ok {
			value = md.Get(field.Name)
		}

		// 字段必传，却不存在
		if field.Require {
			if !ok {
				// 拦截器直接终止，不继续下一个拦截器
				log.WarnContext(ctx, "Request context %s required but not found", field.Name)
				return nil, status.Errorf(codes.Unauthenticated, "Request context %s required but not found", field.Name)
			}

			if len(value) == 0 {
				// 拦截器直接终止，不继续下一个拦截器
				log.WarnContext(ctx, "Request context %s required but empty", field.Name)
				return nil, status.Errorf(codes.Unauthenticated, "Request context %s required but empty", field.Name)
			}

			if len(value) == 1 && value[0] == "" {
				// 拦截器直接终止，不继续下一个拦截器
				log.WarnContext(ctx, "Request context %s required but empty", field.Name)
				return nil, status.Errorf(codes.Unauthenticated, "Request context %s required but empty", field.Name)
			}
		}

		// 透传
		if field.Through {
			pairs := make([]string, 0)
			for _, v := range value {
				pairs = append(pairs, field.Name, v)
				// 单值透传
				if field.One {
					break
				}
			}
			if len(pairs) > 0 {
				ctx = SetOutgoingContext(ctx, pairs)
			}
		}
	}

	return handler(ctx, req)
}
