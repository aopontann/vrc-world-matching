package vrc_world_matching

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

func TestGetWorld(t *testing.T) {
	tests := []struct {
		Name          string
		WorldID       string
		WantCode      int
		PrepareTables Tables
	}{
		{
			Name:     "登録済みのワールドIDを指定",
			WorldID:  "world1",
			WantCode: http.StatusOK,
			PrepareTables: Tables{
				World: []World{
					{ID: "world1", Name: "name", Thumbnail: "sanume1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
			},
		},
		{
			Name:     "登録していないワールドIDを指定",
			WorldID:  "world2",
			WantCode: http.StatusNotFound,
			PrepareTables: Tables{
				World: []World{
					{ID: "world1", Name: "name", Thumbnail: "thumbnail1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
			},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /worlds/{world_id}", AuthMiddleware(GetWorld))

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if err := CleanUp(); err != nil {
				t.Error(err)
			}
			if err := SetUp(test.PrepareTables); err != nil {
				t.Error(err)
			}

			r := httptest.NewRequest(http.MethodGet, "/worlds/"+test.WorldID, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)

			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != test.WantCode {
				t.Fatalf("got %d; want %d", resp.StatusCode, test.WantCode)
			}

			t.Log(string(body))
		})
	}

	if err := CleanUp(); err != nil {
		t.Error(err)
	}
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
				t.Errorf("got %d; want %d", resp.StatusCode, test.WantCode)
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
				t.Fatalf("got %d; want %d", resp.StatusCode, test.WantCode)
			}

			t.Log(string(body))
		})
	}

	if err := CleanUp(); err != nil {
		t.Error(err)
	}
}

func TestGetRecruitListInWorld(t *testing.T) {
	// テストデータのボリュームが多いので、共通するデータを使用するテストケースでは以下の変数を使用
	commonPrepareTables := Tables{
		World: []World{
			{ID: "world1", Name: "name1", Thumbnail: "thumbnail1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		},
		User: []User{
			{ID: "user1", Name: "name1", Icon: "icon1", Bio: "bio1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
			{ID: "user2", Name: "name2", Icon: "icon2", Bio: "bio2", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
			{ID: "user3", Name: "name3", Icon: "icon3", Bio: "bio3", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		},
		Result: []Recruit{
			{ID: "recruit1", Content: "content1", Closed: false, UserID: "user1", WorldID: "world1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
			{ID: "recruit2", Content: "content2", Closed: false, UserID: "user1", WorldID: "world1", CreatedAt: time.Now().UTC().Add(1 * time.Hour), UpdatedAt: time.Now().UTC()},
		},
		JoinMember: []JoinMember{
			{RecruitID: "recruit1", UserID: "user2", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
			{RecruitID: "recruit1", UserID: "user3", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
			{RecruitID: "recruit2", UserID: "user2", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		},
		Message: []Message{
			{ID: "message1", RecruitID: "recruit1", UserID: "user1", Content: "content1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
			{ID: "message2", RecruitID: "recruit1", UserID: "user2", Content: "content2", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
			{ID: "message3", RecruitID: "recruit2", UserID: "user1", Content: "content3", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		},
	}

	tests := []struct {
		Name          string
		WorldID       string
		Order         string
		Asc           string
		WantCode      int
		PrepareTables Tables
		ExpectedValue []RecruitDetail
	}{
		{
			Name:          "ワールドの募集リストを取得(クエリパラメータなし)",
			WorldID:       "world1",
			WantCode:      http.StatusOK,
			PrepareTables: commonPrepareTables,
			ExpectedValue: []RecruitDetail{
				{ID: "recruit1", Content: "content1", UserID: "user1", Username: "name1", UserIcon: "icon1", WorldID: "world1", WorldName: "name1", WorldThumbnail: "thumbnail1", JoiningCount: 2, MessagesCount: 2, CreatedAt: time.Now().UTC()},
				{ID: "recruit2", Content: "content2", UserID: "user1", Username: "name1", UserIcon: "icon1", WorldID: "world1", WorldName: "name1", WorldThumbnail: "thumbnail1", JoiningCount: 1, MessagesCount: 1, CreatedAt: time.Now().UTC().Add(1 * time.Hour)},
			},
		},
		{
			Name:          "ワールドの募集リストを取得(created_at のみ)",
			WorldID:       "world1",
			Order:         "created_at",
			WantCode:      http.StatusOK,
			PrepareTables: commonPrepareTables,
			ExpectedValue: []RecruitDetail{
				{ID: "recruit1", Content: "content1", UserID: "user1", Username: "name1", UserIcon: "icon1", WorldID: "world1", WorldName: "name1", WorldThumbnail: "thumbnail1", JoiningCount: 2, MessagesCount: 2, CreatedAt: time.Now().UTC()},
				{ID: "recruit2", Content: "content2", UserID: "user1", Username: "name1", UserIcon: "icon1", WorldID: "world1", WorldName: "name1", WorldThumbnail: "thumbnail1", JoiningCount: 1, MessagesCount: 1, CreatedAt: time.Now().UTC().Add(1 * time.Hour)},
			},
		},
		{
			Name:          "ワールドの募集リストを取得(created_at asc)",
			WorldID:       "world1",
			Order:         "created_at",
			Asc:           "true",
			WantCode:      http.StatusOK,
			PrepareTables: commonPrepareTables,
			ExpectedValue: []RecruitDetail{
				{ID: "recruit1", Content: "content1", UserID: "user1", Username: "name1", UserIcon: "icon1", WorldID: "world1", WorldName: "name1", WorldThumbnail: "thumbnail1", JoiningCount: 2, MessagesCount: 2, CreatedAt: time.Now().UTC()},
				{ID: "recruit2", Content: "content2", UserID: "user1", Username: "name1", UserIcon: "icon1", WorldID: "world1", WorldName: "name1", WorldThumbnail: "thumbnail1", JoiningCount: 1, MessagesCount: 1, CreatedAt: time.Now().UTC().Add(1 * time.Hour)},
			},
		},
		{
			Name:          "ワールドの募集リストを取得(created_at desc)",
			WorldID:       "world1",
			Order:         "created_at",
			Asc:           "false",
			WantCode:      http.StatusOK,
			PrepareTables: commonPrepareTables,
			ExpectedValue: []RecruitDetail{
				{ID: "recruit2", Content: "content2", UserID: "user1", Username: "name1", UserIcon: "icon1", WorldID: "world1", WorldName: "name1", WorldThumbnail: "thumbnail1", JoiningCount: 1, MessagesCount: 1, CreatedAt: time.Now().UTC().Add(1 * time.Hour)},
				{ID: "recruit1", Content: "content1", UserID: "user1", Username: "name1", UserIcon: "icon1", WorldID: "world1", WorldName: "name1", WorldThumbnail: "thumbnail1", JoiningCount: 2, MessagesCount: 2, CreatedAt: time.Now().UTC()},
			},
		},
		{
			Name:          "ワールドの募集リストを取得(joining_countのみ)",
			WorldID:       "world1",
			Order:         "joining_count",
			WantCode:      http.StatusOK,
			PrepareTables: commonPrepareTables,
			ExpectedValue: []RecruitDetail{
				{ID: "recruit2", Content: "content2", UserID: "user1", Username: "name1", UserIcon: "icon1", WorldID: "world1", WorldName: "name1", WorldThumbnail: "thumbnail1", JoiningCount: 1, MessagesCount: 1, CreatedAt: time.Now().UTC().Add(1 * time.Hour)},
				{ID: "recruit1", Content: "content1", UserID: "user1", Username: "name1", UserIcon: "icon1", WorldID: "world1", WorldName: "name1", WorldThumbnail: "thumbnail1", JoiningCount: 2, MessagesCount: 2, CreatedAt: time.Now().UTC()},
			},
		},
		{
			Name:          "ワールドの募集リストを取得(joining_count asc)",
			WorldID:       "world1",
			Order:         "joining_count",
			Asc:           "true",
			WantCode:      http.StatusOK,
			PrepareTables: commonPrepareTables,
			ExpectedValue: []RecruitDetail{
				{ID: "recruit2", Content: "content2", UserID: "user1", Username: "name1", UserIcon: "icon1", WorldID: "world1", WorldName: "name1", WorldThumbnail: "thumbnail1", JoiningCount: 1, MessagesCount: 1, CreatedAt: time.Now().UTC().Add(1 * time.Hour)},
				{ID: "recruit1", Content: "content1", UserID: "user1", Username: "name1", UserIcon: "icon1", WorldID: "world1", WorldName: "name1", WorldThumbnail: "thumbnail1", JoiningCount: 2, MessagesCount: 2, CreatedAt: time.Now().UTC()},
			},
		},
		{
			Name:          "ワールドの募集リストを取得(joining_count desc)",
			WorldID:       "world1",
			Order:         "joining_count",
			Asc:           "false",
			WantCode:      http.StatusOK,
			PrepareTables: commonPrepareTables,
			ExpectedValue: []RecruitDetail{
				{ID: "recruit1", Content: "content1", UserID: "user1", Username: "name1", UserIcon: "icon1", WorldID: "world1", WorldName: "name1", WorldThumbnail: "thumbnail1", JoiningCount: 2, MessagesCount: 2, CreatedAt: time.Now().UTC()},
				{ID: "recruit2", Content: "content2", UserID: "user1", Username: "name1", UserIcon: "icon1", WorldID: "world1", WorldName: "name1", WorldThumbnail: "thumbnail1", JoiningCount: 1, MessagesCount: 1, CreatedAt: time.Now().UTC().Add(1 * time.Hour)},
			},
		},
		{
			Name:          "ワールドの募集リストを取得(order 不正)",
			WorldID:       "world1",
			Order:         "ERROR",
			WantCode:      http.StatusBadRequest,
			PrepareTables: commonPrepareTables,
			ExpectedValue: nil,
		},
		{
			Name:          "ワールドの募集リストを取得(asc 不正)",
			WorldID:       "world1",
			Order:         "created_at",
			Asc:           "ERROR",
			WantCode:      http.StatusBadRequest,
			PrepareTables: commonPrepareTables,
			ExpectedValue: nil,
		},
		{
			Name:          "テーブルに存在しないワールドIDを指定",
			WorldID:       "world2",
			WantCode:      http.StatusNotFound,
			PrepareTables: commonPrepareTables,
			ExpectedValue: nil,
		},
		{
			Name:     "ワールドの募集リストを取得（募集が一つもない）",
			WorldID:  "world2",
			WantCode: http.StatusOK,
			PrepareTables: Tables{
				World: []World{
					{ID: "world1", Name: "name", Thumbnail: "thumbnail1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
					{ID: "world2", Name: "name", Thumbnail: "thumbnail2", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
				User: []User{
					{ID: "user1", Name: "name1", Icon: "icon1", Bio: "bio1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
				Result: []Recruit{
					{ID: "recruit1", Content: "content1", Closed: false, UserID: "user1", WorldID: "world1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
			},
			ExpectedValue: nil,
		},
		{
			Name:     "ワールドの募集リストを取得（募集が全て締め切り済み）",
			WorldID:  "world1",
			WantCode: http.StatusOK,
			PrepareTables: Tables{
				World: []World{
					{ID: "world1", Name: "name", Thumbnail: "thumbnail1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
				User: []User{
					{ID: "user1", Name: "name1", Icon: "icon1", Bio: "bio1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
				Result: []Recruit{
					{ID: "recruit1", Content: "content1", Closed: true, UserID: "user1", WorldID: "world1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
					{ID: "recruit2", Content: "content2", Closed: true, UserID: "user1", WorldID: "world1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				},
			},
			ExpectedValue: nil,
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /worlds/{world_id}/recruits", AuthMiddleware(GetRecruitListInWorld))

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if err := CleanUp(); err != nil {
				t.Error(err)
			}
			if err := SetUp(test.PrepareTables); err != nil {
				t.Error(err)
			}

			baseURL, err := url.Parse("/worlds/" + test.WorldID + "/recruits")
			q := baseURL.Query()
			if test.Order != "" {
				q.Set("order", test.Order)
			}
			if test.Asc != "" {
				q.Set("asc", test.Asc)
			}
			baseURL.RawQuery = q.Encode()

			r := httptest.NewRequest(http.MethodGet, baseURL.String(), nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)

			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != test.WantCode {
				t.Fatalf("got %d; want %d", resp.StatusCode, test.WantCode)
			}

			// 200以外の場合、レスポンスボディのチェックはしないためここで終了
			if resp.StatusCode != http.StatusOK {
				return
			}

			var result []RecruitDetail
			if err := json.Unmarshal(body, &result); err != nil {
				t.Error(err)
			}

			// created_at, updated_at で テストデータを宣言した時刻とデータ登録時の時刻に差が発生してしまうため、10秒の誤差を許容するようにする
			opts := cmpopts.EquateApproxTime(10 * time.Second)

			if diff := cmp.Diff(test.ExpectedValue, result, opts); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}

	if err := CleanUp(); err != nil {
		t.Error(err)
	}
}
