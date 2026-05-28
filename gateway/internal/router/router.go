package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"opengeo/gateway/internal/handler"
	"opengeo/gateway/internal/middleware"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(s *server.Hertz, hd *handler.Handler, permChecker middleware.PermissionChecker) {
	// 注册全局中间件
	s.Use(middleware.CORS())
	s.Use(middleware.RequestID())
	s.Use(middleware.Logger())
	s.Use(middleware.RateLimiter())
	s.Use(middleware.Recovery())

	// 健康检查
	s.GET("/health", hd.Health)
	s.GET("/ready", hd.Ready)

	// API v1路由组
	v1 := s.Group("/api/v1")
	{
		// 认证相关（不需要JWT）
		auth := v1.Group("/auth")
		{
			auth.POST("/login", hd.Login)
			auth.POST("/register", hd.Register)
			auth.POST("/refresh", hd.RefreshToken)
		}

		// 需要认证的路由
		protected := v1.Group("")
		protected.Use(middleware.JWT())
		{
			// 用户管理（需要 user 权限）
			userGroup := protected.Group("/users")
			userGroup.Use(middleware.RequirePermission(permChecker, "user", "read"))
			{
				userGroup.GET("/:id", hd.GetUser)
				userGroup.GET("", hd.ListUsers)
				userGroup.GET("/:user_id/roles", hd.GetUserRoles)
			}
			userWrite := protected.Group("/users")
			userWrite.Use(middleware.RequirePermission(permChecker, "user", "update"))
			{
				userWrite.PUT("/:id", hd.UpdateUser)
				userWrite.DELETE("/:id", hd.DeleteUser)
			}

			// 角色管理（需要 role 权限）
			roleGroup := protected.Group("/roles")
			roleGroup.Use(middleware.RequirePermission(permChecker, "role", "read"))
			{
				roleGroup.GET("/:id", hd.GetRole)
				roleGroup.GET("", hd.ListRoles)
				roleGroup.GET("/:role_id/permissions", hd.GetRolePermissions)
			}
			roleWrite := protected.Group("/roles")
			roleWrite.Use(middleware.RequirePermission(permChecker, "role", "create"))
			{
				roleWrite.POST("", hd.CreateRole)
				roleWrite.POST("/assign", hd.AssignRole)
				roleWrite.POST("/revoke", hd.RevokeRole)
				roleWrite.POST("/permissions", hd.AddRolePermission)
			}
			roleDelete := protected.Group("/roles")
			roleDelete.Use(middleware.RequirePermission(permChecker, "role", "delete"))
			{
				roleDelete.PUT("/:id", hd.UpdateRole)
				roleDelete.DELETE("/:id", hd.DeleteRole)
				roleDelete.DELETE("/permissions", hd.RemoveRolePermission)
			}

			// 权限管理（需要 role:read 权限）
			permGroup := protected.Group("/permissions")
			permGroup.Use(middleware.RequirePermission(permChecker, "role", "read"))
			{
				permGroup.GET("", hd.ListPermissions)
				permGroup.POST("/check", hd.CheckPermission)
			}

			// 租户管理（需要 tenant 权限）
			tenantGroup := protected.Group("/tenants")
			tenantGroup.Use(middleware.RequirePermission(permChecker, "tenant", "read"))
			{
				tenantGroup.GET("/:id", hd.GetTenant)
				tenantGroup.GET("", hd.ListTenants)
			}
			tenantWrite := protected.Group("/tenants")
			tenantWrite.Use(middleware.RequirePermission(permChecker, "tenant", "update"))
			{
				tenantWrite.POST("", hd.CreateTenant)
				tenantWrite.PUT("/:id", hd.UpdateTenant)
				tenantWrite.DELETE("/:id", hd.DeleteTenant)
			}

			// 内容管理（需要 content 权限）
			contentGroup := protected.Group("/contents")
			contentGroup.Use(middleware.RequirePermission(permChecker, "content", "read"))
			{
				contentGroup.GET("/:id", hd.GetContent)
				contentGroup.GET("", hd.ListContents)
			}
			contentWrite := protected.Group("/contents")
			contentWrite.Use(middleware.RequirePermission(permChecker, "content", "create"))
			{
				contentWrite.POST("", hd.CreateContent)
				contentWrite.POST("/:id/optimize", hd.OptimizeContent)
				contentWrite.POST("/:id/compliance", hd.CheckCompliance)
			}
			contentUpdate := protected.Group("/contents")
			contentUpdate.Use(middleware.RequirePermission(permChecker, "content", "update"))
			{
				contentUpdate.PUT("/:id", hd.UpdateContent)
			}
			contentDelete := protected.Group("/contents")
			contentDelete.Use(middleware.RequirePermission(permChecker, "content", "delete"))
			{
				contentDelete.DELETE("/:id", hd.DeleteContent)
			}
			contentPublish := protected.Group("/contents")
			contentPublish.Use(middleware.RequirePermission(permChecker, "publish", "execute"))
			{
				contentPublish.POST("/:id/publish", hd.PublishContent)
			}

			// 账号管理
			accountGroup := protected.Group("/accounts")
			accountGroup.Use(middleware.RequirePermission(permChecker, "content", "read"))
			{
				accountGroup.GET("/:id", hd.GetAccount)
				accountGroup.GET("", hd.ListAccounts)
				accountGroup.GET("/:id/health", hd.GetAccountHealth)
			}
			accountWrite := protected.Group("/accounts")
			accountWrite.Use(middleware.RequirePermission(permChecker, "content", "create"))
			{
				accountWrite.POST("", hd.CreateAccount)
				accountWrite.PUT("/:id", hd.UpdateAccount)
				accountWrite.DELETE("/:id", hd.DeleteAccount)
			}

			// 账号分组管理
			accountGroupGroup := protected.Group("/account-groups")
			accountGroupGroup.Use(middleware.RequirePermission(permChecker, "content", "read"))
			{
				accountGroupGroup.GET("/:id", hd.GetAccountGroup)
				accountGroupGroup.GET("", hd.ListAccountGroups)
			}
			accountGroupWrite := protected.Group("/account-groups")
			accountGroupWrite.Use(middleware.RequirePermission(permChecker, "content", "create"))
			{
				accountGroupWrite.POST("", hd.CreateAccountGroup)
				accountGroupWrite.PUT("/:id", hd.UpdateAccountGroup)
				accountGroupWrite.DELETE("/:id", hd.DeleteAccountGroup)
				accountGroupWrite.POST("/:id/accounts", hd.AddAccountToGroup)
				accountGroupWrite.DELETE("/:id/accounts/:account_id", hd.RemoveAccountFromGroup)
			}

			// 指纹管理
			fingerprintGroup := protected.Group("/fingerprints")
			fingerprintGroup.Use(middleware.RequirePermission(permChecker, "content", "read"))
			{
				fingerprintGroup.GET("", hd.ListFingerprints)
			}
			fingerprintWrite := protected.Group("/fingerprints")
			fingerprintWrite.Use(middleware.RequirePermission(permChecker, "content", "create"))
			{
				fingerprintWrite.POST("", hd.CreateFingerprint)
				fingerprintWrite.PUT("/:id", hd.UpdateFingerprint)
				fingerprintWrite.DELETE("/:id", hd.DeleteFingerprint)
				fingerprintWrite.POST("/:id/toggle", hd.ToggleFingerprint)
			}

			// 代理IP管理
			proxyGroup := protected.Group("/proxies")
			proxyGroup.Use(middleware.RequirePermission(permChecker, "content", "read"))
			{
				proxyGroup.GET("", hd.ListProxies)
			}
			proxyWrite := protected.Group("/proxies")
			proxyWrite.Use(middleware.RequirePermission(permChecker, "content", "create"))
			{
				proxyWrite.POST("", hd.CreateProxy)
				proxyWrite.DELETE("/:id", hd.DeleteProxy)
				proxyWrite.POST("/:id/check", hd.CheckProxy)
			}

			// 平台管理
			platformGroup := protected.Group("/platforms")
			platformGroup.Use(middleware.RequirePermission(permChecker, "publish", "read"))
			{
				platformGroup.GET("", hd.ListPlatforms)
				platformGroup.GET("/:id", hd.GetPlatform)
			}
			platformWrite := protected.Group("/platforms")
			platformWrite.Use(middleware.RequirePermission(permChecker, "publish", "create"))
			{
				platformWrite.POST("", hd.CreatePlatform)
				platformWrite.PUT("/:id", hd.UpdatePlatform)
				platformWrite.DELETE("/:id", hd.DeletePlatform)
				platformWrite.POST("/:id/enable", hd.EnablePlatform)
				platformWrite.POST("/:id/disable", hd.DisablePlatform)
			}

			// 渠道管理
			channelGroup := protected.Group("/channels")
			channelGroup.Use(middleware.RequirePermission(permChecker, "publish", "read"))
			{
				channelGroup.GET("/:id", hd.GetChannel)
				channelGroup.GET("", hd.ListChannels)
				channelGroup.GET("/platforms", hd.GetChannelPlatforms)
			}
			channelWrite := protected.Group("/channels")
			channelWrite.Use(middleware.RequirePermission(permChecker, "publish", "create"))
			{
				channelWrite.POST("", hd.CreateChannel)
			}

			// 发布管理
			publishGroup := protected.Group("/publish")
			publishGroup.Use(middleware.RequirePermission(permChecker, "publish", "read"))
			{
				publishGroup.GET("/tasks/:id", hd.GetPublishTask)
				publishGroup.GET("/tasks", hd.ListPublishTasks)
				publishGroup.POST("/preview", hd.PreviewPublish)
				publishGroup.POST("/validate", hd.ValidatePublish)
				publishGroup.POST("/dedup/check", hd.CheckDedup)
			}
			publishWrite := protected.Group("/publish")
			publishWrite.Use(middleware.RequirePermission(permChecker, "publish", "execute"))
			{
				publishWrite.POST("/tasks", hd.CreatePublishTask)
				publishWrite.POST("/tasks/:id/cancel", hd.CancelPublishTask)
				publishWrite.POST("/tasks/:id/retry", hd.RetryPublishTask)
			}

			// 错峰策略管理
			staggerGroup := protected.Group("/stagger")
			staggerGroup.Use(middleware.RequirePermission(permChecker, "publish", "read"))
			{
				staggerGroup.GET("/strategies", hd.ListStaggerStrategies)
				staggerGroup.GET("/config", hd.GetStaggerConfig)
			}
			staggerWrite := protected.Group("/stagger")
			staggerWrite.Use(middleware.RequirePermission(permChecker, "publish", "execute"))
			{
				staggerWrite.POST("/strategies", hd.CreateStaggerStrategy)
				staggerWrite.PUT("/strategies/:id", hd.UpdateStaggerStrategy)
				staggerWrite.DELETE("/strategies/:id", hd.DeleteStaggerStrategy)
				staggerWrite.POST("/strategies/:id/toggle", hd.ToggleStaggerStrategy)
				staggerWrite.PUT("/config", hd.UpdateStaggerConfig)
			}

			// 调度管理
			scheduleGroup := protected.Group("/schedules")
			scheduleGroup.Use(middleware.RequirePermission(permChecker, "publish", "read"))
			{
				scheduleGroup.GET("/:id", hd.GetSchedule)
				scheduleGroup.GET("", hd.ListSchedules)
			}
			scheduleWrite := protected.Group("/schedules")
			scheduleWrite.Use(middleware.RequirePermission(permChecker, "publish", "execute"))
			{
				scheduleWrite.POST("", hd.CreateSchedule)
				scheduleWrite.PUT("/:id", hd.UpdateSchedule)
				scheduleWrite.DELETE("/:id", hd.DeleteSchedule)
				scheduleWrite.POST("/:id/enable", hd.EnableSchedule)
				scheduleWrite.POST("/:id/disable", hd.DisableSchedule)
			}

			// 监测管理
			monitorGroup := protected.Group("/monitor")
			monitorGroup.Use(middleware.RequirePermission(permChecker, "content", "read"))
			{
				monitorGroup.GET("/citations", hd.GetAICitations)
				monitorGroup.GET("/scores", hd.GetSourceScores)
				monitorGroup.GET("/competitors", hd.GetCompetitorMonitors)
				monitorGroup.GET("/competitors/our-score", hd.GetOurScore)
				monitorGroup.GET("/roi", hd.GetROIMetrics)
				monitorGroup.POST("/suggestions/generate", hd.GenerateSuggestions)
			}

			// 知识图谱实体管理
			knowledgeGroup := protected.Group("/knowledge")
			knowledgeGroup.Use(middleware.RequirePermission(permChecker, "content", "read"))
			{
				knowledgeGroup.GET("/entities/:id", hd.GetEntity)
				knowledgeGroup.GET("/entities", hd.ListEntities)
				knowledgeGroup.GET("/entities/search", hd.SearchEntities)
			}
			knowledgeWrite := protected.Group("/knowledge")
			knowledgeWrite.Use(middleware.RequirePermission(permChecker, "content", "create"))
			{
				knowledgeWrite.POST("/entities", hd.CreateEntity)
				knowledgeWrite.PUT("/entities/:id", hd.UpdateEntity)
				knowledgeWrite.DELETE("/entities/:id", hd.DeleteEntity)
			}

			// 系统管理（仅 admin 角色可访问，中间件内 admin 自动放行）
			systemGroup := protected.Group("/system")
			systemGroup.Use(middleware.RequirePermission(permChecker, "tenant", "update"))
			{
				systemGroup.GET("/configs", hd.GetSystemConfigs)
				systemGroup.PUT("/configs/:key", hd.UpdateSystemConfig)
				systemGroup.GET("/plugins", hd.ListPlugins)
				systemGroup.POST("/plugins", hd.InstallPlugin)
				systemGroup.PUT("/plugins/:id", hd.UpdatePlugin)
				systemGroup.DELETE("/plugins/:id", hd.DeletePlugin)
				systemGroup.GET("/webhooks", hd.ListWebhooks)
				systemGroup.POST("/webhooks", hd.CreateWebhook)
				systemGroup.PUT("/webhooks/:id", hd.UpdateWebhook)
				systemGroup.DELETE("/webhooks/:id", hd.DeleteWebhook)
				systemGroup.POST("/webhooks/:id/test", hd.TestWebhook)
				systemGroup.GET("/webhooks/:id/history", hd.GetWebhookHistory)
				systemGroup.GET("/plans", hd.GetPlans)
				systemGroup.GET("/permissions", hd.GetPermissionDefinitions)
			}

			// 模板管理
			templateGroup := protected.Group("/templates")
			templateGroup.Use(middleware.RequirePermission(permChecker, "content", "read"))
			{
				templateGroup.GET("", hd.ListTemplates)
				templateGroup.GET("/:id", hd.GetTemplate)
			}
			templateWrite := protected.Group("/templates")
			templateWrite.Use(middleware.RequirePermission(permChecker, "content", "create"))
			{
				templateWrite.POST("", hd.CreateTemplate)
				templateWrite.PUT("/:id", hd.UpdateTemplate)
				templateWrite.DELETE("/:id", hd.DeleteTemplate)
			}
		}
	}
}