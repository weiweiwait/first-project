package routes

import (
	util "MyFirstProject/pkg/utils/jwt"
	"fmt"
	"testing"
)

//func TestMain(m *testing.M) {
//	re := conf.ConfigReader{FileName: "../../../config/locales/config.yaml"}
//	conf.InitConfigForTest(&re)
//	fmt.Println("Write tests on values: ", conf.Config)
//	m.Run()
//}

func TestUserModelEncryptMoney(t *testing.T) {
	claims, _ := util.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTQyMTk2MDksImlzcyI6Im1hbGwifQ.12q-3miWAwzg09RUK52U2YMQrF3ipBkP-CvZcNN_Sok")
	fmt.Println(claims.ID)
}
