package cli

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/eolinker/apinto-dashboard/app/apserver-authorization/entity"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"time"
)

const (
	certFileName = "apinto.cert"
)

func newCertInfo(machineCode, company, edition string, begin, end time.Time, controllerCount, nodeCount int) *entity.CertInfo {
	return &entity.CertInfo{
		Company:         company,
		Edition:         edition,
		BeginTime:       begin.Unix(),
		EndTime:         end.Unix(),
		ControllerCount: controllerCount,
		NodeCount:       nodeCount,
		MachineCode:     machineCode,
	}
}

func buildCert(certInfo *entity.CertInfo, certDirPath string) error {
	//签名
	signature, err := signCert(certInfo)
	if err != nil {
		return err
	}
	//校验签名
	err = verifyCert(certInfo, signature)
	if err != nil {
		return err
	}

	//生成证书
	certInfo.Signature = signature
	return genCertFile(certInfo, certDirPath)
}

func signCert(certInfo *entity.CertInfo) (string, error) {
	signParam := getSignParam(certInfo)

	block, _ := pem.Decode(getPrivateKey())
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		return "", err
	}
	h := sha256.New()
	h.Write([]byte(signParam))
	hash := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA256, hash)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return "", err
	}

	out := base64.StdEncoding.EncodeToString(signature)
	return out, nil
}

func verifyCert(certInfo *entity.CertInfo, signature string) error {
	signParam := getSignParam(certInfo)

	sign, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(getPublicKey())
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	hash := sha256.New()
	hash.Write([]byte(signParam))
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, hash.Sum(nil), sign)
}

func genCertFile(cert *entity.CertInfo, certDirPath string) error {
	fileData, err := yaml.Marshal(cert)
	if err != nil {
		return err
	}

	err = os.MkdirAll(certDirPath, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path.Join(certDirPath, certFileName), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(string(fileData))
	return err
}

func getSignParam(cert *entity.CertInfo) string {
	param, _ := yaml.Marshal(cert)
	return string(param)
}
