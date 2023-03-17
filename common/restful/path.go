package restful

import (
	"net/url"
	"strings"
)

type PathGen interface {
	Gen(arg ...string) string
}
type _path struct {
	paths   []string
	size    int
	index   []int
	argSize int
}

func (p *_path) Gen(arg ...string) string {

	nodes := make([]string, p.size)
	copy(nodes, p.paths)
	size := p.argSize
	if len(arg) < size {
		size = len(arg)
	}
	for i := 0; i < size; i++ {
		v := url.QueryEscape(arg[i])
		nodes[p.index[i]] = v
	}
	return strings.Join(nodes, "/")
}

/*
createPath restful参数构造器， 格式为 /:arg
调用gen 方法构造出实际的 path，gen方法如果传入的参数数量少于需要的值，会使用参数名作为默认值
传入的参数会进行 QueryEscape
*/
func createPath(path string) PathGen {
	offset := 0
	if !strings.HasPrefix(path, "/") {
		offset = 1
	}
	ps := strings.Split(path, "/")
	size := len(ps) + offset
	paths := make([]string, offset, size)

	index := make([]int, 0, size)
	for i, p := range ps {

		if strings.HasPrefix(p, ":") {
			p = p[1:]
			index = append(index, i+offset)
		} else {
			l := len(p)
			if l > 1 && p[0] == '{' && p[l-1] == '}' {
				p = p[1 : l-1]
				index = append(index, i+offset)
			}
		}

		paths = append(paths, p)
	}
	return &_path{paths: paths, size: size, index: index, argSize: len(index)}
}
