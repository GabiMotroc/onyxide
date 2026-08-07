package proj

type Project struct {
	AppType  string `json:"apptype"`
	Location string `json:"location"`
	Score    int    `json:"score"`
}
