package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"opengeo/gateway/internal/auth"
	"opengeo/gateway/internal/client"
	"opengeo/gateway/internal/content"
	"opengeo/gateway/internal/dal"
	"opengeo/gateway/internal/handler"
	"opengeo/gateway/internal/knowledge"
	"opengeo/gateway/internal/model"
	"opengeo/gateway/internal/router"
	"opengeo/pkg/crypto"
)

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:root@tcp(127.0.0.1:3306)/opengeo?charset=utf8mb4&parseTime=True&loc=Local"
	}

	fmt.Println("正在连接数据库...")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		fmt.Println("请确保 MySQL 已启动，连接信息可通过 MYSQL_DSN 环境变量配置")
		os.Exit(1)
	}
	fmt.Println("数据库连接成功")

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("获取数据库实例失败: %v\n", err)
		os.Exit(1)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 初始化 auth 模块（建表 + 种子数据）
	store := auth.NewStore(db)
	if err := store.AutoMigrate(); err != nil {
		fmt.Printf("数据库迁移失败: %v\n", err)
		os.Exit(1)
	}

	// 种子数据：默认租户
	defaultTenant, err := store.SeedTenant(context.Background())
	if err != nil {
		fmt.Printf("创建默认租户失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("默认租户检查完成 (ID: %d)\n", defaultTenant.ID)

	// 种子数据：权限
	if err := store.SeedPermissions(context.Background()); err != nil {
		fmt.Printf("创建默认权限失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("默认权限检查完成")

	// 种子数据：角色
	if err := store.SeedRoles(context.Background(), defaultTenant.ID); err != nil {
		fmt.Printf("创建默认角色失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("默认角色检查完成")

	// 种子数据：角色-权限关联
	if err := store.SeedRolePermissions(context.Background()); err != nil {
		fmt.Printf("创建角色权限关联失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("角色权限关联检查完成")

	// 种子数据：管理员（密码从环境变量读取，默认值仅用于开发环境）
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		if os.Getenv("GO_ENV") == "production" {
			fmt.Println("错误: 生产环境必须设置 ADMIN_PASSWORD 环境变量")
			os.Exit(1)
		}
		adminPassword = "Admin@123456"
	}
	hashedPassword, err := crypto.HashPassword(adminPassword)
	if err != nil {
		fmt.Printf("生成管理员密码失败: %v\n", err)
		os.Exit(1)
	}
	if err := store.SeedAdmin(context.Background(), defaultTenant.ID, hashedPassword); err != nil {
		fmt.Printf("创建管理员账号失败: %v\n", err)
		os.Exit(1)
	}
	// 分配 admin 角色给管理员
	adminUser, _ := store.GetUserByUsername(context.Background(), "admin")
	if adminUser != nil {
		if err := store.SeedAdminRoleAssignment(context.Background(), adminUser.ID); err != nil {
			fmt.Printf("分配管理员角色失败: %v\n", err)
		}
	}
	fmt.Println("管理员账号检查完成")

	authSvc := auth.NewService(store)
	userClient := client.NewUserClient(authSvc)

	contentStore := content.NewStore(db)
	if err := contentStore.AutoMigrate(); err != nil {
		fmt.Printf("内容表迁移失败: %v\n", err)
		os.Exit(1)
	}
	contentSvc := content.NewService(contentStore)
	contentClient := client.NewContentClient(contentSvc)
	knowledgeStore := knowledge.NewStore(db)
	if err := knowledgeStore.AutoMigrate(); err != nil {
		fmt.Printf("知识图谱表迁移失败: %v\n", err)
		os.Exit(1)
	}
	if err := knowledgeStore.Seed(context.Background(), 1); err != nil {
		fmt.Printf("知识图谱种子数据创建失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("知识图谱种子数据检查完成")
	knowledgeSvc := knowledge.NewService(knowledgeStore)
	knowledgeClient := client.NewKnowledgeClient(knowledgeSvc)

	publishClient := client.NewPublishClient(db)
	accountClient := client.NewAccountClient(db)
	scheduleClient := client.NewScheduleClient(db)
	monitorClient := client.NewMonitorClient(db)
	systemClient := client.NewSystemClient(db)

	// 初始化DAL层
	fpRepo := dal.NewBrowserFingerprintRepository(db)
	proxyRepo := dal.NewProxyIPRepository(db)
	tplRepo := dal.NewContentTemplateRepository(db)
	staggerStrategyRepo := dal.NewStaggerStrategyRepository(db)
	staggerConfigRepo := dal.NewStaggerConfigRepository(db)
	brandRepo := dal.NewBrandRepository(db)

	// 自动迁移DAL表
	if err := db.AutoMigrate(
		&model.BrowserFingerprint{},
		&model.ProxyIP{},
		&model.ContentTemplate{},
		&model.StaggerStrategy{},
		&model.StaggerConfig{},
		&model.Brand{},
		&model.BrandMetadata{},
		&model.GlossaryEntry{},
		&model.BrandSnapshot{},
		&model.KnowledgeEntity{},
		&model.KnowledgeRelation{},
	); err != nil {
		fmt.Printf("DAL表迁移失败: %v\n", err)
		os.Exit(1)
	}

	h := handler.NewHandler(
		userClient,
		contentClient,
		knowledgeClient,
		publishClient,
		accountClient,
		scheduleClient,
		monitorClient,
		systemClient,
		fpRepo,
		proxyRepo,
		tplRepo,
		staggerStrategyRepo,
		staggerConfigRepo,
		brandRepo,
	)

	svr := server.Default(
		server.WithHostPorts(":8080"),
		server.WithExitWaitTime(time.Second*5),
	)

	router.RegisterRoutes(svr, h, store)

	go func() {
		if err := svr.Run(); err != nil {
			fmt.Printf("Gateway服务启动失败: %v\n", err)
			os.Exit(1)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("========================================")
	fmt.Println("  OpenGEO Gateway 服务已启动")
	fmt.Printf("  监听端口: %s\n", port)
	fmt.Println("  多租户 + RBAC 权限系统已就绪")
	fmt.Println("========================================")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("正在关闭Gateway服务...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	svr.Shutdown(shutdownCtx)

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}
	fmt.Println("Gateway服务已关闭")
}
