package apinto_module

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStaticRouter(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "noSub",
			args: args{
				prefix: "/prefix/*",
			},
			want: "/prefix/*filepath",
		}, {
			name: "sub",
			args: args{
				prefix: "/prefix/",
			},
			want: "/prefix/*filepath",
		}, {
			name: "noPrefix",
			args: args{
				prefix: "prefix",
			},
			want: "prefix",
		}, {
			name: "noPrefix",
			args: args{
				prefix: "prefix/",
			},
			want: "prefix/*filepath",
		}, {
			name: "empty",
			args: args{
				prefix: "/",
			},
			want: "/*filepath",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, StaticRouter(tt.args.prefix), "StaticRouter(%v)", tt.args.prefix)
		})
	}
}
