package dto

import "time"

type EventName string

type KafkaMessageEvent struct {
	Name EventName `json:"name"`
}

type KafkaMessageMeta struct {
	Sender    string     `json:"sender"`
	SendingAt time.Time  `json:"sendingAt"`
	ExpiredAt *time.Time `json:"expiredAt"`
	Version   *string    `json:"version"`
}

type DataType string

type KafkaMessageBody[T any] struct {
	Type DataType `json:"type"`
	Data T        `json:"data"`
}

type KafkaMessage[T any] struct {
	Event KafkaMessageEvent   `json:"event"`
	Meta  KafkaMessageMeta    `json:"meta"`
	Body  KafkaMessageBody[T] `json:"body"`
}
