package server

const (
	ServerSMSCodeTips     = "200000" // 提示有新的盖楼贴
	ServerSMSCodeBuilding = "200001" // 提示开始盖楼
)

type IMSM interface {
	SendMsg(msg string) error // 同步发送
	GoSendMsg(msg string)     // 异步发送
}
