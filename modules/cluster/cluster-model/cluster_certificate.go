package cluster_model

import (
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
)

type ClusterCertificate struct {
	*cluster_entry.ClusterCertificate
	ClusterName  string
	OperatorName string
}

type ClusterGMCertificate struct {
	*cluster_entry.ClusterGMCertificate
	ClusterName  string
	OperatorName string
}

type Certificate struct {
	ID        int      `json:"id"`
	Uuid      string   `json:"uuid"`
	Name      string   `json:"name"`
	DnsName   []string `json:"dns_name"`
	Key       string   `json:"key"`
	Pem       string   `json:"pem"`
	ValidTime string   `json:"valid_time"`
}

type GMCertificate struct {
	ID       int    `json:"id"`
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	SignKey  string `json:"sign_key"`
	SignCert string `json:"sign_cert"`
	EncKey   string `json:"enc_key"`
	EncCert  string `json:"enc_cert"`
}
