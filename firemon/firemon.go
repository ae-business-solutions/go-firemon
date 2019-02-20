package firemon

type Client struct {
	Username string
	Password string
	BaseURL  string
	Domain   int
}

func NewClient(baseurl, username, password string) *Client {
	return &Client{
		Username: username,
		Password: password,
		BaseURL:  baseurl,
		Domain:   1,
	}
}

func (c *Client) SetDomain(domain int) {
	c.Domain = domain
}
