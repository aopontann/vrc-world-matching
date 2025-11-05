package vrc_world_matching

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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
		Name     string
		WorldID  string
		WantCode int
		PreFunc  func()
	}{
		{
			Name:     "新規登録（正常）",
			WorldID:  "wrld_20821acf-414a-454d-aa3d-be9dcd243b6d",
			WantCode: http.StatusCreated,
			PreFunc:  func() {},
		},
		{
			Name:     "新規登録（ワールドID不正）",
			WorldID:  "wrld_20821acf-414a-454d-aa3d-be9dcd2ERROR",
			WantCode: http.StatusBadRequest,
			PreFunc:  func() {},
		},
		{
			Name:     "新規登録（ワールド既に登録済み、行きたい登録はしていない）",
			WorldID:  "wrld_20821acf-414a-454d-aa3d-be9dcd243b6d",
			WantCode: http.StatusCreated,
			PreFunc: func() {
				err := RegisterWorld("wrld_20821acf-414a-454d-aa3d-be9dcd243b6d")
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name:     "新規登録（ワールド既に登録済み、行きたい登録済み）",
			WorldID:  "wrld_20821acf-414a-454d-aa3d-be9dcd243b6d",
			WantCode: http.StatusBadRequest,
			PreFunc: func() {
				err := RegisterWorld("wrld_20821acf-414a-454d-aa3d-be9dcd243b6d")
				if err != nil {
					t.Fatal(err)
				}
				err = RegisterWantGoWorld("wrld_20821acf-414a-454d-aa3d-be9dcd243b6d", "user1")
				if err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	// テスト開始直前処理
	_, err := db.Exec("TRUNCATE test.worlds")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("TRUNCATE test.want_go")
	if err != nil {
		t.Fatal(err)
	}

	// ハンドラの設定
	mux := http.NewServeMux()
	mux.HandleFunc("POST /worlds/{world_id}", AuthMiddleware(PostWorld))

	// テスト開始
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// テストデータ投入
			test.PreFunc()
			// 1ケース終了ごとにデータを全削除
			defer func() {
				_, err := db.Exec("TRUNCATE test.worlds")
				if err != nil {
					t.Fatal(err)
				}
				_, err = db.Exec("TRUNCATE test.want_go")
				if err != nil {
					t.Fatal(err)
				}
			}()

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
}
