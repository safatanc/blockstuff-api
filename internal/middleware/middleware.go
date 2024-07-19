package middleware

import (
	"context"
	"net/http"

	"github.com/safatanc/blockstuff-api/pkg/jwthelper"
	"github.com/safatanc/blockstuff-api/pkg/response"
)

type Middleware struct {
}

func New() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			response.Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		claims, err := jwthelper.GetClaims(authorization)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(context.Background(), "claims", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
