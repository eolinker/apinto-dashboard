package enum

const (
	MachineCodeDBKey = "machine_code"
	CertDBKey        = "cert"
	DashboardIdDBKey = "dashboard_id"

	editionStandard = "standard"
	editionPremium  = "premium"
	editionUltimate = "ultimate"
)

var (
	editionTitles = map[string]string{
		editionStandard: "标准版",
		editionPremium:  "高级版",
		editionUltimate: "旗舰版",
	}
)

func EditionTitle(edition string) string {
	return editionTitles[edition]
}
