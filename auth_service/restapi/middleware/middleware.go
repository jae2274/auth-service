package middleware

import (
	"context"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
	"github.com/jae2274/auth-service/auth_service/restapi/jwtresolver"
	"github.com/jae2274/goutils/llog"
)

type claimsKey string

const claimsKeyStr claimsKey = "claims"

func SetClaimsMW(jr *jwtresolver.JwtResolver) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString != "" {
				claims, isValid, err := jr.ParseToken(tokenString)
				if err != nil {
					llog.LogErr(r.Context(), err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if !isValid {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), claimsKeyStr, claims)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func CheckJustLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := GetClaims(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CheckHasAuthority(authority string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaims(r.Context())
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if !slices.Contains(claims.Authorities, authority) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetClaims(ctx context.Context) (*jwtresolver.CustomClaims, bool) {
	claims, ok := ctx.Value(claimsKeyStr).(*jwtresolver.CustomClaims)
	return claims, ok
}
