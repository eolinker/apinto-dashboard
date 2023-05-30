package frontend_model

type RouterName string

var (
	RouterNameClusterVariable RouterName = "deploy/cluster/content"          //cluster_name
	RouterNameDiscoveryOnline RouterName = "upstream/serv-discovery/content" //discovery_name
	RouterNameServiceOnline   RouterName = "template/upstream"
	RouterNameTemplateOnline  RouterName = "router/plugin-template/content"
	RouterNameClusterPlugin   RouterName = "deploy/cluster/content/plugin"
)

type Router struct {
	Name   RouterName        `json:"name"`
	Params map[string]string `json:"params"`
	Msg    string            `json:"msg,omitempty"`
}
