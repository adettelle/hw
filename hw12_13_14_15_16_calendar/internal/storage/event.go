package storage

import "time"

type Event struct {
	ID           string
	Title        string
	Date         time.Time     // Дата и время события;
	Duration     time.Duration // Длительность события (или дата и время окончания);
	Description  string        // Описание события - длинный текст, опционально;
	UserID       string        // ID пользователя, владельца события;
	Notification time.Time     // За сколько времени высылать уведомление, опционально.
}
