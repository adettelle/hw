package storage

import "time"

type Event struct {
	ID           string
	Title        string
	Date         time.Time // Дата и время события;
	Duration     time.Time // дата и время окончания (Длительность события);
	Description  string    // Описание события - длинный текст, опционально;
	UserID       string    // ID пользователя, владельца события;
	Notification time.Time // (дата и время, когда высылать уведомление) За сколько времени высылать уведомление, опционально.
}
