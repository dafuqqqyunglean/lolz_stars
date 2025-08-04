package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	tradeLink = "https://steamcommunity.com/tradeoffer/new/?partner=287193132&token=FlOdqM2i"
)

func RunSteamCheat() error {
	url := "https://prod-api.lolz.live/threads?forum_id=849&page=1&limit=80"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+bearerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var data ThreadsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	fmt.Printf("Успешно загружено %d тредов\n", len(data.Threads))

	for _, thread := range data.Threads {
		if thread.ThreadIsClosed {
			continue
		}

		postBody := FormatPost(thread, fmt.Sprintf("[URL=\"%s\"]%s[/URL]", tradeLink, tradeLink))

		if err := sendPost(thread.ThreadId, postBody); err != nil {
			fmt.Printf("Ошибка при отправке в тред %d: %v\n", thread.ThreadId, err)
			continue
		}

		fmt.Printf("Успешно отправлено в тред %d\n", thread.ThreadId)

		time.Sleep(5 * time.Second)
	}

	return nil
}
