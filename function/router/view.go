package router

type View struct {
}

func (v *View) Get(route *Route) *Response {
	return NewResponse(GenericBody{}, "500")
}

func (v *View) Post(route *Route) *Response {
	return NewResponse(GenericBody{}, "500")
}

func (v *View) Put(route *Route) *Response {
	return NewResponse(GenericBody{}, "500")
}

func (v *View) Patch(route *Route) *Response {
	return NewResponse(GenericBody{}, "500")
}

func (v *View) Delete(route *Route) *Response {
	return NewResponse(GenericBody{}, "500")
}
