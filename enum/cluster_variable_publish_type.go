package enum

import "fmt"

type ClusterVariablePublish int

const (
	ClusterVariablePublishNone = iota
	ClusterVariablePublishUnPublished
	ClusterVariablePublishPublished
	ClusterVariablePublishDefect
	ClusterVariablePublishAll
)

var (
	clusterVariablePublishNames = map[ClusterVariablePublish]string{
		ClusterVariablePublishNone:        "NONE",
		ClusterVariablePublishUnPublished: "UNPUBLISHED", //未发布
		ClusterVariablePublishPublished:   "PUBLISHED",   //已发布
		ClusterVariablePublishDefect:      "DEFECT",      //缺失
	}
	ClusterVariablePublishIndex = map[string]ClusterVariablePublish{
		"UNPUBLISHED": ClusterVariablePublishUnPublished,
		"PUBLISHED":   ClusterVariablePublishPublished,
		"DEFECT":      ClusterVariablePublishDefect,
	}
)

func (e ClusterVariablePublish) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", e.String(), "\"")), nil
}

func (e ClusterVariablePublish) String() string {
	if e >= ClusterVariablePublishAll {
		return "unknown"
	}
	return clusterVariablePublishNames[e]
}
