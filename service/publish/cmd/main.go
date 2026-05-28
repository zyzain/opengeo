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

	// 创建Kitex服务器
	svr := server.NewServer(
		server.WithServiceAddr(&net.TCPAddr{Port: 8890}),
		server.WithRegistry(r),
	)

	// 注册服务实现
	// TODO: 注册PublishService实现

	// 启动服务器
	go func() {
		fmt.Println("Publish Service启动中，监听端口: 8890")
		if err := svr.Run(); err != nil {
			fmt.Printf("Publish Service启动失败: %v\n", err)
			os.Exit(1)
		}
	}()

	// 等待Consul注册
	time.Sleep(time.Second * 2)
	fmt.Println("Publish Service已启动并注册到Consul")
	fmt.Println("服务地址: localhost:8890")
	fmt.Println("Metrics地址: localhost:9093")

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("正在关闭Publish Service...")
	_, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := svr.Stop(); err != nil {
		fmt.Printf("Publish Service关闭失败: %v\n", err)
	}
	fmt.Println("Publish Service已关闭")
}

func initTracer() (*sdktrace.TracerProvider, error) {
	exp, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:4317"),
	)
	if err != nil {
		return nil, fmt.Errorf("创建OTLP导出器失败: %w", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("publish-service"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("创建资源失败: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}