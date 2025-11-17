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

// GetWorld 行きたいワールド一覧取得
func GetWorld(w http.ResponseWriter, r *http.Request) {
	worldID := r.PathValue("world_id")

	worlds, err := GetWorldInfo(worldID)
	if err != nil {
		if errors.Is(err, NotFoundError) {
			slog.Warn(err.Error())
			http.Error(w, NotFoundError.Error(), http.StatusNotFound)
			return
		}
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

// DeleteWorld 行きたいワールドの登録解除
func DeleteWorld(w http.ResponseWriter, r *http.Request) {
	v := (r.Context()).Value("user_id")
	userID, _ := v.(string)
	worldID := r.PathValue("world_id")

	err := UnregisterWantGoWorld(worldID, userID)
	if err != nil {
		if errors.Is(err, NotRegisteredError) {
			slog.Warn(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func GetRecruitListInWorld(w http.ResponseWriter, r *http.Request) {
	worldID := r.PathValue("world_id")

	orderColumn := r.URL.Query().Get("order")
	if orderColumn == "" {
		orderColumn = "created_at"
	}
	if orderColumn != "created_at" && orderColumn != "joining_count" {
		msg := "sort に不正な値が指定されています"
		slog.Warn(msg, slog.String("value", orderColumn))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var isAsc bool
	switch r.URL.Query().Get("asc") {
	case "true":
		isAsc = true
	case "false":
		isAsc = false
	case "":
		isAsc = true
	default:
		msg := "type に不正な値が指定されています"
		slog.Warn(msg, slog.String("value", r.URL.Query().Get("asc")))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	recruits, err := ListRecruitInWorld(worldID, orderColumn, isAsc)
	if err != nil {
		if errors.Is(err, NotFoundError) {
			slog.Warn(err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(recruits)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
