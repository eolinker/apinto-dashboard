package plugin_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/modules/base/version-entry"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IClusterPluginPublishVersionStore interface {
	store.IBaseStore[plugin_entry.ClusterPluginPublishVersion]
}

type clusterPluginPublishVersionHandler struct {
}

func (s *clusterPluginPublishVersionHandler) Kind() string {
	return "p_v_cluster_plugin"
}

func (s *clusterPluginPublishVersionHandler) Encode(sv *plugin_entry.ClusterPluginPublishVersion) *version_entry.Version {

	v := new(version_entry.Version)
	v.Id = sv.Id
	v.Kind = s.Kind()
	v.Target = sv.ClusterId
	v.Operator = sv.Operator
	v.CreateTime = sv.CreateTime
	v.NamespaceID = sv.NamespaceId
	bytes, _ := json.Marshal(sv.PublishedPluginsList)
	v.Data = bytes

	return v
}

func (s *clusterPluginPublishVersionHandler) Decode(v *version_entry.Version) *plugin_entry.ClusterPluginPublishVersion {
	sv := new(plugin_entry.ClusterPluginPublishVersion)
	sv.Id = v.Id
	sv.ClusterId = v.Target
	sv.Operator = v.Operator
	sv.NamespaceId = v.NamespaceID
	sv.CreateTime = v.CreateTime
	_ = json.Unmarshal(v.Data, &sv.PublishedPluginsList)

	return sv
}

func newClusterPluginPublishVersionStore(db store.IDB) IClusterPluginPublishVersionStore {
	var h store.BaseKindHandler[plugin_entry.ClusterPluginPublishVersion, version_entry.Version] = &clusterPluginPublishVersionHandler{}
	return store.CreateBaseKindStore(h, db)
}
