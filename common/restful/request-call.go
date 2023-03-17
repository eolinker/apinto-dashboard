package restful

type _RequestCallBuilder[O any] struct {
	builderCore
}

func (b *_RequestCallBuilder[O]) Build() RequestCall[O] {
	return &_RequestCall[O]{
		requestCore: b.build(),
	}
}

type _RequestCall[O any] struct {
	*requestCore
}

func (R *_RequestCall[O]) Header(name, value string) RequestCall[O] {
	R.requestCore.header.Set(name, value)
	return R
}

func (R *_RequestCall[O]) Query(name, value string) RequestCall[O] {
	R.requestCore.query.Set(name, value)
	return R
}
func (R *_RequestCall[O]) Request(args ...RestArgs) (*ResponseData[O], error) {

	resp, err := doRequest(R.address, R.method, R.path, args, R.query, R.header, nil)
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
