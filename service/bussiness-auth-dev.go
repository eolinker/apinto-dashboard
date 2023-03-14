//go:build dev
// +build dev

package service

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	machine_code "github.com/eolinker/apinto-dashboard/app/apserver/machine-code"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/entry/system-entry"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-basic/uuid"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"time"
)

type IBussinessAuthService interface {
	GetActivationInfo(ctx context.Context) (*model.ActivationInfo, error)
	ActivateCert(ctx context.Context, certFile []byte) (*model.ActivationInfo, error)
	ReActivateCert(ctx context.Context, certFile []byte) (*model.ActivationInfo, error)
	GetMachineCode(ctx context.Context) (string, error)
	CheckCertValid(ctx context.Context) (bool, error)
	examineCert(ctx context.Context, certInfo *model.BussinessCertInfo, mac string) error
}

type bussinessAuthService struct {
	systemInfoStore store.ISystemInfoStore
	systemInfoCache cache.ISystemInfoCache
}

func newBussinessAuthService() IBussinessAuthService {
	b := &bussinessAuthService{}
	bean.Autowired(&b.systemInfoStore)
	bean.Autowired(&b.systemInfoCache)

	return b
}

func (b *bussinessAuthService) GetActivationInfo(ctx context.Context) (*model.ActivationInfo, error) {
	certData, err := b.systemInfoStore.GetSystemInfoByKey(ctx, enum.CertDBKey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &model.ActivationInfo{
				BussinessCertInfo: model.BussinessCertInfo{},
				DashboardID:       "",
			}, nil
		}
		return nil, err
	}
	cert := new(model.BussinessCertInfo)
	_ = yaml.Unmarshal(certData.Value, cert)

	//TODO 控制台多节点后这里要改
	dashboardIdData, err := b.systemInfoStore.GetSystemInfoByKey(ctx, enum.DashboardIdDBKey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &model.ActivationInfo{
				BussinessCertInfo: model.BussinessCertInfo{},
				DashboardID:       "",
			}, nil
		}
		return nil, err
	}

	return &model.ActivationInfo{
		BussinessCertInfo: *cert,
		DashboardID:       string(dashboardIdData.Value),
	}, nil
}

func (b *bussinessAuthService) ActivateCert(ctx context.Context, certFile []byte) (*model.ActivationInfo, error) {
	//校验证书
	machineCode := machine_code.GetMachineCode()
	certInfo := new(model.BussinessCertInfo)
	err := yaml.Unmarshal(certFile, certInfo)
	if err != nil {
		return nil, fmt.Errorf("unmarshal apinto.cert fail. %s", err)
	}

	err = b.examineCert(ctx, certInfo, machineCode)
	if err != nil {
		return nil, err
	}

	err = b.systemInfoStore.Transaction(ctx, func(txCtx context.Context) error {
		//保存证书
		cert := &system_entry.SystemInfo{
			Key:   enum.CertDBKey,
			Value: certFile,
		}
		err := b.systemInfoStore.Save(txCtx, cert)
		if err != nil {
			return err
		}
		certKey := b.systemInfoCache.Key(enum.CertDBKey)
		err = b.systemInfoCache.Set(txCtx, certKey, cert, time.Hour)
		if err != nil {
			return err
		}

		//保存机器码
		err = b.systemInfoStore.Save(txCtx, &system_entry.SystemInfo{
			Key:   enum.MachineCodeDBKey,
			Value: []byte(machineCode),
		})
		if err != nil {
			return err
		}

		//初始化dashboard_id
		err = b.systemInfoStore.InitDashboardID(txCtx, uuid.New())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	dashboardID := ""
	dashboardIdData, _ := b.systemInfoStore.GetSystemInfoByKey(ctx, enum.DashboardIdDBKey)
	if dashboardIdData != nil {
		dashboardID = string(dashboardIdData.Value)
	}

	return &model.ActivationInfo{
		BussinessCertInfo: *certInfo,
		DashboardID:       dashboardID,
	}, nil
}

func (b *bussinessAuthService) ReActivateCert(ctx context.Context, certFile []byte) (*model.ActivationInfo, error) {
	//校验证书
	machineCode := machine_code.GetMachineCode()
	certInfo := new(model.BussinessCertInfo)
	err := yaml.Unmarshal(certFile, certInfo)
	if err != nil {
		return nil, fmt.Errorf("unmarshal apinto.cert fail. %s", err)
	}

	err = b.examineCert(ctx, certInfo, machineCode)
	if err != nil {
		return nil, err
	}

	err = b.systemInfoStore.Transaction(ctx, func(txCtx context.Context) error {
		//保存证书
		cert := &system_entry.SystemInfo{
			Key:   enum.CertDBKey,
			Value: certFile,
		}
		err := b.systemInfoStore.Save(txCtx, cert)
		if err != nil {
			return err
		}

		certKey := b.systemInfoCache.Key(enum.CertDBKey)
		err = b.systemInfoCache.Set(txCtx, certKey, cert, time.Hour)
		if err != nil {
			return err
		}

		//保存机器码
		err = b.systemInfoStore.Save(txCtx, &system_entry.SystemInfo{
			Key:   enum.MachineCodeDBKey,
			Value: []byte(machineCode),
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	dashboardID := ""
	dashboardIdData, _ := b.systemInfoStore.GetSystemInfoByKey(ctx, enum.DashboardIdDBKey)
	if dashboardIdData != nil {
		dashboardID = string(dashboardIdData.Value)
	}

	return &model.ActivationInfo{
		BussinessCertInfo: *certInfo,
		DashboardID:       dashboardID,
	}, nil
}

func (b *bussinessAuthService) GetMachineCode(ctx context.Context) (string, error) {
	//未初次激活，或者证书过期的情况下才能返回
	info, err := b.systemInfoStore.GetSystemInfoByKey(ctx, enum.CertDBKey)
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}
	//若证书不为空，需要校验是否过期，过期才能返回
	if err == nil {
		certInfo := new(model.BussinessCertInfo)
		_ = yaml.Unmarshal(info.Value, certInfo)
		endTime := time.Unix(certInfo.EndTime, 0)
		//若现在未过期，则返回星号
		if time.Now().Before(endTime) {
			return "********", nil
		}
	}

	mac := machine_code.GetMachineCode()
	return mac, nil
}

// CheckCertValid 校验证书是否有效
func (b *bussinessAuthService) CheckCertValid(ctx context.Context) (bool, error) {
	return true, nil
	////获取证书信息
	//certInfo, err := b.getCertInfo(ctx)
	//if err != nil {
	//	if err == gorm.ErrRecordNotFound {
	//		return false, nil
	//	}
	//	return false, err
	//}
	////检查机器码. 即数据库有无修改过
	//machineCode := machine_code.GetMachineCode()
	//if certInfo.MachineCode != machineCode {
	//	return false, nil
	//}
	//
	////检查是否过期
	//now := time.Now()
	//beginTime := time.Unix(certInfo.BeginTime, 0)
	//endTime := time.Unix(certInfo.EndTime, 0)
	//if now.Before(beginTime) {
	//	return false, nil
	//}
	//return endTime.After(now), nil
}

func (b *bussinessAuthService) getCertInfo(ctx context.Context) (*model.BussinessCertInfo, error) {

	certCacheKey := b.systemInfoCache.Key(enum.CertDBKey)
	info, err := b.systemInfoCache.Get(ctx, certCacheKey)
	if err != nil {
		info, err = b.systemInfoStore.GetSystemInfoByKey(ctx, enum.CertDBKey)
		if err != nil {
			return nil, err
		}
		_ = b.systemInfoCache.Set(ctx, certCacheKey, info, time.Minute*30)
	}

	cert := new(model.BussinessCertInfo)
	_ = yaml.Unmarshal(info.Value, cert)

	return cert, nil
}

// ExamineCert 验证证书
func (b *bussinessAuthService) examineCert(ctx context.Context, certInfo *model.BussinessCertInfo, mac string) error {
	//校验机器码
	if mac != certInfo.MachineCode {
		return errors.New("machine_code isn't match. ")
	}
	//验签
	signParam := getSignParam(certInfo)

	sign, err := base64.StdEncoding.DecodeString(certInfo.Signature)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(publicKeyBytes)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	hash := sha256.New()
	hash.Write([]byte(signParam))
	err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, hash.Sum(nil), sign)
	if err != nil {
		return errors.New("verify signature fail. ")
	}

	now := time.Now()
	beginTime := time.Unix(certInfo.BeginTime, 0)
	endTime := time.Unix(certInfo.EndTime, 0)
	if now.Before(beginTime) {
		return fmt.Errorf("cert valid begin at %s. ", beginTime.Format(time.RFC3339))
	}
	if now.After(endTime) {
		return fmt.Errorf("cert is already out of date. ")
	}

	return nil
}

func getSignParam(certInfo *model.BussinessCertInfo) string {
	SignParam := &model.CertSignatureInfo{
		MachineCode:     certInfo.MachineCode,
		Company:         certInfo.Company,
		Edition:         certInfo.Edition,
		BeginTime:       certInfo.BeginTime,
		EndTime:         certInfo.EndTime,
		ControllerCount: certInfo.ControllerCount,
		NodeCount:       certInfo.NodeCount,
	}
	param, _ := yaml.Marshal(SignParam)
	return string(param)
}

var (
	publicKeyBytes = []byte(`-----BEGIN RSA Public Key-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDKWmkbciQ75lShbNBHWj94MFby
nrGOUI9Yk7NwirIeZAqdByoUIZmVuc1d1+9u6fSWPRHgXFzGI1OxExxi/ZGPgHFe
djbE+v5a1m66KdJSV9sOzF2YIQsXfxOdHGPPVxkdBrzVRvHnttdqq3bqIVEaymlm
8mr6ProjXdmjeBFIEwIDAQAB
-----END RSA Public Key-----
`)
)
