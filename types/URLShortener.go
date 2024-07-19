package types

type Request struct {
	Url string `json:"url"`
}

type Response struct {
	ShortenUrl string `json:"shorten_url"`
	ExpireTime string `json:"expire_time"`
}
