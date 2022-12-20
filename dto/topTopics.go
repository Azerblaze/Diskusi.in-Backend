package dto

type TopTopics struct {
	TopicID   uint   `json:"topic_id"`
	TopicName string `json:"topic_name"`
	PostCount int    `json:"post_count"`
}
