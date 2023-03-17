package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func CBCEncrypter(text []byte, key []byte, iv []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}
	// 填充
	paddText := PKCS7Padding(text, block.BlockSize())

	blockMode := cipher.NewCBCEncrypter(block, iv)

	// 加密
	result := make([]byte, len(paddText))
	blockMode.CryptBlocks(result, paddText)
	// 返回密文
	return result
}

// CBCDecrypter 解密
func CBCDecrypter(encrypter []byte, key []byte, iv []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	result := make([]byte, len(encrypter))
	blockMode.CryptBlocks(result, encrypter)
	// 去除填充
	result = UnPKCS7Padding(result)
	return result

}

func PKCS7Padding(text []byte, blockSize int) []byte {
	// 计算待填充的长度
	padding := blockSize - len(text)%blockSize
	var paddingText []byte
	if padding == 0 {
		// 已对齐，填充一整块数据，每个数据为 blockSize
		paddingText = bytes.Repeat([]byte{byte(blockSize)}, blockSize)
	} else {
		// 未对齐 填充 padding 个数据，每个数据为 padding
		paddingText = bytes.Repeat([]byte{byte(padding)}, padding)
	}
	return append(text, paddingText...)
}

func UnPKCS7Padding(text []byte) []byte {
	// 取出填充的数据 以此来获得填充数据长度
	unPadding := int(text[len(text)-1])
	return text[:(len(text) - unPadding)]
}
