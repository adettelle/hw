package storage

import "time"

type Event struct {
	ID           string
	Title        string
	CreatedAt    time.Time
	Date         time.Time // Дата и время события;
	Duration     time.Time // дата и время окончания (Длительность события);
	Description  string    // Описание события - длинный текст, опционально;
	UserID       string    // ID пользователя, владельца события;
	Notification time.Time
	// (дата и время, когда высылать уведомление) За сколько времени высылать уведомление, опционально.
}

type EventCreateDTO struct { // EventCreateRequestDTO
	Title       string    `json:"title" validate:"required,min=1"`
	DateStart   time.Time `json:"dateStart" validate:"required"` // Дата и время события;
	DateEnd     time.Time `json:"dateEnd" validate:"required"`   // дата и время окончания (Длительность события);
	Description string    `json:"description"`                   // Описание события - длинный текст, опционально;
	// UserID       string    `json:"userID" validate:"required,min=1"`// ID пользователя, владельца события;
	Notification time.Time `json:"notification"` // За сколько времени высылать уведомление, опционально.
}

type EventUpdateDTO struct { // EventUpdateParams
	Title        *string    `json:"title" validate:"required,min=1"`
	Date         *time.Time `json:"dateStart" validate:"required"` // Дата и время события;
	Duration     *time.Time `json:"dateEnd" validate:"required"`   // Длительность события (или дата и время окончания);
	Description  *string    `json:"description"`                   // Описание события - длинный текст, опционально;
	Notification *time.Time `json:"notification"`                  // За сколько времени высылать уведомление, опционально.
}

type EventGetDTO struct { // EventCreateRequestDTO
	ID          string    `json:"id" validate:"required,min=1"`
	Title       string    `json:"title" validate:"required,min=1"`
	DateStart   time.Time `json:"dateStart" validate:"required"` // Дата и время события;
	DateEnd     time.Time `json:"dateEnd" validate:"required"`   // дата и время окончания (Длительность события);
	Description string    `json:"description"`                   // Описание события - длинный текст, опционально;
	// UserID       string    `json:"userID" validate:"required,min=1"`// ID пользователя, владельца события;
	Notification time.Time `json:"notification"` // За сколько времени высылать уведомление, опционально.
}
