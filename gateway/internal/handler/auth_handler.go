package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

// Login 登录
func (h *Handler) Login(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Username string `json:"username" vd:"len($)>0"`
		Password string `json:"password" vd:"len($)>0"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	resp, err := h.userClient.Login(ctx, req.Username, req.Password)
	if err != nil {
		errResponse(c, err, "login failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// Register 注册
func (h *Handler) Register(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Username string `json:"username" vd:"len($)>0"`
		Password string `json:"password" vd:"len($)>0"`
		Email    string `json:"email"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	resp, err := h.userClient.Register(ctx, req.Username, req.Password, req.Email)
	if err != nil {
		errResponse(c, err, "registration failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// RefreshToken 刷新Token
func (h *Handler) RefreshToken(ctx context.Context, c *app.RequestContext) {
	var req struct {
		RefreshToken string `json:"refresh_token" vd:"len($)>0"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	resp, err := h.userClient.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		errResponse(c, err, "refresh token failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}
