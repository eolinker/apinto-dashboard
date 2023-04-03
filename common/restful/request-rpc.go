package restful

import "encoding/json"

type _RpcBuilder[I any, O any] struct {
	builderCore
}

func (b *_RpcBuilder[I, O]) Build() RequestRPC[I, O] {
	return &RPCRequest[I, O]{
		requestCore: b.build(),
	}
}

type RPCRequest[I any, O any] struct {
	*requestCore
}

func (R *RPCRequest[I, O]) Header(name, value string) RequestRPC[I, O] {
	R.requestCore.header.Set(name, value)
	return R
}

func (R *RPCRequest[I, O]) Query(name, value string) RequestRPC[I, O] {
	R.requestCore.query.Set(name, value)
	return R
}
func (R *RPCRequest[I, O]) Request(input *I, args ...RestArgs) (*ResponseData[O], error) {
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	resp, err := doRequest(R.address, R.method, R.path, args, R.query, R.header, data)
	if err != nil {
		return nil, err
	}
	out := new(ResponseData[O])
	err = UnmarshalResponse(resp, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
