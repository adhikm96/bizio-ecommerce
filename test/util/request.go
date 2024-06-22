package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func GetServerUrl(url string) string {
	return "http://localhost:8000" + url
}

func MakeReq(reqType string, url string, payload []byte, cookie *http.Cookie) ([]byte, *http.Response, error) {
	client := &http.Client{}

	request, err := http.NewRequest(reqType, GetServerUrl(url), bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println(err.Error())
		return nil, nil, err
	}

	if cookie != nil {
		request.AddCookie(cookie)
	}

	return getResp(client, request)
}

func getResp(client *http.Client, request *http.Request) ([]byte, *http.Response, error) {
	resp, err := client.Do(request)

	if err != nil {
		fmt.Println(err.Error())
		return nil, nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err.Error())
		return nil, resp, err
	}

	if resp.StatusCode != 200 {
		fmt.Println("status code is not 200")
		return data, resp, errors.New("status code is not 200")
	}

	return data, resp, nil
}
