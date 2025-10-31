package cluster_service

import (
	"context"
	"github.com/eolinker/eosc/common/bean"
	"time"

	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-store"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/go-basic/uuid"
)

func newGmClusterCertificateService() cluster.IGmCertificateService {
	s := &gmCertificateService{}
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.clusterCertificateStore)
	bean.Autowired(&s.userInfoService)
	bean.Autowired(&s.apintoClient)

	return s
}

type gmCertificateService struct {
	clusterCertificateStore cluster_store.IGMCertificateStore
	clusterService          cluster.IClusterService
	userInfoService         user.IUserInfoService
	apintoClient            cluster.IApintoClient
}

func (c *gmCertificateService) Info(ctx context.Context, namespaceId, certificateId int, clusterName string) (*cluster_model.GMCertificate, error) {
	clusterId, err := c.clusterService.CheckByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, err
	}
	info, err := c.clusterCertificateStore.First(ctx, map[string]interface{}{
		"namespace": namespaceId,
		"id":        certificateId,
		"cluster":   clusterId,
	})
	if err != nil {
		return nil, err
	}
	signCert, err := common.ParseGMCert(info.SignKey, info.SignCert)
	if err != nil {
		return nil, err
	}

	return &cluster_model.GMCertificate{
		ID:       info.Id,
		Uuid:     info.Uuid,
		Name:     signCert.Leaf.Subject.CommonName,
		SignKey:  common.Base64Encode([]byte(info.SignKey)),
		SignCert: common.Base64Encode([]byte(info.SignCert)),
		EncKey:   common.Base64Encode([]byte(info.EncKey)),
		EncCert:  common.Base64Encode([]byte(info.EncCert)),
	}, nil
}

// Insert 新增证书
func (c *gmCertificateService) Insert(ctx context.Context, operator, namespaceId int, clusterName, signKey, signCert, encKey, encCert string) error {
	t := time.Now()
	clusterId, err := c.clusterService.CheckByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	client, err := c.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		return err
	}

	//校验证书是否合法
	sc, err := common.ParseGMCert(signKey, signCert)
	if err != nil {
		return err
	}

	//校验证书是否合法
	_, err = common.ParseGMCert(encKey, encCert)
	if err != nil {
		return err
	}

	value := &cluster_entry.ClusterGMCertificate{
		Uuid:        uuid.New(),
		ClusterId:   clusterId,
		NamespaceId: namespaceId,
		Operator:    operator,
		SignKey:     signKey,
		SignCert:    signCert,
		EncKey:      encKey,
		EncCert:     encCert,
		CreateTime:  t,
		UpdateTime:  t,
	}

	return c.clusterCertificateStore.Transaction(ctx, func(txCtx context.Context) error {

		if err = c.clusterCertificateStore.Insert(txCtx, value); err != nil {
			return err
		}

		config := v1.GMCertConfig{
			Name:        value.Uuid,
			SignCert:    common.Base64Encode([]byte(signCert)),
			SignKey:     common.Base64Encode([]byte(signKey)),
			EncCert:     common.Base64Encode([]byte(encCert)),
			EncKey:      common.Base64Encode([]byte(encKey)),
			Description: sc.Leaf.Subject.CommonName,
			Driver:      "gm-server",
		}
		return client.ForGMCert().Save(value.Uuid, &config)
	})

}

// Update 修改证书
func (c *gmCertificateService) Update(ctx context.Context, operator, namespaceId, certificateId int, clusterName, signKey, signCert, encKey, encCert string) error {

	clusterId, err := c.clusterService.CheckByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	client, err := c.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		return err
	}

	//校验证书是否合法
	sc, err := common.ParseGMCert(signKey, signCert)
	if err != nil {
		return err
	}

	//校验证书是否合法
	_, err = common.ParseGMCert(encKey, encCert)
	if err != nil {
		return err
	}

	//查询证书
	cc, err := c.clusterCertificateStore.Get(ctx, certificateId)
	if err != nil {
		return err
	}

	var value = &cluster_entry.ClusterGMCertificate{
		Id:          certificateId,
		Uuid:        cc.Uuid,
		ClusterId:   clusterId,
		NamespaceId: namespaceId,
		Operator:    operator,
		SignKey:     signKey,
		SignCert:    signCert,
		EncKey:      encKey,
		EncCert:     encCert,
		UpdateTime:  time.Now(),
	}
	return c.clusterCertificateStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = c.clusterCertificateStore.Update(txCtx, value); err != nil {
			return err
		}

		config := v1.GMCertConfig{
			Name:        cc.Uuid,
			SignCert:    common.Base64Encode([]byte(signCert)),
			SignKey:     common.Base64Encode([]byte(signKey)),
			EncCert:     common.Base64Encode([]byte(encCert)),
			EncKey:      common.Base64Encode([]byte(encKey)),
			Description: sc.Leaf.Subject.CommonName,
			Driver:      "gm-server",
		}
		return client.ForGMCert().Save(cc.Uuid, &config)
	})

}

// QueryList 查询证书列表
func (c *gmCertificateService) QueryList(ctx context.Context, namespaceId int, clusterName string) ([]*cluster_model.ClusterGMCertificate, error) {

	clusterId, err := c.clusterService.CheckByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, err
	}

	list, err := c.clusterCertificateStore.QueryListByClusterId(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	userIds := common.SliceToSliceIds(list, func(t *cluster_entry.ClusterGMCertificate) int {
		return t.Operator
	})

	userInfoMaps, _ := c.userInfoService.GetUserInfoMaps(ctx, userIds...)

	modelList := make([]*cluster_model.ClusterGMCertificate, 0, len(list))
	for _, certificate := range list {

		operatorName := ""
		if userInfo, ok := userInfoMaps[certificate.Operator]; ok {
			operatorName = userInfo.NickName
		}

		modelList = append(modelList, &cluster_model.ClusterGMCertificate{ClusterGMCertificate: certificate, OperatorName: operatorName})
	}
	return modelList, nil
}

// DeleteById 删除证书
func (c *gmCertificateService) DeleteById(ctx context.Context, namespaceId int, clusterName string, id int) error {

	clusterId, err := c.clusterService.CheckByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	client, err := c.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		return err
	}

	//查询证书
	cc, err := c.clusterCertificateStore.Get(ctx, id)
	if err != nil {
		return err
	}

	//todo 删除前的逻辑校验
	return c.clusterCertificateStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = c.clusterCertificateStore.Delete(txCtx, id); err != nil {
			return err
		}

		return client.ForGMCert().Del(cc.Uuid)

	})

}
