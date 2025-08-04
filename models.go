package main

type ThreadsResponse struct {
	Threads []Thread `json:"threads"`
}

type ThreadResponse struct {
	Thread Thread `json:"thread"`
}

type Thread struct {
	ThreadId        int64  `json:"thread_id,omitempty"`
	ThreadTitle     string `json:"thread_title,omitempty"`
	CreatorUserID   int64  `json:"creator_user_id,omitempty"`
	CreatorUsername string `json:"creator_username,omitempty"`
	ThreadIsClosed  bool   `json:"thread_is_closed,omitempty"`
	LastReplyDate   int64
	Captcha         string
	Restrictions    Restrictions   `json:"restrictions,omitempty"`
	ThreadPrefixes  []ThreadPrefix `json:"thread_prefixes,omitempty"`
}

type Restrictions struct {
	ReplyDelay int64 `json:"reply_delay,omitempty"`
}

type ThreadPrefix struct {
	PrefixID    int64  `json:"prefix_id,omitempty"`
	PrefixTitle string `json:"prefix_title,omitempty"`
}

type WhitelistEntry struct {
	ThreadID int64  `json:"thread_id"`
	Captcha  string `json:"captcha,omitempty"`
}
