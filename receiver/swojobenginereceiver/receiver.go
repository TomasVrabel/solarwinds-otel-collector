package swojobenginereceiver

import (
	"context"
	"errors"
	"fmt"
	"net"

	jobEngineEvents "github.com/solarwinds/solarwinds-otel-collector/receiver/swojobenginereceiver/internal/job-engine-events"
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

	obsrepGRPC *receiverhelper.ObsReport
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	nextConsumer consumer.Metrics
	obsrep       *receiverhelper.ObsReport

	jobEngineEvents.UnimplementedJobEngineEventsServer
}

// OnJobFinished implements jobEngineEvents.OnJobFinished
func (s *server) NotifyJobFinished(_ context.Context, in *jobEngineEvents.NotifyJobFinishedRequest) (*jobEngineEvents.NotifyJobFinishedResponse, error) {

	ctx := s.obsrep.StartMetricsOp(context.Background())
	metrics, err := model.Transform_PCU_toMetrics(in)
	metricsRecordCount := 1
	err = s.nextConsumer.ConsumeMetrics(ctx, *metrics)
	s.obsrep.EndMetricsOp(ctx, "protobuf", metricsRecordCount, err)

	return &jobEngineEvents.NotifyJobFinishedResponse{}, nil
}

func (r *jobengineGRPCReceiver) Start(ctx context.Context, host component.Host) error {
	var err error

	r.serverGRPC = grpc.NewServer()
	if err != nil {
		return fmt.Errorf("failed create grpc server error: %w", err)
	}

	jobEngineEvents.RegisterJobEngineEventsServer(r.serverGRPC, &server{
		nextConsumer: r.nextConsumer,
		obsrep:       r.obsrepGRPC,
	})

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
