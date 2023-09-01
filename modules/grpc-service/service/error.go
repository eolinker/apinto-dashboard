package service

type SendMsgError struct {
	UUID              string `json:"uuid"`
	NoticeChannelUUID string `json:"notice_channel_uuid"`
	Msg               string `json:"msg"`
}
