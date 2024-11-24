package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080/users"

type Response struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

func TestUserFlow(t *testing.T) {
	client := &http.Client{}
	var token string
	var userID int = 1

	// 1. 注册用户
	t.Run("RegisterUser", func(t *testing.T) {
		payload := `{"nickname":"TestUser","phone_number":"+12345678901","password":"securePassword123"}`
		resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer([]byte(payload)))
		if err != nil {
			t.Fatalf("Failed to register user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := ioutil.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, string(body))
		}
	})

	// 2. 登录用户
	t.Run("LoginUser", func(t *testing.T) {
		payload := `{"phone_number":"+12345678901","password":"securePassword123"}`
		resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer([]byte(payload)))
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
		}

		var response Response
		json.NewDecoder(resp.Body).Decode(&response)
		token = response.Token
		fmt.Printf("Token: %s\n", token) // 调试打印 Token
		if token == "" {
			t.Fatal("Token not received")
		}
	})

	// 3. 更新用户名
	t.Run("UpdateNickname", func(t *testing.T) {
		payload := `{"nickname":"UpdatedUser"}`
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%d/nickname", baseURL, userID), bytes.NewBuffer([]byte(payload)))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+token) // 确保格式正确
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to update nickname: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
		}
	})
}