package enum

import "fmt"

type ClusterNodeStatus int

const (
	ClusterNodeStatusNone = iota
	ClusterNodeStatusNotRunning
	ClusterNodeStatusRunning
	ClusterNodeStatusAll
)

var (
	clusterNodeStatusNames = map[ClusterNodeStatus]string{
		ClusterNodeStatusNone:       "NONE",
		ClusterNodeStatusNotRunning: "NOTRUNNING", //未运行
		ClusterNodeStatusRunning:    "RUNNING",    //运行中
	}
)

func (e ClusterNodeStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", e.String(), "\"")), nil
}

func (e ClusterNodeStatus) String() string {
	if e >= ClusterNodeStatusAll {
		return "unknown"
	}
	return clusterNodeStatusNames[e]
}
