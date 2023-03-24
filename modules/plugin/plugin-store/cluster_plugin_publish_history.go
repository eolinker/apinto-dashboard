package plugin_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/modules/base/publish-entry"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IClusterPluginPublishHistoryStore interface {
	store.BasePublishHistoryStore[plugin_entry.ClusterPluginPublishHistory]
}

type clusterPluginPublishHistoryHandler struct {
}

func (s *clusterPluginPublishHistoryHandler) Kind() string {
	return "p_h_cluster_plugin"
}

func (s *clusterPluginPublishHistoryHandler) Encode(sr *plugin_entry.ClusterPluginPublishHistory) *publish_entry.PublishHistory {
	val, _ := json.Marshal(sr.CluPluginPublishHistoryInfo)
	history := &publish_entry.PublishHistory{
		Id:          sr.Id,
		Kind:        s.Kind(),
		ClusterId:   sr.ClusterId,
		NamespaceId: sr.NamespaceId,
		Target:      sr.ClusterId,
		VersionId:   sr.VersionId,
		Data:        string(val),
		Desc:        sr.Desc,
		OptType:     sr.OptType,
		OptTime:     sr.OptTime,
		VersionName: sr.VersionName,
		Operator:    sr.Operator,
	}
	return history
}

func (s *clusterPluginPublishHistoryHandler) Decode(r *publish_entry.PublishHistory) *plugin_entry.ClusterPluginPublishHistory {
	val := new(plugin_entry.CluPluginPublishHistoryInfo)
	_ = json.Unmarshal([]byte(r.Data), val)
	history := &plugin_entry.ClusterPluginPublishHistory{
		Id:                          r.Id,
		VersionName:                 r.VersionName,
		Desc:                        r.Desc,
		NamespaceId:                 r.NamespaceId,
		ClusterId:                   r.ClusterId,
		VersionId:                   r.VersionId,
		OptTime:                     r.OptTime,
		OptType:                     r.OptType,
		CluPluginPublishHistoryInfo: *val,
		Operator:                    r.Operator,
	}
	return history
}

func newClusterPluginPublishHistoryStore(db store.IDB) IClusterPluginPublishHistoryStore {
	var historyHandler store.BaseKindHandler[plugin_entry.ClusterPluginPublishHistory, publish_entry.PublishHistory] = new(clusterPluginPublishHistoryHandler)
	return store.CreatePublishHistory[plugin_entry.ClusterPluginPublishHistory](historyHandler, db)
}
