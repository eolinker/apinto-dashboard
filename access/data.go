package access

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
)

//go:embed config/access.data
var accessTitleData []byte
var (
	accessAll          []Access
	accessTitles       map[Access]string
	accessKeys         map[Access]string
	accessParse        map[string]Access
	ErrorAccessUnknown = errors.New("unknown")
)

func All() []Access {
	return accessAll
}
func initData() {
	lines := bytes.Split(accessTitleData, []byte("\n"))
	l := len(lines)
	if l < int(lastId) {
		panic("access data less")
	}
	if l > int(lastId) {
		panic("access data over")
	}
	accessParse = make(map[string]Access, l)
	accessKeys = make(map[Access]string, l)
	accessTitles = make(map[Access]string, l)
	accessAll = make([]Access, l)
	for i, line := range lines {
		line = bytes.TrimSpace(line)
		index := bytes.IndexByte(line, ' ')
		if index < 1 {
			panic("access data invalid")
		}
		key := string(bytes.TrimSpace(line[:index]))
		title := string(bytes.TrimSpace(line[index:]))
		ac := Access(i)
		accessAll[i] = ac
		if _, has := accessParse[key]; has {
			panic("access key:" + key + "  duplicate")
		}
		accessParse[key] = ac
		accessTitles[ac] = title
		accessKeys[ac] = key
	}

}

func (a Access) String() string {
	return accessKeys[a]
}
func (a Access) Key() string {
	return accessKeys[a]
}
func (a Access) Title() string {
	return accessTitles[a]
}
func Parse(v string) (Access, error) {
	access, has := accessParse[v]
	if has {
		return access, nil
	}
	return unknown, fmt.Errorf("access %s:%w", v, ErrorAccessUnknown)
}
