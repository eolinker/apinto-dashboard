package dto

type SystemRunInfo struct {
	Title       string `json:"title"`
	DashboardID string `json:"dashboard_id"`
	Version     string `json:"version"`
}

type ActivationInfoItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
