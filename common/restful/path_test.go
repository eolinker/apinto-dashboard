package restful

import "testing"

func Test__path_gen(t *testing.T) {

	type args struct {
		arg []string
	}
	tests := []struct {
		name string
		path string
		args args
		want string
	}{
		{
			name: "base",
			path: "/",
			args: args{},
			want: "/",
		}, {
			name: "none",
			path: "",
			args: args{},
			want: "/",
		},
		{
			name: "/path",
			path: "/path",
			args: args{},
			want: "/path",
		},
		{
			name: "path",
			path: "path",
			args: args{},
			want: "/path",
		},
		{
			name: "args",
			path: "path/:a/:b",
			args: args{
				arg: []string{"001", "b"},
			},
			want: "/path/001/b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := createPath(tt.path)
			if got := p.Gen(tt.args.arg...); got != tt.want {
				t.Errorf("gen() = %v, want %v", got, tt.want)
			}
		})
	}
}
