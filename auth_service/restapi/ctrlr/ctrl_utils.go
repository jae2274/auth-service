package ctrlr

import (
	"context"
	"net/http"

	"github.com/jae2274/auth-service/auth_service/restapi/jwtresolver"
	"github.com/jae2274/auth-service/auth_service/restapi/middleware"
	"github.com/jae2274/goutils/llog"
)

func errorHandler(ctx context.Context, w http.ResponseWriter, err error) bool {
	if err != nil {
		llog.LogErr(ctx, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}
	return false
}

// 주의: ctx에 claims가 존재한다는 가정하에 사용되는 함수로, 존재하지 않을 경우 panic이 발생
func GetClaimsOrPatal(ctx context.Context) *jwtresolver.CustomClaims {
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		panic("no claims in context")
	}

	return claims
}
