package git

import "os/exec"

type Client struct {
	Dir string
	Env map[string]string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) WithDir(dir string) *Client {
	c.Dir = dir
	return c
}

func (c *Client) WithEnv(env map[string]string) *Client {
	c.Env = env
	return c
}

func (c *Client) Commit(msg string) error {
	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Dir = c.Dir
	return cmd.Run()
}

func (c *Client) Add(filename string) error {
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = c.Dir
	return cmd.Run()
}

func (c *Client) Push() error {
	cmd := exec.Command("git", "push")
	cmd.Dir = c.Dir
	return cmd.Run()
}
