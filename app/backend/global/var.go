package global

import (
	"log/slog"
	"os"
	"sync"

	"github.com/TravisRoad/blog-edit/internal/git"
	"github.com/joho/godotenv"
)

var (
	Config Conf
	Git    *git.Client
)

var (
	once         sync.Once
	defaltClient *git.Client
)

func GitClient() *git.Client {
	if defaltClient == nil {
		once.Do(func() {
			defaltClient = git.NewClient().
				WithDir(Config.GitRepoPath).
				WithEnv(map[string]string{
					"GIT_AUTHOR_NAME":  Config.GitAuthorName,
					"GIT_AUTHOR_EMAIL": Config.GitAuthorEmail,
				}).
				WithAuthor(Config.GitAuthorName).
				WithEmail(Config.GitAuthorEmail)
		})
	}
	return defaltClient
}

const (
	Addr           = "ADDR"
	GitAuthorEmail = "GIT_AUTHOR_EMAIL"
	GitAuthorName  = "GIT_AUTHOR_NAME"
	GitRepoPath    = "GIT_REPO_PATH"
	DataPath       = "DATA_PATH"
	AuthToken      = "AUTH_TOKEN"
)

type Conf struct {
	Addr           string
	GitAuthorEmail string
	GitAuthorName  string
	GitRepoPath    string
	DataPath       string
	AuthToken      string
}

func defaultConfig() Conf {
	return Conf{
		Addr:           ":8080",
		GitAuthorEmail: "",
		GitAuthorName:  "",
		GitRepoPath:    "./repo",
		DataPath:       "./repo/data",
		AuthToken:      "",
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
		Config.GitRepoPath = os.Getenv(GitRepoPath)
		Config.DataPath = os.Getenv(DataPath)
		Config.AuthToken = os.Getenv(AuthToken)

		slog.Info("current", slog.Any("config", Config))
	}
	return nil
}
