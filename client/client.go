package client

type Client struct {
}

type Options struct {
	Args string
}

func New(opts *Options) *Client {
	return &Client{}
}

func (c *Client) Run() error {
	return nil
}
