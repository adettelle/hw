package internalhttp

import (
	"github.com/adettelle/hw/hw12_13_14_15_calendar/pkg/mware/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(h *EventHandlers, logg *zap.Logger) chi.Router {
	r := chi.NewRouter()

	r.Get(`/`, logger.WithLogging(h.mainPage, logg))
	r.Get(`/user/{userid}/event/{id}`, logger.WithLogging(h.GetEventByID, logg))
	r.Put(`/user/{userid}/event/`, logger.WithLogging(h.AddEvent, logg))
	r.Post(`/update/user/{userid}/event/{id}`, logger.WithLogging(h.UpdateEventeByID, logg)) // TODO !!!
	r.Delete(`/user/{userid}/event/{id}`, logger.WithLogging(h.DeleteEventByID, logg))
	r.Get(`/user/{userid}/events/`, logger.WithLogging(h.GetEventListingByUserID, logg))

	return r
}
