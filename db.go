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
	baseGetThreadURL = "https://api.lolz.live/threads/"
	whitelistPath    = "whitelist.json"
	outputPath       = "threads_db.json"
)

func LoadDatabase() {
	fmt.Println("Loading database...")

	wl, err := loadWhitelist(whitelistPath)
	if err != nil {
		fmt.Println("Ошибка загрузки whitelist:", err)
		return
	}

	var threads []Thread

	for i, entry := range wl {
		thread, err := fetchThread(entry.ThreadID)
		if err != nil {
			fmt.Printf("Ошибка получения треда %d: %v\n", entry.ThreadID, err)
			continue
		}

		thread.Captcha = entry.Captcha

		threads = append(threads, *thread)

		printProgress(i+1, len(wl))
		fmt.Println(*thread)
		time.Sleep(5 * time.Second)
	}

	filteredThreads := FilterThreads(threads)

	if err := saveThreadsDB(filteredThreads); err != nil {
		fmt.Println("Ошибка сохранения threads_db.json:", err)
		return
	}

	fmt.Println("БД успешно загружена и сохранена.")
}

func UpdateDatabaseWithNewWhitelistThreads() {
	fmt.Println("🔍 Проверка новых тредов в whitelist.json...")

	whitelist, err := loadWhitelist(whitelistPath)
	if err != nil {
		fmt.Println("Ошибка загрузки whitelist.json:", err)
		return
	}

	var db []Thread
	data, err := os.ReadFile(outputPath)
	if err == nil {
		json.Unmarshal(data, &db)
	}

	existing := make(map[int64]bool)
	for _, thread := range db {
		existing[thread.ThreadId] = true
	}

	var newCount int
	for _, entry := range whitelist {
		if existing[entry.ThreadID] {
			continue
		}

		thread, err := fetchThread(entry.ThreadID)
		if err != nil {
			fmt.Printf("⚠️ Ошибка при получении треда %d: %v\n", entry.ThreadID, err)
			continue
		}

		thread.Captcha = entry.Captcha
		db = append(db, *thread)
		newCount++

		fmt.Printf("✅ Добавлен тред %d (%s)\n", thread.ThreadId, thread.ThreadTitle)
		time.Sleep(5 * time.Second)
	}

	if newCount > 0 {
		if err := saveThreadsDB(db); err != nil {
			fmt.Println("❌ Ошибка сохранения threads_db.json:", err)
			return
		}
		fmt.Printf("💾 Добавлено %d новых тредов.\n", newCount)
	} else {
		fmt.Println("🟢 Новых тредов нет. БД актуальна.")
	}
}

func loadWhitelist(path string) ([]WhitelistEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var list []WhitelistEntry
	err = json.NewDecoder(f).Decode(&list)
	return list, err
}

func fetchThread(threadID int64) (*Thread, error) {
	req, err := http.NewRequest("GET", baseGetThreadURL+fmt.Sprint(threadID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", "Bearer "+bearerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var data ThreadResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data.Thread, err
}

func saveThreadsDB(data []Thread) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func printProgress(current, total int) {
	const steps = 10
	progressPerStep := total / steps
	if progressPerStep == 0 {
		progressPerStep = 1
	}

	completed := current / progressPerStep
	if completed > steps {
		completed = steps
	}
	bar := "[" + strings.Repeat("■", completed) + strings.Repeat("-", steps-completed) + "]"
	fmt.Printf("\rЗагрузка %d/%d %s", current, total, bar)
}
