package frontend_model

type RouterName string

var (
	RouterNameClusterVariable RouterName = "deploy/cluster/content"          //cluster_name
	RouterNameDiscoveryOnline RouterName = "upstream/serv-discovery/content" //discovery_name
	RouterNameServiceOnline   RouterName = "upstream/upstream/content/publish"
	RouterNameTemplateOnline  RouterName = "router/plugin/content"
)

type Router struct {
	Name   RouterName        `json:"name"`
	Params map[string]string `json:"params"`
}
