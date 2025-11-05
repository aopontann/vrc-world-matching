package vrc_world_matching

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

// GetWorldList 行きたいワールド一覧取得
func GetWorldList(w http.ResponseWriter, r *http.Request) {
	v := (r.Context()).Value("user_id")
	userID, _ := v.(string)
	slog.Info(userID)

	worlds, err := ListWorld()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(worlds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PostWorld 行きたいワールド登録
func PostWorld(w http.ResponseWriter, r *http.Request) {
	v := (r.Context()).Value("user_id")
	userID, _ := v.(string)
	worldID := r.PathValue("world_id")

	// ワールド情報を登録
	err := RegisterWorld(worldID)
	if err != nil {
		// 存在しないワールドIDの場合、400を返す
		if errors.Is(err, VRChatAPINotExistWorldError) {
			slog.Warn(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 行きたいワールドを登録
	err = RegisterWantGoWorld(worldID, userID)
	if err != nil {
		// 既に行きたいワールドとして登録済みの場合、400を返す
		if errors.Is(err, AlreadyRegisteredError) {
			slog.Warn(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("OK"))
}
