package qbitorrent

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type Config struct {
	Host        string `env:"BPM_TORRENT_QBITORRENT_HOST,required"`
	Port        int    `env:"BPM_TORRENT_QBITORRENT_PORT,required"`
	Credentials struct {
		Username string `env:"BPM_TORRENT_QBITORRENT_USER,required"`
		Password string `env:"BPM_TORRENT_QBITORRENT_PASS,required"`
	}
}

type QBitorrent struct {
	client *http.Client
	config *Config
}

func New(c *Config) *QBitorrent {
	return &QBitorrent{client: &http.Client{}, config: c}
}

func (q *QBitorrent) host() string {
	return fmt.Sprintf("%v:%v", q.config.Host, q.config.Port)
}

func (q *QBitorrent) endpoint(p string) string {
	return fmt.Sprintf("%v%v", q.host(), p)
}

func (q *QBitorrent) Authenticate() ([]*http.Cookie, error) {
	loginData := url.Values{"username": {q.config.Credentials.Username}, "password": {q.config.Credentials.Password}}
	res, err := q.client.PostForm(q.endpoint("/api/v2/auth/login"), loginData)
	if err != nil {
		return nil, fmt.Errorf("login failed: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		data, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("login failed: %s: %s", res.Status, data)
	}

	return res.Cookies(), nil
}

func (q *QBitorrent) authenticateRequest(req *http.Request) error {
	cookies, err := q.Authenticate()
	if err != nil {
		return err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	return nil
}

func (q *QBitorrent) AddTorrent(fURL string) error {
	bs, err := downloadFile(fURL)
	if err != nil {
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("torrents", "upload.torrent")
	if _, err := part.Write(bs); err != nil {
		return err
	}
	defer writer.Close()

	req, _ := http.NewRequest("POST", q.endpoint("/api/v2/torrents/add"), body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	if err := q.authenticateRequest(req); err != nil {
		return err
	}

	res, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		data, _ := io.ReadAll(res.Body)
		return fmt.Errorf("qBittorrent returned %s: %s", res.Status, data)
	}

	return nil
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
