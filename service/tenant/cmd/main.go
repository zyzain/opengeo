package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"

	"opengeo/service/tenant/internal/adapter/mysql"
	"opengeo/service/tenant/internal/application"
	"opengeo/service/tenant/internal/handler"
)

func main() {
	// 从环境变量获取配置
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "root")
	dbName := getEnv("DB_NAME", "opengeo")
	servicePort := getEnv("SERVICE_PORT", "8001")

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database")

	// 设置连接池
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// 创建仓储
	tenantRepo := mysql.NewTenantRepository(db)
	quotaRepo := mysql.NewTenantQuotaRepository(db)

	// 创建服务
	tenantService := application.NewTenantService(tenantRepo, quotaRepo, nil)

	// 创建 Handler
	tenantHandler := handler.NewTenantHandler(tenantService)

	// 启动 HTTP 服务（简化版本，实际应该使用 Hertz 或 gRPC）
	listener, err := net.Listen("tcp", ":"+servicePort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	log.Printf("Tenant service listening on :%s", servicePort)

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down tenant service...")
	_ = tenantHandler // 避免未使用警告
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
