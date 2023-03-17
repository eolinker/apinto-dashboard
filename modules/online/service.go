package online

import "context"

type IResetOnlineService interface {
	ResetOnline(ctx context.Context, namespaceId, clusterId int)
}
