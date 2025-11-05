package vrc_world_matching

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type WorldInfoFromVRChatAPI struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Thumbnail string `json:"thumbnailImageUrl" db:"thumbnail"`
}

func GetWorldInfoFromVRChatAPI(worldID string) (*WorldInfoFromVRChatAPI, error) {
	url := "https://api.vrchat.cloud/api/1/worlds/" + worldID
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	req.Header.Add("User-Agent", os.Getenv("USER_AGENT"))

	res, err := client.Do(req)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	fmt.Println(string(body))

	// ワールドが存在しない場合
	if res.StatusCode == http.StatusNotFound {
		return nil, VRChatAPINotExistWorldError
	}
	if res.StatusCode != 200 {
		return nil, VRChatAPIError
	}

	var worldInfo WorldInfoFromVRChatAPI
	err = json.Unmarshal(body, &worldInfo)
	if err != nil {
		return nil, err
	}
	return &worldInfo, nil
}
