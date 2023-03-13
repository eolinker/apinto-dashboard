package api_store

import (
	"context"
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
	api_entry "github.com/eolinker/apinto-dashboard/modules/api/api-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IAPIVersionStore interface {
	store.IBaseStore[api_entry.APIVersion]
	GetAPIVersionByApiIds(ctx context.Context, ids []int) ([]*api_entry.APIVersion, error)
}

type apiVersionStore struct {
	*store.BaseKindStore[api_entry.APIVersion, entry.Version]
}

type apiVersionKindHandler struct {
}

func (s *apiVersionKindHandler) Kind() string {
	return "api"
}

func (s *apiVersionKindHandler) Encode(av *api_entry.APIVersion) *entry.Version {
	data, _ := json.Marshal(av.APIVersionConfig)
	v := &entry.Version{
		Id:          av.Id,
		Target:      av.ApiID,
		NamespaceID: av.NamespaceID,
		Kind:        s.Kind(),
		Data:        data,
		Operator:    av.Operator,
		CreateTime:  av.CreateTime,
	}

	return v
}

func (s *apiVersionKindHandler) Decode(v *entry.Version) *api_entry.APIVersion {
	av := &api_entry.APIVersion{
		Id:               v.Id,
		ApiID:            v.Target,
		NamespaceID:      v.NamespaceID,
		APIVersionConfig: api_entry.APIVersionConfig{},
		Operator:         v.Operator,
		CreateTime:       v.CreateTime,
	}
	_ = json.Unmarshal(v.Data, &av.APIVersionConfig)

	return av
}

func (s *apiVersionStore) GetAPIVersionByApiIds(ctx context.Context, ids []int) ([]*api_entry.APIVersion, error) {
	versions := make([]*entry.Version, 0, len(ids))
	results := make([]*api_entry.APIVersion, 0, len(ids))

	err := s.DB(ctx).Model(&entry.Version{}).
		Select(`version.*`).
		Joins("right join stat on stat.version = version.id").
		Where("stat.target in (?) and stat.kind = ?", ids, s.BaseKindStore.Kind()).
		Find(&versions).Error

	if err != nil {
		return results, err
	}

	for _, version := range versions {
		s.BaseKindStore.Decode(version)
		results = append(results, s.BaseKindStore.Decode(version))
	}

	return results, nil
}

func newAPIVersionStore(db store.IDB) IAPIVersionStore {
	var h store.BaseKindHandler[api_entry.APIVersion, entry.Version] = new(apiVersionKindHandler)
	return &apiVersionStore{store.CreateBaseKindStore(h, db)}
}
