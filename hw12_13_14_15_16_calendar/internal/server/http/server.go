package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/adettelle/hw/hw12_13_14_15_calendar/configs"
	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Server struct { // TODO
	cfg      *configs.Config
	logg     *zap.Logger
	storager app.Storager
}

type Logger interface { // TODO
}

type Application interface { // TODO
}

func NewServer(cfg *configs.Config, logg *zap.Logger, _ Application, storager app.Storager) *Server {
	return &Server{cfg: cfg, logg: logg, storager: storager}
}

func (s *Server) Start(ctx context.Context, logg *zap.Logger) error {
	// mux := http.NewServeMux()

	// mux := Router()
	eventHandlers := New(s.storager, logg)
	router := NewRouter(eventHandlers, logg)

	srv := &http.Server{
		Addr:         s.cfg.Address,
		Handler:      router, // mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	// mux.HandleFunc(`/`, mainPage)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	// TODO
	s.logg.Info("Gracefully shutting down server")
	return nil
}

func (eh *EventHandlers) mainPage(res http.ResponseWriter, _ *http.Request) {
	res.Write([]byte("Hello!"))
}

type EventHandlers struct {
	Storager app.Storager
	//	DBCon    database.DBConnector
	Logg *zap.Logger
}

func New(storager app.Storager, logg *zap.Logger) *EventHandlers {
	return &EventHandlers{
		Storager: storager,
		Logg:     logg,
	}
}

func (eh *EventHandlers) GetEventByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.PathValue("userid")
	eventID := r.PathValue("id")

	event, err := eh.Storager.GetEventByID(eventID, userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := json.Marshal(event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (eh *EventHandlers) AddEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.PathValue("userid")

	var buf bytes.Buffer
	var event storage.EventCreateDTO

	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		eh.Logg.Error("error in reading body:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &event); err != nil {
		eh.Logg.Error("error in unmarshalling json:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(event)
	if err != nil {
		eh.Logg.Error("error in validating:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = eh.Storager.AddEventByID(context.Background(), event, userID)
	if err != nil {
		eh.Logg.Error("error in adding event:", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (eh *EventHandlers) DeleteEventByID(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("id")
	err := eh.Storager.DeleteEventByID(context.Background(), eventID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (eh *EventHandlers) UpdateEventeByID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userid")
	eventID := r.PathValue("id")

	var buf bytes.Buffer
	var event storage.EventUpdateDTO

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		eh.Logg.Error("error in reading body:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &event); err != nil {
		eh.Logg.Error("error in unmarshalling json:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = eh.Storager.UpdateEventByID(context.Background(), eventID, event, userID)
	if err != nil {
		eh.Logg.Error("error in updating event:", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func NewEventsListResponse(events []storage.EventGetDTO) []*storage.EventGetDTO {
	res := []*storage.EventGetDTO{}
	for _, event := range events {
		res = append(res, &event)
	}
	return res
}

func (eh *EventHandlers) GetEventListingByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userid")
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "day"
	}

	var parsedTime time.Time
	var err error

	date := r.URL.Query().Get("date")

	if date == "" {
		now := time.Now()
		parsedTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 1, 0, time.Local)
	} else {
		parsedTime, err = time.Parse("2006-01-02", date)
		if err != nil {
			eh.Logg.Error("error in parsing time:", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	events, err := eh.Storager.GetEventListingByUserID(userID, parsedTime, period)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if len(events) == 0 {
		eh.Logg.Info("no events")
		w.WriteHeader(http.StatusNoContent) // нет данных для ответа
		return
	}

	resp, err := json.Marshal(events)
	if err != nil {
		eh.Logg.Error("error in marshalling json", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		eh.Logg.Error("error in writing resp:", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
