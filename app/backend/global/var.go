package global

import (
	"log/slog"
	"os"

	"github.com/TravisRoad/blog-edit/internal/git"
	"github.com/joho/godotenv"
)

var (
	Config Conf
	Git    *git.Client
)

const (
	Addr           = "ADDR"
	GitAuthorEmail = "GIT_AUTHOR_EMAIL"
	GitAuthorName  = "GIT_AUTHOR_NAME"
)

type Conf struct {
	Addr           string
	GitAuthorEmail string
	GitAuthorName  string
}

func defaultConfig() Conf {
	return Conf{
		Addr:           ":8080",
		GitAuthorEmail: "",
		GitAuthorName:  "",
	}
}

func LoadConfig() error {
	Config = defaultConfig()
	slog.Info("loading .env file...")
	if err := godotenv.Load(".env"); err != nil {
		slog.Info("Error loading .env file", "Using default", err)
		slog.Info("current", slog.Any("config", Config))
		return nil
	}
	{
		Config.Addr = os.Getenv(Addr)
		Config.GitAuthorEmail = os.Getenv(GitAuthorEmail)
		Config.GitAuthorName = os.Getenv(GitAuthorName)
		slog.Info("current", slog.Any("config", Config))
	}
	return nil
}
