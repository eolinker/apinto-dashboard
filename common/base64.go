package common

import (
	"encoding/base64"
)

//Base64Encode 加密
func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

//Base64Decode 解密
func Base64Decode(src string) ([]byte, error) {
	// 对上面的编码结果进行base64解码
	decodeBytes, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}
	return decodeBytes, nil
}
