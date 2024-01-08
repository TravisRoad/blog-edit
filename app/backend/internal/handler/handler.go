package handler

import (
	"encoding/base64"
	"encoding/json"
	"io"
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
	path := filepath.Join(global.Config.DataPath, string(filename))
	return path, nil
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
		ans = append(ans, d.Name())
		return nil
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(map[string]interface{}{
		"code": 0,
		"msg":  "",
		"data": ans,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	enc := mux.Vars(r)["filename"]
	filename, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	path := filepath.Join(global.Config.DataPath, string(filename))
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
	// path, err := decodeFilename(r)
	// if err != nil {
	// 	fail(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// cli := global.GitClient()
}

func Sync(w http.ResponseWriter, r *http.Request) {

}
