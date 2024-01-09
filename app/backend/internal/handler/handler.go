package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/TravisRoad/blog-edit/global"
	"github.com/gorilla/mux"
)

func fail(w http.ResponseWriter, msg string, statusCode int) {
	_, _ = w.Write([]byte(msg))
	w.WriteHeader(http.StatusInternalServerError)
}

func decodeFilename(r *http.Request) (string, error) {
	enc := mux.Vars(r)["filename"]
	filename, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}
	path := string(filename)
	// path := filepath.Join(global.Config.DataPath, string(filename))
	return filepath.Abs(path)
}

func GetFileList(w http.ResponseWriter, r *http.Request) {
	ans := []string{}
	err := filepath.WalkDir(global.Config.DataPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ans = append(ans, path)
		return nil
	})

	if err != nil {
		fail(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(map[string]interface{}{
		"code": 0,
		"msg":  "",
		"data": ans,
	})
	if err != nil {
		fail(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		fail(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	path, err := decodeFilename(r)
	if err != nil {
		fail(w, err.Error(), http.StatusInternalServerError)
		return
	}
	content, err := os.ReadFile(path)
	if err != nil {
		fail(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(content)
	if err != nil {
		fail(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateFile(w http.ResponseWriter, r *http.Request) {
	path, err := decodeFilename(r)
	if err != nil {
		fail(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		fail(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		fail(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Apply(w http.ResponseWriter, r *http.Request) {
	path, err := decodeFilename(r)
	if err != nil {
		fail(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cli := global.GitClient()

	if out, err := cli.Add(path); err != nil {
		slog.Error(out, err)
		fail(w, out, http.StatusInternalServerError)
		return
	}
	if out, err := cli.Commit("apply change"); err != nil {
		slog.Error(out, err)
		fail(w, out, http.StatusInternalServerError)
		return
	}
}

func Sync(w http.ResponseWriter, r *http.Request) {
	cli := global.GitClient()
	if out, err := cli.Pull(); err != nil {
		slog.Error(out, err)
		fail(w, out, http.StatusInternalServerError)
		return
	}
	if out, err := cli.Push(); err != nil {
		slog.Error(out, err)
		fail(w, fmt.Sprintf("failed to push: %s", out), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
