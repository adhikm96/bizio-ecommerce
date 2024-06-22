package user

import (
	"encoding/json"
	"fmt"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/server"
	"github.com/Digital-AIR/bizio-ecommerce/test/util"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	go server.StartServer()
	code := m.Run()
	os.Exit(code)
}

func TestUserLoginFlow(t *testing.T) {

	time.Sleep(time.Second * 1)

	_, resp, _ := util.MakeReq("GET", "/ping", []byte(""), nil)

	if resp != nil && resp.StatusCode != 401 {
		t.Fail()
	}

	buff, _ := json.Marshal(&common.UserReg{
		Username: "ak00029",
		Email:    "ak00029@example.com",
		Password: "password",
	})

	_, _, err := util.MakeReq("POST", "/register", buff, nil)

	if err != nil {
		t.Fail()
	}

	buff, _ = json.Marshal(&common.UserReg{
		Username: "ak00029",
		Password: "password",
	})

	_, res, err := util.MakeReq("POST", "/login", buff, nil)

	if err != nil {
		t.Fail()
	}

	userCookie := ""

	for _, cookie := range res.Cookies() {
		if cookie.Name == "user-session" {
			userCookie = cookie.Value
		}
	}

	cookie := &http.Cookie{
		Name:    "user-session",
		Value:   userCookie,
		Expires: time.Now().Add(24 * time.Hour),
	}

	_, _, err = util.MakeReq("GET", "/ping", buff, cookie)

	if err != nil {
		t.Fail()
	}
}

func TestHandleUserCreate(t *testing.T) {
	time.Sleep(time.Second * 1)

	buff, _ := json.Marshal(&common.UserReg{
		Username: "ak00029",
		Email:    "ak00029@example.com",
		Password: "password",
	})

	data, _, err := util.MakeReq("POST", "/register", buff, nil)

	if err != nil {
		t.Fail()
	}

	fmt.Println(string(data))
}
