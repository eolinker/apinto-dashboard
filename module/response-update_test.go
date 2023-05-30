package apinto_module

import (
	"fmt"
	_ "github.com/stretchr/testify/assert"
)

func ExampleDoEvent() {
	type args struct {
		event string
		v     any
	}
	type testCase struct {
		name string
		args args
	}
	tests := []testCase{
		{
			name: "login",
			args: args{
				event: "login",
				v:     "username=a&nickname=b",
			},
		},
	}
	RegisterEventHandler("login", func(event string, v any) {
		fmt.Printf("get event %s:%s", event, v)
	})
	for _, tt := range tests {

		DoEvent(tt.args.event, tt.args.v)

	}
	//output: get event login:username=a&nickname=b
}
