//go:build ignore

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

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

func main() {
	DEBUG := false
	// 環境変数からAPIキーを取得
	apiKey := os.Getenv("OPENAI_API_KEY")

	// OpenAIのエンドポイントURL
	endpointURL := "https://api.openai.com/v1/chat/completions"

	// ユーザーからの入力を受け取る
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("You: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

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

	// ChatGPTからの応答を出力
	fmt.Println("ChatGPT:", output)
}
