package handler

import (
	"fmt"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/service"
	"github.com/Digital-AIR/bizio-ecommerce/internal/session"
	"net/http"
	"strconv"
)

//func HandleUserLogout(writer http.ResponseWriter, request *http.Request) {
//	if !session.IsAuthenticated(request) {
//		common.HandleUnAuthRes(writer, map[string]string{})
//		return
//	}
//
//	err := session.Logout(writer, request)
//
//	if err != nil {
//		common.HandleErrorRes(writer, map[string]string{"message": "failed to logout"})
//	}
//}

func HandleUserLogin(writer http.ResponseWriter, request *http.Request) {
	userLoginDto := &common.UserLoginDto{}

	// read body
	if ok := common.ReadReqPayload(writer, request, &userLoginDto); !ok {
		return
	}
	// validate body
	if errMap := validateLoginReq(userLoginDto); len(errMap) > 1 {
		common.HandleErrorRes(writer, errMap)
		return
	}

	// check username & password is present
	if !service.CheckUserNamePassword(userLoginDto.Username, userLoginDto.Password) {
		common.HandleUnAuthRes(writer, map[string]string{})
		return
	}

	// create user session
	err := session.CreateUserSession(writer, request)

	if err != nil {
		fmt.Println("failed to create session")
	}

}

func validateLoginReq(dto *common.UserLoginDto) map[string]string {
	errMap := make(map[string]string)

	if dto.Username == "" {
		errMap["username"] = "username is required"
	}

	if dto.Password == "" {
		errMap["password"] = "password is required"
	}
	return errMap
}

func HandleUserCreate(writer http.ResponseWriter, request *http.Request) {
	userDto := &common.UserReg{}

	// read body
	if ok := common.ReadReqPayload(writer, request, &userDto); !ok {
		return
	}

	// validate body
	if errMap := validateCreateReq(userDto); len(errMap) > 1 {
		common.HandleErrorRes(writer, errMap)
		return
	}

	// save user
	user := service.SaveUser(userDto)

	common.SendOkRes(writer, map[string]string{"id": strconv.Itoa(int(user.ID))})
}

func validateCreateReq(dto *common.UserReg) map[string]string {
	errMap := make(map[string]string)

	if dto.Username == "" {
		errMap["username"] = "username is required"
	}

	if dto.Email == "" {
		errMap["email"] = "email is required"
	}

	if dto.Password == "" {
		errMap["password"] = "password is required"
	}

	return errMap
}
