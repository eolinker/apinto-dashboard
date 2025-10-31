package cluster_dto

type ClusterCertificateInput struct {
	Key string `json:"key"`
	Pem string `json:"pem"`
}

type ClusterCertificateOut struct {
	Id           int      `json:"id"`
	ClusterId    int      `json:"cluster_id"`
	Name         string   `json:"name"`
	DNSName      []string `json:"dns_name"`
	OperatorName string   `json:"operator"`
	ValidTime    string   `json:"valid_time"`
	CreateTime   string   `json:"create_time"`
	UpdateTime   string   `json:"update_time"`
}

type GMCertificateInput struct {
	SignKey  string `json:"sign_key"`
	SignCert string `json:"sign_cert"`
	EncKey   string `json:"enc_key"`
	EncCert  string `json:"enc_cert"`
}

type GMCertificateOut struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	SignKey  string `json:"sign_key"`
	SignCert string `json:"sign_cert"`
	EncKey   string `json:"enc_key"`
	EncCert  string `json:"enc_cert"`
}

type GMCertificateItem struct {
	Id           int              `json:"id"`
	ClusterId    int              `json:"cluster_id"`
	Name         string           `json:"name"`
	SignCert     *CertificateInfo `json:"sign_cert"`
	EncCert      *CertificateInfo `json:"enc_cert"`
	OperatorName string           `json:"operator"`
	CreateTime   string           `json:"create_time"`
	UpdateTime   string           `json:"update_time"`
}

type CertificateInfo struct {
	DNSName   []string `json:"dns_name"`
	ValidTime string   `json:"valid_time"`
}
