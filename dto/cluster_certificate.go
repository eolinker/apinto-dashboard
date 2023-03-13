package dto

type ClusterCertificateInput struct {
	Key string `json:"key"`
	Pem string `json:"pem"`
}

type ClusterCertificateOut struct {
	Id           int    `json:"id"`
	ClusterId    int    `json:"cluster_id"`
	Name         string `json:"name"`
	OperatorName string `json:"operator"`
	ValidTime    string `json:"valid_time"`
	CreateTime   string `json:"create_time"`
	UpdateTime   string `json:"update_time"`
}
