package twitch

type ChallengeWebhook struct {
	Challenge    string        `json:"challenge"`
	Subscription *Subscription `json:"subscription"`
}

type NotificationWebhook struct {
	Subscription *Subscription      `json:"subscription"`
	Event        *map[string]string `json:"event"`
}

type RevocationWebhook struct {
	Subscription *Subscription `json:"subscription"`
}

type Subscription struct {
	Id        string            `json:"id"`
	Status    string            `json:"status"`
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Cost      int64             `json:"cost"`
	Condition map[string]string `json:"condition"`
	Transport *Transport        `json:"transport"`
	CreatedAt string            `json:"created_at"`
}

type Transport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
}

type Event struct {
	UserId               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	BroadcasterUserId    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
}
