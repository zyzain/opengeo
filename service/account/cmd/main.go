package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-consul"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"opengeo/service/account/internal/dal"
	"opengeo/service/account/internal/handler"
	"opengeo/service/account/internal/service"
)

func main() {
	// 初始化OpenTelemetry
	tp, err := initTracer()
	if err != nil {
		fmt.Printf("初始化Tracer失败: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			fmt.Printf("关闭Tracer失败: %v\n", err)
		}
	}()

	// 创建Consul注册器
	r, err := consul.NewConsulRegister("localhost:8500")
	if err != nil {
		fmt.Printf("连接Consul失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库
	if err := dal.Init(nil); err != nil {
		fmt.Printf("数据库初始化失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("数据库连接成功")

	// 种子默认租户
	defaultTenant := dal.SeedTenant(context.Background(), dal.DB)
	fmt.Printf("默认租户 (ID: %d)\n", defaultTenant.ID)

	// 初始化仓储
	userRepo := dal.NewUserRepository(dal.DB)
	accountRepo := dal.NewAccountRepository(dal.DB)
	_ = dal.NewAccountGroupRepository(dal.DB)
	tenantRepo := dal.NewTenantRepository(dal.DB)
	roleRepo := dal.NewRoleRepository(dal.DB)
	permRepo := dal.NewPermissionRepository(dal.DB)
	userRoleRepo := dal.NewUserRoleRepository(dal.DB)
	rolePermRepo := dal.NewRolePermissionRepository(dal.DB)

	// 初始化服务
	userSvc := service.NewUserService(userRepo)
	accountSvc := service.NewAccountService(accountRepo)
	rbacSvc := service.NewRBACService(tenantRepo, roleRepo, permRepo, userRoleRepo, rolePermRepo)

	// 初始化处理器
	accountHandler := handler.NewAccountHandler(accountSvc, userSvc, rbacSvc)

	// 创建Kitex服务器
	svr := server.NewServer(
		server.WithServiceAddr(&net.TCPAddr{Port: 8888}),
		server.WithRegistry(r),
	)

	// 注册服务实现
	// TODO: 使用生成的Kitex代码注册服务
	// accountSvcImpl := &AccountServiceImpl{handler: accountHandler}
	// svr.RegisterService(accountSvcImpl.ServiceName(), accountSvcImpl)

	_ = accountHandler

	// 启动服务器
	go func() {
		fmt.Println("Account Service启动中，监听端口: 8888")
		if err := svr.Run(); err != nil {
			fmt.Printf("Account Service启动失败: %v\n", err)
			os.Exit(1)
		}
	}()

	// 等待Consul注册
	time.Sleep(time.Second * 2)
	fmt.Println("Account Service已启动并注册到Consul")
	fmt.Println("服务地址: localhost:8888")
	fmt.Println("Metrics地址: localhost:9091")

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("正在关闭Account Service...")
	_, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := svr.Stop(); err != nil {
		fmt.Printf("Account Service关闭失败: %v\n", err)
	}
	fmt.Println("Account Service已关闭")
}

func initTracer() (*sdktrace.TracerProvider, error) {
	// 创建OTLP导出器
	exp, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:4317"),
	)
	if err != nil {
		return nil, fmt.Errorf("创建OTLP导出器失败: %w", err)
	}

	// 创建资源
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("account-service"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("创建资源失败: %w", err)
	}

	// 创建TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// 设置全局TracerProvider
	otel.SetTracerProvider(tp)

	return tp, nil
}

// AccountServiceImpl 实现AccountService接口
// 这是一个占位符实现，需要根据生成的代码进行完善
type AccountServiceImpl struct{}

// TODO: 实现AccountService的所有方法
// func (s *AccountServiceImpl) CreateUser(ctx context.Context, req *account.CreateUserRequest) (*account.User, error) {
//     // 实现创建用户逻辑
//     return nil, fmt.Errorf("未实现")
// }

// func (s *AccountServiceImpl) GetUser(ctx context.Context, req *account.GetUserRequest) (*account.User, error) {
//     // 实现获取用户逻辑
//     return nil, fmt.Errorf("未实现")
// }

// ... 其他方法实现