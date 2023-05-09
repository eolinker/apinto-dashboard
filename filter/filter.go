package filter

import (
	"github.com/gin-gonic/gin"
	"sort"
	"strings"
)

type Builder interface {
	Exclude(path ...string) Builder
	Root(path ...string) Builder
	Build(hs ...gin.HandlerFunc) gin.HandlerFunc
}
type _Builder struct {
	excludes []string
	roots    []string
}
type PathSort []string

func (ps PathSort) Len() int {
	return len(ps)
}

func (ps PathSort) Less(i, j int) bool {
	if len(ps[i]) == len(ps[j]) {
		return ps[i] < ps[j]
	}
	return len(ps[i]) > len(ps[j])
}

func (ps PathSort) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (b *_Builder) Build(hs ...gin.HandlerFunc) gin.HandlerFunc {
	tmp := make(map[string]struct{})

	for _, r := range b.roots {
		tmp[r] = struct{}{}
	}
	roots := make([]string, 0, len(tmp))
	for r := range tmp {
		roots = append(roots, r)
	}
	sort.Sort(PathSort(roots))
	excludes := make(map[string]struct{})
	for _, e := range b.excludes {
		excludes[e] = struct{}{}
	}

	return func(ginCtx *gin.Context) {
		path := ginCtx.FullPath()
		hasRoot := false
		// 这里以后优化
		for _, root := range roots {
			if hasRoot = strings.HasPrefix(path, root); hasRoot {
				break
			}
		}
		//

		if !hasRoot {
			return
		}

		if _, ok := excludes[path]; ok {
			return
		}

		for _, h := range hs {
			h(ginCtx)
		}

	}
}

func (b *_Builder) Exclude(path ...string) Builder {
	b.excludes = append(b.excludes, path...)
	return b
}

func (b *_Builder) Root(path ...string) Builder {
	b.roots = append(b.roots, path...)
	return b
}

func NewBuilder() Builder {
	return &_Builder{}
}
