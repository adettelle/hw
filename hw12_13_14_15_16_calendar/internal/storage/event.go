package storage

import "time"

type Event struct {
	ID           string
	Title        string
	CreatedAt    time.Time
	Start        time.Time // Дата и время события;
	End          time.Time // дата и время окончания (Длительность события);
	Description  string    // Описание события - длинный текст, опционально;
	UserID       string    // ID пользователя, владельца события;
	Notification time.Time
	// (дата и время, когда высылать уведомление) За сколько времени высылать уведомление
	Notified bool
}

type EventCreateDTO struct {
	Title        string    `json:"title" validate:"required,min=1"`
	Start        time.Time `json:"dateStart" validate:"required"` // Дата и время события;
	End          time.Time `json:"dateEnd" validate:"required"`   // дата и время окончания (Длительность события);
	Description  string    `json:"description"`                   // Описание события - длинный текст, опционально;
	Notification time.Time `json:"notification"`                  // За сколько времени высылать уведомление, опционально.
	Notified     bool      `json:"notified"`
}

type EventUpdateDTO struct {
	Title        *string    `json:"title" validate:"required,min=1"`
	Start        *time.Time `json:"dateStart" validate:"required"` // Дата и время события;
	End          *time.Time `json:"dateEnd" validate:"required"`   // Длительность события (или дата и время окончания);
	Description  *string    `json:"description"`                   // Описание события - длинный текст, опционально;
	Notification *time.Time `json:"notification"`                  // За сколько времени высылать уведомление, опционально.
	Notified     bool       `json:"notified"`
}

type EventGetDTO struct {
	ID           string    `json:"id" validate:"required,min=1"`
	Title        string    `json:"title" validate:"required,min=1"`
	Start        time.Time `json:"dateStart" validate:"required"` // Дата и время события;
	End          time.Time `json:"dateEnd" validate:"required"`   // дата и время окончания (Длительность события);
	Description  string    `json:"description"`                   // Описание события - длинный текст, опционально;
	Notification time.Time `json:"notification"`                  // За сколько времени высылать уведомление, опционально.
	Notified     bool      `json:"notified"`
}

type EventToNotify struct {
	ID     string    `json:"id"`
	Title  string    `json:"title"`
	Start  time.Time `json:"dateStart"` // Дата и время события;
	UserID string    `json:"userId" `
}
