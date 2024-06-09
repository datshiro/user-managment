package response

type Opts struct {
  Code int
}

func defaultOpts() Opts {
  return Opts{}
}

type OptFunc func(*Opts)

func WithCode(code int ) OptFunc  {
  return func(o *Opts) {
    o.Code = code
  }
}
