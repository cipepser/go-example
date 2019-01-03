package functionalOptions

type RequestOptions struct {
	page    int
	perPage int
	sort    string
}

type Option func(request *RequestOptions)

func Page(p int) Option {
	return func(r *RequestOptions) {
		if r != nil {
			r.page = p
		}
	}
}

func PerPage(pp int) Option {
	return func(r *RequestOptions) {
		if r != nil {
			r.perPage = pp
		}
	}
}

func Sort(s string) Option {
	return func(r *RequestOptions) {
		if r != nil {
			r.sort = s
		}
	}
}

func NewRequest(opts ...Option) *RequestOptions {
	// set default values
	r := &RequestOptions{
		page:    1,
		perPage: 30,
		sort:    "desc",
	}

	for _, opt := range opts {
		opt(r)
	}
	return r
}
