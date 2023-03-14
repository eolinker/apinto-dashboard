package cluster_service

import (
	"context"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-store"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-basic/uuid"
	"time"
)

type clusterCertificateService struct {
	clusterCertificateStore cluster_store.IClusterCertificateStore
	clusterService          cluster.IClusterService
	userInfoService         user.IUserInfoService
	apintoClient            cluster.IApintoClient
}

func newClusterCertificateService() cluster.IClusterCertificateService {
	s := &clusterCertificateService{}
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.clusterCertificateStore)
	bean.Autowired(&s.userInfoService)
	bean.Autowired(&s.apintoClient)

	return s
}

// Insert 新增证书
func (c *clusterCertificateService) Insert(ctx context.Context, operator, namespaceId int, clusterName, key, pem string) error {
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
	_, err = common.ParseCert(key, pem)
	if err != nil {
		return err
	}

	value := &cluster_entry.ClusterCertificate{
		ClusterId:   clusterId,
		NamespaceId: namespaceId,
		Operator:    operator,
		Key:         key,
		Pem:         pem,
		UUID:        uuid.New(),
		CreateTime:  t,
		UpdateTime:  t,
	}

	return c.clusterCertificateStore.Transaction(ctx, func(txCtx context.Context) error {

		if err = c.clusterCertificateStore.Insert(txCtx, value); err != nil {
			return err
		}

		config := v1.CertConfig{
			Name:   value.UUID,
			Key:    common.Base64Encode([]byte(key)),
			Pem:    common.Base64Encode([]byte(pem)),
			Driver: "server",
		}
		return client.ForCert().Save(value.UUID, &config)
	})

}

// Update 修改证书
func (c *clusterCertificateService) Update(ctx context.Context, operator, namespaceId, certificateId int, clusterName, key, pem string) error {

	clusterId, err := c.clusterService.CheckByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	client, err := c.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		return err
	}

	//校验证书是否合法
	_, err = common.ParseCert(key, pem)
	if err != nil {
		return err
	}

	//查询证书
	cc, err := c.clusterCertificateStore.Get(ctx, certificateId)
	if err != nil {
		return err
	}

	value := &cluster_entry.ClusterCertificate{
		Id:          certificateId,
		ClusterId:   clusterId,
		NamespaceId: namespaceId,
		Operator:    operator,
		Key:         key,
		Pem:         pem,
		UUID:        cc.UUID,
		UpdateTime:  time.Now(),
	}
	return c.clusterCertificateStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = c.clusterCertificateStore.Update(txCtx, value); err != nil {
			return err
		}

		config := v1.CertConfig{
			Name:   cc.UUID,
			Key:    common.Base64Encode([]byte(key)),
			Pem:    common.Base64Encode([]byte(pem)),
			Driver: "server",
		}
		return client.ForCert().Save(cc.UUID, &config)
	})

}

// QueryList 查询证书列表
func (c *clusterCertificateService) QueryList(ctx context.Context, namespaceId int, clusterName string) ([]*cluster_model.ClusterCertificate, error) {

	clusterId, err := c.clusterService.CheckByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, err
	}

	list, err := c.clusterCertificateStore.QueryListByClusterId(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	userIds := common.SliceToSliceIds(list, func(t *cluster_entry.ClusterCertificate) int {
		return t.Operator
	})

	userInfoMaps, _ := c.userInfoService.GetUserInfoMaps(ctx, userIds...)

	modelList := make([]*cluster_model.ClusterCertificate, 0, len(list))
	for _, certificate := range list {

		operatorName := ""
		if userInfo, ok := userInfoMaps[certificate.Operator]; ok {
			operatorName = userInfo.NickName
		}

		modelList = append(modelList, &cluster_model.ClusterCertificate{ClusterCertificate: certificate, OperatorName: operatorName})
	}
	return modelList, nil
}

// DeleteById 删除证书
func (c *clusterCertificateService) DeleteById(ctx context.Context, namespaceId int, clusterName string, id int) error {

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

		return client.ForCert().Del(cc.UUID)

	})

}
