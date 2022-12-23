package dto

type TopTopics struct {
	TopicID          uint   `json:"topicId"`
	TopicName        string `json:"topicName"`
	TopicDescription string `json:"topicDescription"`
	PostCount        int    `json:"postCount"`
}
