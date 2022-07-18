package aliyun

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	dysmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type SendMsg struct {
	mobile          string
	accessKeyId     string
	accessKeySecret string
}

func NewSendMsg(mobile, accessKeyId, accessKeySecret string) (*SendMsg, error) {
	if mobile == "" || accessKeyId == "" || accessKeySecret == "" {
		return nil, errors.New("aliyun sms server parameter exception")
	}
	return &SendMsg{
		mobile:          mobile,
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
	}, nil
}

func (sm *SendMsg) SendMsg(msg string) error {
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(sm.accessKeyId, sm.accessKeySecret)
	client, err := dysmsapi.NewClientWithOptions("cn-hangzhou", config, credential)
	if err != nil {
		return err
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = "阿里云短信测试"
	request.TemplateCode = "SMS_154950909"
	request.PhoneNumbers = sm.mobile
	param := map[string]string{"code": msg}
	p, _ := json.Marshal(param)
	request.TemplateParam = string(p)

	resp, err := client.SendSms(request)
	if err != nil {
		return err
	}
	log.Printf("resp code: %s, message: %s", resp.Code, resp.Code)
	return nil
}

// GoSendMsg 异步发送
// NOTICE: 这里使用协程不影响后续盖楼的业务
func (sm *SendMsg) GoSendMsg(msg string) {
	go func() {
		err := sm.SendMsg(msg)
		if err != nil {
			log.Print("[send sms err] err: %s", err)
		}
	}()
}
