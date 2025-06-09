package swojobenginereceiver

import (
	"context"
	"errors"
	"fmt"
	"net"

	jobEngineEvents "github.com/solarwinds/solarwinds-otel-collector/pkg/job-engine-events"
	model "github.com/solarwinds/solarwinds-otel-collector/receiver/swojobenginereceiver/internal/model"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componentstatus"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type jobengineGRPCReceiver struct {
	conf         *Config
	nextConsumer consumer.Metrics
	settings     receiver.Settings
	serverGRPC   *grpc.Server
	server       *server

	obsrepGRPC *receiverhelper.ObsReport
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	nextConsumer consumer.Metrics
	obsrep       *receiverhelper.ObsReport
	logger       *zap.Logger

	jobEngineEvents.UnimplementedJobEngineEventsServer
}

// OnJobFinished implements jobEngineEvents.OnJobFinished
func (s *server) NotifyJobFinished(_ context.Context, in *jobEngineEvents.NotifyJobFinishedRequest) (*jobEngineEvents.NotifyJobFinishedResponse, error) {
	metrics, err := model.Transform_toMetrics(in, s.logger)

	if err != nil {
		message := "Failed to transform metrics."
		s.logger.Error(message, zap.Error(err))

		return &jobEngineEvents.NotifyJobFinishedResponse{}, nil
	}

	ctx := s.obsrep.StartMetricsOp(context.Background())
	err = s.nextConsumer.ConsumeMetrics(ctx, *metrics)
	s.obsrep.EndMetricsOp(ctx, "protobuf", metrics.DataPointCount(), err)

	return &jobEngineEvents.NotifyJobFinishedResponse{}, nil
}

func (r *jobengineGRPCReceiver) Start(ctx context.Context, host component.Host) error {
	var err error

	r.serverGRPC = grpc.NewServer()
	if err != nil {
		return fmt.Errorf("failed create grpc server error: %w", err)
	}

	r.server = &server{
		nextConsumer: r.nextConsumer,
		obsrep:       r.obsrepGRPC,
		logger:       r.settings.Logger,
	}

	jobEngineEvents.RegisterJobEngineEventsServer(r.serverGRPC, r.server)

	err = r.startGRPCServer(ctx, host)
	if err != nil {
		return fmt.Errorf("failed to start grpc server error: %w", err)
	}

	return err
}

func (r *jobengineGRPCReceiver) startGRPCServer(ctx context.Context, host component.Host) error {
	r.settings.Logger.Info("Starting GRPC server", zap.Int("endpoint.port", r.conf.EndpointPort))
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", r.conf.EndpointPort))
	if err != nil {
		return err
	}

	go func() {
		if errGRPC := r.serverGRPC.Serve(listener); !errors.Is(errGRPC, grpc.ErrServerStopped) && errGRPC != nil {
			componentstatus.ReportStatus(host, componentstatus.NewFatalErrorEvent(errGRPC))
		}
	}()
	return nil
}

func (r *jobengineGRPCReceiver) Shutdown(_ context.Context) error {
	if r.serverGRPC != nil {
		r.serverGRPC.GracefulStop()
	}

	return nil
}
