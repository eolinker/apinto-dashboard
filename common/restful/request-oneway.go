package restful

import "encoding/json"

type _OnewayBuilder[I any] struct {
	builderCore
}

func (b *_OnewayBuilder[I]) Build() RequestOneWay[I] {
	return &_OnewayRequest[I]{
		requestCore: b.build(),
	}
}

type _OnewayRequest[I any] struct {
	*requestCore
}

func (R *_OnewayRequest[I]) Header(name, value string) RequestOneWay[I] {
	R.requestCore.header.Set(name, value)
	return R
}

func (R *_OnewayRequest[I]) Query(name, value string) RequestOneWay[I] {
	R.requestCore.query.Set(name, value)
	return R
}
func (R *_OnewayRequest[I]) Request(input *I, args ...RestArgs) (*Response, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	resp, err := doRequest(R.address, R.method, R.path, args, R.query, R.header, data)
	if err != nil {
		return nil, err
	}
	out := new(Response)
	err = UnmarshalResponse(resp, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
