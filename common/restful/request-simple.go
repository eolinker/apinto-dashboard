package restful

type _SimpleBuilder struct {
	builderCore
}

func (b *_SimpleBuilder) Build() Request {
	return &_RequestSimple{
		requestCore: b.build(),
	}
}

type _RequestSimple struct {
	*requestCore
}

func (r *_RequestSimple) Header(name, value string) Request {
	r.requestCore.header.Set(name, value)
	return r
}

func (r *_RequestSimple) Query(name, value string) Request {
	r.requestCore.query.Set(name, value)
	return r
}

func (r *_RequestSimple) Request(args ...RestArgs) (*Response, error) {
	response, err := doRequest(r.address, r.method, r.path, args, r.query, r.header, nil)
	if err != nil {
		return nil, err
	}

	resp := new(Response)
	err = UnmarshalResponse(response, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
