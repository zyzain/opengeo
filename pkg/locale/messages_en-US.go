package locale

var enUSMessages = map[string]string{
	// General
	"success":                 "Success",
	"failed":                  "Failed",
	"operation_success":       "Operation successful",
	"operation_failed":        "Operation failed",
	"internal_error":          "Internal error",
	"invalid_parameter":       "Invalid parameter",
	"unauthorized":            "Unauthorized",
	"forbidden":               "Forbidden",
	"not_found":               "Resource not found",
	"already_exists":          "Resource already exists",
	"rate_limit_exceeded":     "Rate limit exceeded, please try again later",
	"server_busy":             "Server is busy, please try again later",

	// Auth
	"login_success":           "Login successful",
	"login_failed":            "Login failed",
	"register_success":        "Registration successful",
	"register_failed":         "Registration failed",
	"token_refreshed":         "Token refreshed",
	"invalid_credentials":     "Invalid username or password",
	"user_disabled":           "User account is disabled",
	"user_not_found":          "User not found",
	"username_exists":         "Username already exists",
	"missing_auth_header":     "Missing authorization header",
	"invalid_auth_format":     "Invalid authorization format",
	"invalid_token":           "Invalid or expired token",
	"password_mismatch":       "Passwords do not match",

	// Content
	"content_created":         "Content created successfully",
	"content_updated":         "Content updated successfully",
	"content_deleted":         "Content deleted successfully",
	"content_not_found":       "Content not found",
	"content_optimized":       "AI optimization completed",
	"content_optimize_failed": "AI optimization failed",
	"compliance_passed":       "Compliance check passed",
	"compliance_failed":       "Compliance check failed",
	"compliance_issues":       "Compliance issues found, please check details",

	// Publish
	"publish_task_created":    "Publish task created",
	"publish_cancelled":       "Task cancelled",
	"publish_cancel_failed":   "Cancellation failed",
	"publish_retry_submitted": "Retry task submitted",
	"publish_retry_failed":    "Retry failed",
	"publish_preview_failed":  "Failed to fetch preview data",

	// Account
	"account_created":         "Account created successfully",
	"account_updated":         "Account updated successfully",
	"account_deleted":         "Account deleted successfully",
	"account_not_found":       "Account not found",

	// Brand
	"brand_created":           "Brand created successfully",
	"brand_updated":           "Brand updated successfully",
	"brand_deleted":           "Brand deleted successfully",
	"brand_not_found":         "Brand not found",
	"metadata_saved":          "Metadata saved successfully",
	"metadata_save_failed":    "Metadata save failed",
	"entity_added":            "Entity added successfully",
	"entity_add_failed":       "Failed to add entity",
	"relation_added":          "Relation added successfully",
	"relation_add_failed":     "Failed to add relation",
	"entity_deleted":          "Entity deleted successfully",
	"entity_delete_failed":    "Failed to delete entity",
	"snapshot_created":        "Snapshot created successfully",
	"snapshot_create_failed":  "Failed to create snapshot",

	// Monitor
	"sync_created":            "Sync task created",
	"sync_failed":             "Sync failed",

	// Tenant
	"tenant_created":          "Tenant created successfully",
	"tenant_updated":          "Tenant updated successfully",
	"tenant_deleted":          "Tenant deleted",
	"tenant_not_found":        "Tenant not found",

	// System
	"plugin_installed":        "Plugin installed successfully",
	"config_updated":          "Configuration updated",
	"role_created":            "Role created successfully",
	"role_updated":            "Role updated successfully",
	"role_deleted":            "Role deleted successfully",
	"permission_updated":      "Permissions updated successfully",
	"permission_update_failed":"Permission update failed",
	"user_created":            "User created successfully",
	"user_updated":            "User updated successfully",
	"user_deleted":            "User deleted successfully",

	// Knowledge
	"entity_name_required":    "Entity name is required",
	"entity_type_required":    "Entity type is required",

	// Publish validation
	"title_required":          "Title is required",
	"body_required":           "Body is required",
	"title_too_long":          "Title exceeds 64 characters and may be truncated",
	"content_library":         "Content library",
	"published_content":       "Published content",
	"similar_content_found":   "Similar content found, consider adding unique perspectives and analysis",
	"use_synonyms":            "Consider using synonyms to replace some high-frequency words",
	"restructure_content":     "Consider rearranging paragraphs or restructuring content",
	"content_original":        "Content is highly original, keep it up",
	"content_too_short":       "Content is short, consider adding more details",

	// Roles
	"role_admin":              "Administrator",
	"role_operator":           "Operator",
	"role_viewer":             "Viewer",

	// Permissions
	"perm_user_create":        "Create user",
	"perm_user_read":          "View user",
	"perm_user_update":        "Update user",
	"perm_user_delete":        "Delete user",
	"perm_content_create":     "Create content",
	"perm_content_read":       "View content",
	"perm_content_update":     "Edit content",
	"perm_content_delete":     "Delete content",
	"perm_publish_create":     "Create publish",
	"perm_publish_read":       "View publish",
	"perm_publish_execute":    "Execute publish",
	"perm_role_create":        "Create role",
	"perm_role_read":          "View role",
	"perm_role_update":        "Update role",
	"perm_role_delete":        "Delete role",
	"perm_tenant_read":        "View tenant",
	"perm_tenant_update":      "Update tenant",

	// Permission groups
	"perm_group_content":      "Content Management",
	"perm_group_account":      "Account Management",
	"perm_group_publish":      "Publish Management",
	"perm_group_monitor":      "Monitor Management",
	"perm_group_system":       "System Management",

	// Plans
	"plan_starter":            "Starter",
	"plan_professional":       "Professional",
	"plan_enterprise":         "Enterprise",

	// Platforms
	"platform_wechat":         "WeChat Official Account",
	"platform_weibo":          "Weibo",
	"platform_douyin":         "Douyin",
	"platform_xiaohongshu":    "Xiaohongshu",
	"platform_zhihu":          "Zhihu",
	"platform_toutiao":        "Toutiao",

	// Validation
	"title_empty":             "Title cannot be empty",
	"body_empty":              "Body cannot be empty",

	// Startup messages
	"connecting_db":           "Connecting to database...",
	"db_connect_failed":       "Database connection failed",
	"db_connect_hint":         "Please ensure MySQL is running. Connection info can be configured via MYSQL_DSN environment variable",
	"db_connected":            "Database connected successfully",
	"db_instance_failed":      "Failed to get database instance",
	"db_migration_failed":     "Database migration failed",
	"tenant_create_failed":    "Failed to create default tenant",
	"tenant_check_done":       "Default tenant check completed",
	"perm_create_failed":      "Failed to create default permissions",
	"perm_check_done":         "Default permissions check completed",
	"role_create_failed":      "Failed to create default roles",
	"role_check_done":         "Default roles check completed",
	"role_perm_failed":        "Failed to create role-permission associations",
	"role_perm_check_done":    "Role-permission associations check completed",
	"admin_password_required": "ERROR: ADMIN_PASSWORD environment variable must be set in production",
	"admin_password_failed":   "Failed to generate admin password",
	"admin_create_failed":     "Failed to create admin account",
	"admin_role_failed":       "Failed to assign admin role",
	"admin_check_done":        "Admin account check completed",
	"content_table_failed":    "Content table migration failed",
	"knowledge_table_failed":  "Knowledge graph table migration failed",
	"knowledge_seed_failed":   "Knowledge graph seed data creation failed",
	"knowledge_seed_done":     "Knowledge graph seed data check completed",
	"gateway_shutting_down":   "Shutting down Gateway service...",

	// Content type
	"optimized_content":       "Optimized content",
	"content_structure_good":  "Content structure is good",
	"publish_task_created_msg":"Publish task created",
	"content_compliant":       "Content is compliant",
	"plugin_install_success":  "Plugin installed successfully",
}
