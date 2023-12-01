package rbac

import (
	"context"

	"order-service/constant"
	"order-service/middlewares"
)

func GetUserLogin(ctx context.Context) *middlewares.RBACData {
	return ctx.Value(constant.UserLogin).(*middlewares.RBACData) //nolint:forcetypeassert
}
