# OpenGEO BrandOS Makefile (云端 API + SaaS 架构)

# 变量定义
GO := go
KITEX := kitex
HZ := hz
DOCKER := docker
DOCKER_COMPOSE := docker-compose
BUF := buf

# 默认目标
.PHONY: all
all: help

# 帮助信息
.PHONY: help
help:
	@echo "OpenGEO BrandOS - 开发命令 (云端 API + SaaS 架构)"
	@echo ""
	@echo "架构说明："
	@echo "  gateway:   标准分层 (router/handler/middleware/client)"
	@echo "  tenant:    六边形架构 (domain/application/port/adapter)"
	@echo "  brand:     六边形架构 (domain/application/port/adapter)"
	@echo "  content:   六边形架构 (domain/application/port/adapter)"
	@echo "  publish:   六边形 + EDA (事件驱动)"
	@echo "  monitor:   CQRS读写分离 (command/query)"
	@echo ""
	@echo "Proto 命令："
	@echo "  make proto-gen      生成 Go/TS 代码"
	@echo "  make proto-lint     Proto 规范检查"
	@echo "  make proto-breaking Breaking Change 检测"
	@echo ""
	@echo "基础设施命令："
	@echo "  make dev-up          启动基础设施"
	@echo "  make dev-down        停止基础设施"
	@echo "  make dev             一键启动全栈（基础设施+后端+前端）"
	@echo ""
	@echo "构建命令："
	@echo "  make build           构建所有服务"
	@echo "  make build-gateway   构建Gateway服务"
	@echo "  make build-tenant    构建Tenant服务"
	@echo "  make build-brand     构建Brand服务"
	@echo "  make build-content   构建Content服务"
	@echo "  make build-publish   构建Publish服务"
	@echo "  make build-monitor   构建Monitor服务"
	@echo ""
	@echo "开发运行命令："
	@echo "  make dev-gateway     本地运行Gateway"
	@echo "  make dev-tenant      本地运行Tenant"
	@echo "  make dev-brand       本地运行Brand"
	@echo "  make dev-content     本地运行Content"
	@echo "  make dev-publish     本地运行Publish"
	@echo "  make dev-monitor     本地运行Monitor"
	@echo ""
	@echo "Docker命令："
	@echo "  make docker          Docker全栈启动"
	@echo "  make docker-down     停止Docker"
	@echo "  make docker-restart  重启Docker容器"
	@echo ""
	@echo "测试命令："
	@echo "  make test            运行所有测试"
	@echo "  make test-unit       运行单元测试"
	@echo "  make test-integration 运行集成测试"
	@echo "  make test-oss-only   验证开源版完整性（无企业插件）"

# ==================== Proto 命令 ====================

.PHONY: proto-gen
proto-gen:
	@echo "生成 Proto 代码..."
	cd proto && $(BUF) generate
	@echo "Proto 代码生成完成！"

.PHONY: proto-lint
proto-lint:
	@echo "Proto 规范检查..."
	cd proto && $(BUF) lint
	@echo "Proto 规范检查完成！"

.PHONY: proto-breaking
proto-breaking:
	@echo "Breaking Change 检测..."
	cd proto && $(BUF) breaking --against '.git#branch=main'
	@echo "Breaking Change 检测完成！"

.PHONY: proto
proto: proto-gen proto-lint
	@echo "Proto 处理完成！"

# ==================== 基础设施命令 ====================

.PHONY: dev-up
dev-up:
	@echo "启动基础设施..."
	$(DOCKER_COMPOSE) up -d mysql redis consul jaeger prometheus grafana
	@echo "等待基础设施就绪..."
	@sleep 10
	@echo "基础设施启动完成！"

.PHONY: dev-down
dev-down:
	@echo "停止基础设施..."
	$(DOCKER_COMPOSE) down

.PHONY: dev
dev: dev-up
	@echo "启动后端服务..."
	$(GO) run ./gateway/cmd/ &
	$(GO) run ./service/tenant/cmd/ &
	$(GO) run ./service/brand/cmd/ &
	$(GO) run ./service/content/cmd/ &
	$(GO) run ./service/publish/cmd/ &
	$(GO) run ./service/monitor/cmd/ &
	$(GO) run ./service/scheduler/cmd/ &
	$(GO) run ./service/account/cmd/ &
	$(GO) run ./service/system/cmd/ &
	@echo "启动前端..."
	cd web && npm run dev &

# ==================== 构建命令 ====================

.PHONY: build
build: build-gateway build-tenant build-brand build-content build-publish build-monitor build-scheduler build-account build-system
	@echo "所有服务构建完成！"

.PHONY: build-gateway
build-gateway:
	@echo "构建Gateway服务..."
	$(GO) build -o bin/gateway ./gateway/cmd/

.PHONY: build-tenant
build-tenant:
	@echo "构建Tenant服务..."
	$(GO) build -o bin/tenant-service ./service/tenant/cmd/

.PHONY: build-brand
build-brand:
	@echo "构建Brand服务..."
	$(GO) build -o bin/brand-service ./service/brand/cmd/

.PHONY: build-account
build-account:
	@echo "构建Account服务..."
	$(GO) build -o bin/account-service ./service/account/cmd/

.PHONY: build-content
build-content:
	@echo "构建Content服务..."
	$(GO) build -o bin/content-service ./service/content/cmd/

.PHONY: build-publish
build-publish:
	@echo "构建Publish服务..."
	$(GO) build -o bin/publish-service ./service/publish/cmd/

.PHONY: build-monitor
build-monitor:
	@echo "构建Monitor服务..."
	$(GO) build -o bin/monitor-service ./service/monitor/cmd/

.PHONY: build-scheduler
build-scheduler:
	@echo "构建Scheduler服务..."
	$(GO) build -o bin/scheduler-service ./service/scheduler/cmd/

.PHONY: build-system
build-system:
	@echo "构建System服务..."
	$(GO) build -o bin/system-service ./service/system/cmd/

# ==================== 开发运行命令 ====================

.PHONY: dev-gateway
dev-gateway:
	@echo "启动Gateway服务..."
	$(GO) run ./gateway/cmd/

.PHONY: dev-account
dev-account:
	@echo "启动Account服务..."
	$(GO) run ./service/account/cmd/

.PHONY: dev-content
dev-content:
	@echo "启动Content服务..."
	$(GO) run ./service/content/cmd/

.PHONY: dev-publish
dev-publish:
	@echo "启动Publish服务..."
	$(GO) run ./service/publish/cmd/

.PHONY: dev-monitor
dev-monitor:
	@echo "启动Monitor服务..."
	$(GO) run ./service/monitor/cmd/

.PHONY: dev-scheduler
dev-scheduler:
	@echo "启动Scheduler服务..."
	$(GO) run ./service/scheduler/cmd/

.PHONY: dev-system
dev-system:
	@echo "启动System服务..."
	$(GO) run ./service/system/cmd/

# ==================== Docker命令 ====================

.PHONY: docker
docker:
	@echo "Docker全栈启动..."
	$(DOCKER_COMPOSE) up -d

.PHONY: docker-down
docker-down:
	@echo "停止Docker..."
	$(DOCKER_COMPOSE) down

.PHONY: docker-restart
docker-restart:
	@echo "重启Docker容器..."
	$(DOCKER_COMPOSE) restart

# ==================== 测试命令 ====================

.PHONY: test
test: test-unit test-integration
	@echo "所有测试完成！"

.PHONY: test-unit
test-unit:
	@echo "运行单元测试..."
	$(GO) test ./... -short -v

.PHONY: test-integration
test-integration:
	@echo "运行集成测试..."
	$(GO) test ./... -run Integration -v

.PHONY: test-oss-only
test-oss-only:
	@echo "验证开源版完整性..."
	$(GO) build ./...
	$(GO) test ./... -short
	@echo "开源版验证完成！"

# ==================== Mock 生成 ====================

.PHONY: mock-gen
mock-gen:
	@echo "生成 MSW Mock..."
	./scripts/mock-gen.sh
	@echo "Mock 生成完成！"

# ==================== 代码质量 ====================

.PHONY: lint
lint:
	@echo "运行代码检查..."
	golangci-lint run ./...

.PHONY: fmt
fmt:
	@echo "格式化代码..."
	gofmt -w .

.PHONY: vet
vet:
	@echo "运行go vet..."
	$(GO) vet ./...

# ==================== 数据库命令 ====================

.PHONY: db-init
db-init:
	@echo "初始化数据库..."
	mysql -h 127.0.0.1 -P 3306 -u root -proot < scripts/init.sql

.PHONY: db-reset
db-reset:
	@echo "重置数据库..."
	$(MAKE) db-init

# ==================== 清理命令 ====================

.PHONY: clean
clean:
	@echo "清理构建产物..."
	rm -rf bin/
	rm -rf tmp/

.PHONY: tidy
tidy:
	@echo "整理依赖..."
	$(GO) mod tidy