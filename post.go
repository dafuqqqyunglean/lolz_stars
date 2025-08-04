package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	bearerToken       = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzUxMiJ9.eyJzdWIiOjM0MjE4ODksImlzcyI6Imx6dCIsImlhdCI6MTc1MzkwNDczMCwianRpIjoiODIyODMyIiwic2NvcGUiOiJiYXNpYyByZWFkIHBvc3QgY29udmVyc2F0ZSBwYXltZW50IGludm9pY2UgY2hhdGJveCBtYXJrZXQifQ.Pe8qiGVF23yG90PSQsW35NIdEMZ73Ryq0P1Umy5nj0bch-_7pZVku34CNacrUqELupLHdHz6gSdyZpxITAXc9mXNMyEe6pmRAhuv91YeaG1rylnAwgHRrX3wqoQ6vkTG3ZGA9JdYarACSYX_gpBehDslXbNIkCNAoI6Wq8dlE_Q"
	baseCreatePostURL = "https://api.lolz.live/posts?thread_id="
	mediaLink         = "yataklyblymamy/4"
)

func RunPostLoop() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("❌ Ошибка загрузки часового пояса:", err)
		loc = time.FixedZone("MSK", 3*60*60)
	}

	for {
		fmt.Println("\n🔄 Проверка тредов для постинга...")
		PostToThreads()

		nowMSK := time.Now().In(loc)
		fmt.Printf("(%s) ⏳ Ожидание 1 час \n", nowMSK.Format("02.01.2006 15:04:05 MST"))

		time.Sleep(61 * time.Minute)
	}
}

func PostToThreads() {
	fmt.Println("Загрузка базы...")

	data, err := os.ReadFile(outputPath)
	if err != nil {
		fmt.Println("Ошибка чтения threads_db.json:", err)
		return
	}

	var threads []Thread
	if err := json.Unmarshal(data, &threads); err != nil {
		fmt.Println("Ошибка разбора threads_db.json:", err)
		return
	}

	fmt.Printf("Загружено %d записей \n", len(threads))

	now := time.Now().Unix()
	updated := false

	for i, thread := range threads {
		if thread.ThreadIsClosed {
			continue
		}

		if thread.LastReplyDate != 0 && (now-thread.LastReplyDate) < thread.Restrictions.ReplyDelay {
			continue
		}

		postBody := FormatPost(thread, fmt.Sprintf("[MEDIA=telegram]%s[/MEDIA]", mediaLink))

		if err := sendPost(thread.ThreadId, postBody); err != nil {
			fmt.Printf("Ошибка при отправке в тред %d: %v\n", thread.ThreadId, err)
			continue
		}

		fmt.Printf("Успешно отправлено в тред %d\n", thread.ThreadId)
		threads[i].LastReplyDate = now
		updated = true

		time.Sleep(5 * time.Second)
	}

	if updated {
		if err := saveThreadsDB(threads); err != nil {
			fmt.Println("Ошибка при сохранении базы:", err)
		} else {
			fmt.Println("База обновлена.")
		}
	} else {
		fmt.Println("Нет подходящих тредов для ответа.")
	}
}

func sendPost(threadID int64, postContent string) error {
	url := fmt.Sprintf("%s%d", baseCreatePostURL, threadID)

	payload := map[string]string{
		"post_body": postContent,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", "Bearer "+bearerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ошибка API (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}
