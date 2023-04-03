package enum

import "fmt"

type ClusterStatus int

const (
	ClusterStatusNone = iota
	ClusterStatusNormal
	ClusterStatusPartiallyNormal
	ClusterStatusAbnormal
	ClusterStatusAbnormalAll
)

var (
	clusterStatusNames = map[ClusterStatus]string{
		ClusterStatusNone:            "NONE",
		ClusterStatusNormal:          "NORMAL",
		ClusterStatusPartiallyNormal: "PARTIALLY_NORMAL",
		ClusterStatusAbnormal:        "ABNORMAL",
	}
	clusterStatusIndex = map[string]ClusterStatus{}
)

func (e ClusterStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", e.String(), "\"")), nil
}

func (e ClusterStatus) String() string {
	if e >= ClusterStatusAbnormalAll {
		return "unknown"
	}
	return clusterStatusNames[e]
}
