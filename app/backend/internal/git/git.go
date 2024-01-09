package git

import (
	"fmt"
	"os/exec"
)

type Client struct {
	Dir    string
	Env    []string
	Author string
	Email  string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) WithDir(dir string) *Client {
	c.Dir = dir
	return c
}

func (c *Client) WithAuthor(author string) *Client {
	c.Author = author
	return c
}

func (c *Client) WithEmail(email string) *Client {
	c.Email = email
	return c
}

func (c *Client) WithEnv(env map[string]string) *Client {
	for k, v := range env {
		c.Env = append(c.Env, k+"="+v)
	}
	return c
}

func (c *Client) Commit(msg string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", msg, "--author", fmt.Sprintf("%s <%s>", c.Author, c.Email))
	cmd.Dir = c.Dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (c *Client) Add(filename string) (string, error) {
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = c.Dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (c *Client) Push() (string, error) {
	cmd := exec.Command("git", "push")
	cmd.Dir = c.Dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (c *Client) Pull() (string, error) {
	cmd := exec.Command("git", "pull", "--rebase")
	cmd.Dir = c.Dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}
