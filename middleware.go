package vrc_world_matching

import (
	"context"
	"net/http"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ヘッダーから取得
		// sessionKey := ""
		// インメモリDBからセッションキーから紐づくユーザIDを取得
		userID := "user1"
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
