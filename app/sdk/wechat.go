package sdk

import (
	"log"
	"time"
	"wechatRobot/app/config"
	"wechatRobot/grpc/wcf"

	"golang.org/x/sys/windows"
)

var (
	modWxSDK      = windows.MustLoadDLL(config.AppConfig.App.SdkPath + "sdk.dll")
	procWxInitSDK = modWxSDK.MustFindProc("WxInitSDK")
	wxClient      *wcf.Client
)

func Initialize(enableDebug bool) error {
	debug := uint32(0)
	if enableDebug {
		debug = 1
	}

	ret, _, err := procWxInitSDK.Call(uintptr(debug), uintptr(10086))
	if ret != 0 {
		return err
	}

	log.Println("SDK initialized successfully")

	time.Sleep(3 * time.Second)

	client, err := wcf.NewWCF("")
	if err != nil {
		return err
	}
	wxClient = client
	return nil
}

func Cleanup() {
	if wxClient != nil {
		_ = wxClient.Close()
	}
}

func GetClient() *wcf.Client {
	return wxClient
}

func WaitForLogin(client *wcf.Client) error {
	for !client.IsLogin() {
		time.Sleep(2 * time.Second)
	}
	log.Println("Login successful, initializing...")
	return nil
}
