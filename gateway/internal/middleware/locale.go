package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"opengeo/pkg/locale"
)

func Locale() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		acceptLang := string(c.GetHeader("Accept-Language"))
		loc := locale.ParseAcceptLanguage(acceptLang)
		ctx = locale.WithLocale(ctx, loc)
		c.Next(ctx)
	}
}
