package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func main() {

	jsonData := []byte(`{
		"username": "testuser2",
		"password": "1234"
	}`)

	resp, err := http.Post(
		"http://localhost:8000/register",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// ✅ Read response body (IMPORTANT — you were missing this)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:", string(body))
}
