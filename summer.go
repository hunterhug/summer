package summer

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/hunterhug/golog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"math"
	"net"
	"net/http"
	_ "net/http/pprof"
)

// RPCMicroService RPC类型的微服务
type RPCMicroService struct {
	Name       string       // 服务名
	GRPCServer *grpc.Server // gRPC实例
}

// RPCMicroServiceRunHooks 微服务勾子，用来注册gRPC微服务
type RPCMicroServiceRunHooks struct {
	// 非空时且是调试模式时使用
	TCPAddress string

	// 注册gRPC服务
	RegisterRPCServer func(gRPCServer *grpc.Server)

	// gRPC服务一元拦截器
	RegisterRPCServerUnaryInterceptor []grpc.UnaryServerInterceptor

	// gRPC服务流拦截器
	RegisterRPCServerStreamInterceptor []grpc.StreamServerInterceptor

	// 注册HTTP代理服务
	RegisterHTTPGateWayServer func(httpServer *runtime.ServeMux, endpoint string, opts []grpc.DialOption)

	// 切面到最后执行的函数
	AopEndFunc func() error

	// 跳过HTTP转化选项
	SkipHttpGatewayOption bool

	// 跳过日志打印
	SkipLogInterceptor bool
}

// NewRPCMicroService 创建一个微服务
func NewRPCMicroService(name string) (*RPCMicroService, error) {
	if name == "" {
		return nil, errors.New("microService name is empty")
	}

	return &RPCMicroService{
		Name: name,
	}, nil
}

// Run 运行微服务
func (s *RPCMicroService) Run(vars *RPCMicroServiceRunHooks) error {
	ctx := context.Background()

	if s == nil {
		log.ErrorContext(ctx, "RPCMicroService is nil")
		return errors.New("RPCMicroService is nil")
	}

	// 非空时且是调试模式时使用
	if vars.TCPAddress != "" && Debug {
		ListeningAddress = vars.TCPAddress
	}

	// TCP监听
	lis, err := net.Listen("tcp", ListeningAddress)

	if err != nil {
		log.ErrorContextWithFields(ctx, map[string]interface{}{
			"address":      ListeningAddress,
			"error string": err.Error(),
		}, "Run(): failed to tcp listen")
		return err
	}

	log.InfoContext(ctx, "Run: tcp listening:%s", ListeningAddress)

	// gRPC服务设置，开启自定义拦截器
	// 一元拦截器
	interceptor := make([]grpc.UnaryServerInterceptor, 0)

	// 先打上追踪ID（一定成功），再打日志（必须在第二个拦截器，不然可能不会执行）
	// 再传上下文（上下文不通过的话会直接返回，所以第二个拦截器打印的日志没有后面拦截器添加的上下文）
	interceptor = append(interceptor, TraceInterceptor)

	if !vars.SkipLogInterceptor {
		interceptor = append(interceptor, LogInterceptor)
	}

	interceptor = append(interceptor, ContextFieldInterceptor)

	for _, o := range vars.RegisterRPCServerUnaryInterceptor {
		interceptor = append(interceptor, o)
	}

	options := make([]grpc.ServerOption, 0)
	if len(interceptor) == 1 {
		options = append(options, grpc.UnaryInterceptor(interceptor[0]))
	} else {
		options = append(options, grpc.ChainUnaryInterceptor(interceptor...))
	}

	// 流式拦截器
	streamInterceptorNum := len(vars.RegisterRPCServerStreamInterceptor)
	if streamInterceptorNum > 0 {
		if streamInterceptorNum == 1 {
			options = append(options, grpc.StreamInterceptor(vars.RegisterRPCServerStreamInterceptor[0]))
		} else {
			options = append(options, grpc.ChainStreamInterceptor(vars.RegisterRPCServerStreamInterceptor...))
		}
	}

	options = append(options, grpc.MaxRecvMsgSize(math.MaxInt32), grpc.MaxSendMsgSize(math.MaxInt32))
	gRPCServer := grpc.NewServer(options...)

	// 健康检查
	//grpc_health_v1.RegisterHealthServer(gRPCServer, health.NewServer())

	//reflection.Register(gRPCServer) // We don't need this feature, see more at https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md

	s.GRPCServer = gRPCServer

	// 注册自己的服务
	vars.RegisterRPCServer(s.GRPCServer)

	// 可开启转发
	if ListeningHTTPGateWay {
		mux := runtime.NewServeMux()

		// HTTP代理生成，用于调试使用
		if !vars.SkipHttpGatewayOption {
			option := runtime.WithMarshalerOption("application/json", &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					AllowPartial:    true,
					UseProtoNames:   true,
					UseEnumNumbers:  true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					AllowPartial:   true,
					DiscardUnknown: true,
				},
			})

			mux = runtime.NewServeMux(option)
		}

		go func() {
			err := http.ListenAndServe(ListeningHTTPGateWayAddress, mux)
			if err != nil {
				log.ErrorContext(ctx, "Run: http listening:%s err:%s", ListeningHTTPGateWayAddress, err.Error())
			}
		}()

		log.InfoContext(ctx, "Run: http listening:%s", ListeningHTTPGateWayAddress)

		// HTTP转gRPC注入
		vars.RegisterHTTPGateWayServer(mux, ListeningAddress, []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)),
		})
	}

	// 可开启性能分析
	if ListeningHTTPProf {
		// pprof分析，固定死端口
		go func() {
			err := http.ListenAndServe(ListeningHTTPProfAddress, nil)
			if err != nil {
				log.ErrorContext(ctx, "Run: http pprof listening:%s err:%s", ListeningHTTPProfAddress, err.Error())
			}
		}()

		log.InfoContext(ctx, "Run: http pprof listening:%s", ListeningHTTPProfAddress)
	}

	// 注入
	err = vars.AopEndFunc()
	if err != nil {
		log.ErrorContext(ctx, "Run: aopEndFunc err:%s", err.Error())
		return err
	}

	// 启动gRPC服务
	err = gRPCServer.Serve(lis)
	if err != nil {
		log.ErrorContext(ctx, "Run: gRPC serve err:%s", err.Error())
		return err
	}

	return nil
}
