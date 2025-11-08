package vrc_world_matching

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetWorldList(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/worlds", nil)
	w := httptest.NewRecorder()

	GetWorldList(w, r)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got %d; want %d", resp.StatusCode, http.StatusOK)
	}

	t.Log(string(body))
}

func TestPostWorld(t *testing.T) {
	tests := []struct {
		Name          string
		WorldID       string
		WantCode      int
		PrepareTables Tables
	}{
		{
			Name:          "新規登録（正常）",
			WorldID:       "wrld_20821acf-414a-454d-aa3d-be9dcd243b6d",
			WantCode:      http.StatusCreated,
			PrepareTables: Tables{},
		},
		{
			Name:          "新規登録（ワールドID不正）",
			WorldID:       "wrld_20821acf-414a-454d-aa3d-be9dcd2ERROR",
			WantCode:      http.StatusBadRequest,
			PrepareTables: Tables{},
		},
		{
			Name:     "新規登録（ワールド既に登録済み、行きたい登録はしていない）",
			WorldID:  "world1",
			WantCode: http.StatusCreated,
			PrepareTables: Tables{
				World: []World{
					{ID: "world1", Name: "name", Thumbnail: "thumbnail", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
			},
		},
		{
			Name:     "新規登録（ワールド既に登録済み、行きたい登録済み）",
			WorldID:  "world1",
			WantCode: http.StatusBadRequest,
			PrepareTables: Tables{
				World: []World{
					{ID: "world1", Name: "name", Thumbnail: "thumbnail", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
				WantGo: []WantGo{
					{UserID: "user1", WorldID: "world1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
			},
		},
	}

	// ハンドラの設定
	mux := http.NewServeMux()
	mux.HandleFunc("POST /worlds/{world_id}", AuthMiddleware(PostWorld))

	// テスト開始
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if err := CleanUp(); err != nil {
				t.Error(err)
			}
			if err := SetUp(test.PrepareTables); err != nil {
				t.Error(err)
			}

			r := httptest.NewRequest(http.MethodPost, "/worlds/"+test.WorldID, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)

			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != test.WantCode {
				t.Errorf("got %d; want %d", resp.StatusCode, http.StatusOK)
			}

			t.Log(string(body))
		})
	}

	if err := CleanUp(); err != nil {
		t.Error(err)
	}
}

func TestDeleteWorld(t *testing.T) {
	tests := []struct {
		Name          string
		WorldID       string
		WantCode      int
		PrepareTables Tables
	}{
		{
			Name:     "行きたい登録したワールドの解除",
			WorldID:  "world1",
			WantCode: http.StatusOK,
			PrepareTables: Tables{
				World: []World{
					{ID: "world1", Name: "name", Thumbnail: "sanume1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
				WantGo: []WantGo{
					{UserID: "user1", WorldID: "world1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
			},
		},
		{
			Name:     "行きたい登録していないワールドIDを指定",
			WorldID:  "world1",
			WantCode: http.StatusBadRequest,
			PrepareTables: Tables{
				World: []World{
					{ID: "world1", Name: "name", Thumbnail: "thumbnail1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
					{ID: "world2", Name: "name", Thumbnail: "thumbnail2", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
				WantGo: []WantGo{
					{UserID: "user1", WorldID: "world2", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
			},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /worlds/{world_id}", AuthMiddleware(DeleteWorld))

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if err := CleanUp(); err != nil {
				t.Error(err)
			}
			if err := SetUp(test.PrepareTables); err != nil {
				t.Error(err)
			}

			r := httptest.NewRequest(http.MethodDelete, "/worlds/"+test.WorldID, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)

			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != test.WantCode {
				t.Fatalf("got %d; want %d", resp.StatusCode, http.StatusOK)
			}

			t.Log(string(body))
		})
	}

	if err := CleanUp(); err != nil {
		t.Error(err)
	}
}
