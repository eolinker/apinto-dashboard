package cli

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/urfave/cli/v2"
	"os"
	"path"
)

const (
	rsaBits     = 1024
	keysDirFlag = "dir"
)

func GenRSAKeys() *cli.Command {
	return &cli.Command{
		Name:  "gen-rsa",
		Usage: "generate 1024 bits RSA keys",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     keysDirFlag,
				Usage:    "证书文件生成的目录, 默认为 程序当前目录下 ./export/rsa/",
				Required: false,
				Value:    "./export/rsa",
			},
		},
		Action: GenerateRSAKeys,
	}
}

func GenerateRSAKeys(c *cli.Context) error {
	keysDir := c.String(keysDirFlag)

	privateKey, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		panic(err)
	}

	X509PrivateKey, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}
	//使用pem格式对x509输出的内容进行编码
	//创建文件保存私钥
	err = os.MkdirAll(keysDir, os.ModePerm)
	if err != nil {
		return err
	}
	privateFile, err := os.OpenFile(path.Join(keysDir, "private.pem"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer privateFile.Close()
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	//将数据保存到文件
	err = pem.Encode(privateFile, &privateBlock)
	if err != nil {
		return err
	}

	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//pem格式编码
	//创建用于保存公钥的文件
	publicFile, err := os.OpenFile(path.Join(keysDir, "public.pem"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer publicFile.Close()
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	//保存到文件
	err = pem.Encode(publicFile, &publicBlock)
	if err != nil {
		return err
	}

	return nil
}
