package entity

type Post struct {
	Id    string    `json:"id" bson:"_id,omitempty"`
	Title string `json:"title"`
	Text  string `json:"text"`
}
