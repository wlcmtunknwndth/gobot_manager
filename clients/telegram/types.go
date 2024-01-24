package telegram

// UpdatesResponse — returns a list of updates.
type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// Update — update json from tg.api.
type Update struct {
	ID      int               `json:"update_id"`
	Message *IncomningMessage `json:"message"`
}

// IncomingMessage — json with messages from user.
type IncomningMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

// From — response from api with Username.
type From struct {
	Username string `json:"username"`
}

// Chat — response from api with ChatID
type Chat struct {
	ID int `json:"id"`
}
