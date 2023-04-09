//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Chat struct {
	Input  string
	Output string
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

func requestOpenAI(input string) string {
	DEBUG := false
	// 環境変数からAPIキーを取得
	apiKey := os.Getenv("OPENAI_API_KEY")

	// OpenAIのエンドポイントURL
	endpointURL := "https://api.openai.com/v1/chat/completions"

	// APIリクエストのためのパラメータを設定
	data := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "user", "content": input},
		},
		"temperature": 0.7,
	}

	// APIリクエストを作成
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	req, _ := http.NewRequest("POST", endpointURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// APIリクエストを送信してレスポンスを取得
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if DEBUG {
		fmt.Println(string(respBody))
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
	if err := tmpl.Execute(w, nil); err != nil {
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

	output := requestOpenAI(message)
	chat := Chat{
		Input:  message,
		Output: output,
	}

	// メッセージをログに出力する
	tmpl := template.Must(template.ParseFiles("templates/log.html"))
	if err := tmpl.Execute(w, chat); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func main() {

	http.HandleFunc("/", topHandler)
	http.HandleFunc("/send", sendHandler)
	http.ListenAndServe(":8080", nil)

}
