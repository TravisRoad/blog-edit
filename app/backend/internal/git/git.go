package git

import (
	"os/exec"
)

type Client struct {
	Dir string
	Env []string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) WithDir(dir string) *Client {
	c.Dir = dir
	return c
}

func (c *Client) WithEnv(env map[string]string) *Client {
	for k, v := range env {
		c.Env = append(c.Env, k+"="+v)
	}
	return c
}

func (c *Client) Commit(msg string) error {
	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Dir = c.Dir
	cmd.Env = c.Env
	return cmd.Run()
}

func (c *Client) Add(filename string) error {
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = c.Dir
	cmd.Env = c.Env
	return cmd.Run()
}

func (c *Client) Push() error {
	cmd := exec.Command("git", "push")
	cmd.Dir = c.Dir
	cmd.Env = c.Env
	return cmd.Run()
}
