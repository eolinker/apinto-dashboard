package navigation_model

type Navigation struct {
	Uuid  string `json:"uuid"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
	Index int    `json:"-"`
}

type Navigations []*Navigation

func (n Navigations) Len() int {
	return len(n)
}

func (n Navigations) Less(i, j int) bool {
	return n[i].Index < n[j].Index
}

func (n Navigations) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
