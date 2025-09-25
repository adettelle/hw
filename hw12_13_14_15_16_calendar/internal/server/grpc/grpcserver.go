package internalgrpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	pb "github.com/adettelle/hw/hw12_13_14_15_calendar/api"
	"github.com/adettelle/hw/hw12_13_14_15_calendar/configs"
	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCServer struct {
	pb.UnimplementedStoragerServer
	grpcServer *grpc.Server
	cfg        *configs.Config
	logg       *zap.Logger
	Storager   app.Storager
}

func NewGRPCServer(cfg *configs.Config, logg *zap.Logger, storager app.Storager) *GRPCServer {
	server := grpc.NewServer()
	return &GRPCServer{cfg: cfg, logg: logg, Storager: storager, grpcServer: server}
}

func (s *GRPCServer) AddEventByID(ctx context.Context, in *pb.AddEventByIDRequest) (*pb.AddEventByIDResponse, error) {
	var response pb.AddEventByIDResponse
	event := storage.EventCreateDTO{
		Title:        in.EventCreateDTO.Title,
		DateStart:    in.EventCreateDTO.StartTime.AsTime(),
		DateEnd:      in.EventCreateDTO.EndTime.AsTime(),
		Description:  in.EventCreateDTO.Description,
		Notification: in.EventCreateDTO.Notification.AsTime(),
	}

	res, err := s.Storager.AddEventByID(ctx, event, in.UserID)
	if err != nil {
		response.Error = err.Error()
		return &response, err
	}

	response.Id = res

	return &response, nil
}

func (s *GRPCServer) UpdateEventByID(ctx context.Context,
	in *pb.UpdateEventByIDRequest,
) (*pb.UpdateEventByIDResponse, error) {
	var response pb.UpdateEventByIDResponse
	date := in.EventCreateDTO.StartTime.AsTime()
	duration := in.EventCreateDTO.EndTime.AsTime()
	notification := in.EventCreateDTO.Notification.AsTime()

	event := storage.EventUpdateDTO{
		Title:        &in.EventCreateDTO.Title,
		Date:         &date,
		Duration:     &duration,
		Description:  &in.EventCreateDTO.Description,
		Notification: &notification,
	}

	err := s.Storager.UpdateEventByID(ctx, in.Id, event, in.UserID)
	if err != nil {
		response.Error = err.Error()
		return &response, err
	}

	return &response, nil
}

func (s *GRPCServer) DeleteEventByID(ctx context.Context,
	in *pb.DeleteEventByIDRequest,
) (*pb.DeleteEventByIDResponse, error) {
	var response pb.DeleteEventByIDResponse

	err := s.Storager.DeleteEventByID(ctx, in.Id)
	if err != nil {
		response.Error = err.Error()
		return &response, err
	}

	return &response, nil
}

func (s *GRPCServer) GetEventListingByUserID(_ context.Context,
	in *pb.GetEventListingByUserIDRequest,
) (*pb.GetEventListingByUserIDResponse, error) {
	var response pb.GetEventListingByUserIDResponse

	events, err := s.Storager.GetEventListingByUserID(in.UserID, in.Date.AsTime(), in.Period.String())
	if err != nil {
		response.Error = err.Error()
		return &response, err
	}

	resEvents := make([]*pb.Event, 0, len(events))
	for _, e := range events {
		resEvents = append(resEvents, &pb.Event{
			Id:           e.ID,
			Title:        e.Title,
			Description:  e.Description,
			CreatedAt:    timestamppb.New(e.CreatedAt),
			StartTime:    timestamppb.New(e.Date),
			EndTime:      timestamppb.New(e.Duration),
			Notification: timestamppb.New(e.Notification),
		})
	}

	response.Event = resEvents

	return &response, nil
}

func (s *GRPCServer) GetEventByID(_ context.Context, in *pb.GetEventByIDRequest) (*pb.GetEventByIDResponse, error) {
	var response pb.GetEventByIDResponse

	e, err := s.Storager.GetEventByID(in.Id, in.UserID)
	if err != nil {
		response.Error = err.Error()
		return &response, err
	}
	response.Event = &pb.Event{
		Id:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		StartTime:   timestamppb.New(e.Date),
		EndTime:     timestamppb.New(e.Duration),
		CreatedAt:   timestamppb.New(e.CreatedAt),
	}

	return &response, nil
}

func (s *GRPCServer) Notify(_ context.Context, _ *pb.NotifyRequest) (*pb.NotifyResponse, error) {
	var response *pb.NotifyResponse

	return response, nil
}

func (s *GRPCServer) Start(ctx context.Context, logg *zap.Logger) error { // port string storager app.Storager,
	// определяем порт для сервера
	_, port, err := net.SplitHostPort(s.cfg.GRPCAddress)
	if err != nil {
		log.Fatal(err)
	}
	_, err = strconv.Atoi(port)
	if err != nil {
		log.Fatal(fmt.Errorf("invalid port: '%s'", port))
	}

	lc := net.ListenConfig{}
	grpcListen, err := lc.Listen(ctx, "tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}

	// регистрируем сервис
	pb.RegisterStoragerServer(s.grpcServer, &GRPCServer{Storager: s.Storager})
	logg.Info("start grpc server success ", zap.Any("endpoint", grpcListen.Addr()))

	// получаем запрос gRPC
	if err := s.grpcServer.Serve(grpcListen); err != nil {
		logg.Info("failed to grpc server serve")
		return err
	}

	return nil
}

func (s *GRPCServer) Close() error {
	s.grpcServer.GracefulStop()
	return nil
}
