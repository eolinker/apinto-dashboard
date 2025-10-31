package common

import "context"

type _UserKeyType struct {
}

var (
	userKey = _UserKeyType{}
)

func WidthUser(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, userKey, id)
}
func GetUser(ctx context.Context) int {
	value := ctx.Value(userKey)
	if value == nil {
		return 0
	}
	if id, ok := value.(int); ok {
		return id
	}
	return 0
}
