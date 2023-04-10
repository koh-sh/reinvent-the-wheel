//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var messages []map[string]string

func returnSystemPrompt() string {
	// systemプロンプトを指定できます。空欄の場合は設定しません。
	// 20230410 "gpt-3.5-turbo-0301 does not always pay strong attention to system messages. Future models will be trained to pay stronger attention to system messages."
	systemPrompt := ``
	return systemPrompt
}

func parseOpenAIResponse(responseBody []byte) (string, error) {
	var responseMap map[string]interface{}

	if err := json.Unmarshal(responseBody, &responseMap); err != nil {
		return "", err
	}
	content := responseMap["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	if len(content) == 0 {
		return "", errors.New("unexpected response choices format")
	}

	return content, nil
}

func requestOpenAI() string {
	// 環境変数からDEBUGかどうか取得
	debug := os.Getenv("DEBUG")
	isDebug, err := strconv.ParseBool(debug)
	if err != nil {
		isDebug = false
	}
	// 環境変数からAPIキーを取得
	apiKey := os.Getenv("OPENAI_API_KEY")

	// OpenAIのエンドポイントURL
	endpointURL := "https://api.openai.com/v1/chat/completions"

	// APIリクエストのためのパラメータを設定
	data := map[string]interface{}{
		"model":       "gpt-3.5-turbo",
		"messages":    messages,
		"temperature": 0.7,
	}

	if isDebug {
		log.Printf("[DEBUG] Conversation sending to ChatGPT: %s", messages)
	}

	// APIリクエストを作成
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	req, _ := http.NewRequest("POST", endpointURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// APIリクエストを送信してレスポンスを取得
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if isDebug {
		log.Printf("[DEBUG] Response from ChatGPT: %s", string(respBody))
	}

	output, err := parseOpenAIResponse(respBody)
	if err != nil {
		log.Fatal(err)
	}
	return output

}

func topHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)
	tmpl := template.Must(template.ParseFiles("templates/top.html"))
	if err := tmpl.Execute(w, messages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	message := r.FormValue("message")
	if message == "" {
		http.Error(w, "empty message", http.StatusBadRequest)
		return
	}

	log.Printf("Sent Message to ChatGPT: %s", message)

	// 会話の最初だった場合プロンプトを追加
	systemPrompt := returnSystemPrompt()
	if messages == nil && systemPrompt != "" {
		chat := map[string]string{
			"role":    "system",
			"content": systemPrompt,
		}
		messages = append(messages, chat)
	}

	chat := map[string]string{
		"role":    "user",
		"content": message,
	}

	messages = append(messages, chat)
	output := requestOpenAI()
	log.Printf("Answer from ChatGPT: %s", output)
	chat = map[string]string{
		"role":    "assistant",
		"content": output,
	}

	messages = append(messages, chat)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func clearHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)
	log.Printf("Clearing messages history")
	messages = nil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {

	http.HandleFunc("/", topHandler)
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/clear", clearHandler)
	log.Printf("Accepting Web access on http://0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)

}
