package store

type Handler func(db IDB)

var (
	handlers []Handler
)

func RegisterStore(h Handler) {
	handlers = append(handlers, h)
}
func runHandler(db IDB) {
	for _, h := range handlers {
		h(db)
	}
}
