package server

type Server struct {
}

type Options struct {
}

func New(opts *Options) *Server {
	return &Server{}
}

func (s *Server) Run() error {
	return nil
}
