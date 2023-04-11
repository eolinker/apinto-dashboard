package navigation_model

type NavigationBasicInfo struct {
	ID        int    `json:"-"`
	Uuid      string `json:"uuid"`
	Title     string `json:"title"`
	Icon      string `json:"icon"`
	IconType  string `json:"icon_type"`
	CanDelete bool   `json:"can_delete"`
	Sort      int    `json:"-"`
}

type Navigations []*NavigationBasicInfo

func (n Navigations) Len() int {
	return len(n)
}

func (n Navigations) Less(i, j int) bool {
	return n[i].Sort < n[j].Sort
}

func (n Navigations) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

type Navigation struct {
	*NavigationBasicInfo
	Modules []*Module `json:"modules"`
}

type Module struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}
