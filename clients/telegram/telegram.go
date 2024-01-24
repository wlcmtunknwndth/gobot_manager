package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/wlcmtunknwndth/gobot_manager/lib/error_handler"
)

// Client — client to work with api.
type Client struct {
	host     string      // api host-server
	basePath string      // path to storage
	client   http.Client // local client to work with host
} /////

const (
	getUpdatesMethod  = "getUpdates"
	SendMessageMethod = "sendMessage"
)

// tg-bot.com/bot<token>

func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

// newBasePath — "bot" + token.
func newBasePath(token string) string {
	return "bot" + token // https://api.telegram.org/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/
}

// Updates — get updates from host.
func (c *Client) Updates(offset, limit int) ([]Update, error) {
	q := url.Values{} // map[string] []string for request

	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q) // do req by a chosen command(sendMessage or getUpdates)
	if err != nil {
		return nil, err
	}
	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil { //parses json
		return nil, err
	}

	return res.Result, nil
}

// SendMessage —— sends messages to proper chat
func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest(SendMessageMethod, q)

	if err != nil {
		return error_handler.Wrap("can't send message", err)
	}

	return nil
}

// doRequest —— request bt a chosend method(sendMessage or getUpdates)
func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = error_handler.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
