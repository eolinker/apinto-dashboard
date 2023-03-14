package bussiness_model

type BussinessCertInfo struct {
	MachineCode     string `json:"machine_code" yaml:"machine_code"`
	Company         string `json:"company" yaml:"company"`
	Edition         string `json:"edition" yaml:"edition"`
	BeginTime       int64  `json:"begin_time" yaml:"begin_time"`
	EndTime         int64  `json:"end_time" yaml:"end_time"`
	ControllerCount int    `json:"controller_count" yaml:"controller_count"`
	NodeCount       int    `json:"node_count" yaml:"node_count"`
	Signature       string `json:"signature,omitempty" yaml:"signature,omitempty"`
}

type ActivationInfo struct {
	BussinessCertInfo
	DashboardID string
}

type CertSignatureInfo struct {
	MachineCode     string `json:"machine_code" yaml:"machine_code"`
	Company         string `json:"company" yaml:"company"`
	Edition         string `json:"edition" yaml:"edition"`
	BeginTime       int64  `json:"begin_time" yaml:"begin_time"`
	EndTime         int64  `json:"end_time" yaml:"end_time"`
	ControllerCount int    `json:"controller_count" yaml:"controller_count"`
	NodeCount       int    `json:"node_count" yaml:"node_count"`
}
