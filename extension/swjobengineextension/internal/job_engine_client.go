package internal

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/solarwinds/solarwinds-otel-collector/extension/swjobengineextension/internal/job_engine_service"
	"github.com/solarwinds/solarwinds-otel-collector/extension/swjobengineextension/internal/models"
)

func pointerTo[T ~string](s T) *T {
	return &s
}

func CreateGrpcConnection(config *Config, logger *zap.Logger) (conn *grpc.ClientConn, err error) {
	cert, err := tls.LoadX509KeyPair(config.TLS_PublicKey, config.TLS_PrivateKey)
	if err != nil {
		log.Fatalf("failed to load client cert: %v", err)
	}

	logger.Info("Connecting..")

	ca := x509.NewCertPool()
	caFilePath := config.TLS_PublicKey
	caBytes, err := os.ReadFile(caFilePath)
	if err != nil {
		log.Fatalf("failed to read ca cert %q: %v", caFilePath, err)
	}
	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		log.Fatalf("failed to parse %q", caFilePath)
	}

	tlsConfig := &tls.Config{
		ServerName:   config.TLS_ServerNamme,
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
	}

	// Set up a connection to the server.
	conn, err = grpc.NewClient(config.JobEngineServiceEndpoint, grpc.WithTimeout(5*time.Second), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		log.Fatalf("failed to create connection %q", caFilePath)
	}
	return conn, err
}

func CreateScheduledJobs(jobNamespace string, jobType string, credentialsXml string, jobDescription string) *pb.ScheduledJob {
	runOnce := false
	use64Bits := true
	isCustomDebugLogEnabled := true

	return &pb.ScheduledJob{
		Frequency:           durationpb.New(time.Second * 120),
		IsOneShot:           false,
		InitialWait:         durationpb.New(time.Second * 1),
		RunOnce:             &runOnce,
		NotificationAddress: "grpc://localhost:18733",
		State:               "",
		CronExpression:      pointerTo("* * * * *"),
		TimeZone: &pb.TimeZoneInfoType{
			Value: "UTC;0;(UTC) Coordinated Universal Time;Coordinated Universal Time;Coordinated Universal Time;;",
		},
		Start: timestamppb.New(time.Now()),
		End:   timestamppb.New(time.Now().Add(time.Hour * 240)),
		Job: &pb.JobDescription{
			IsCustomDebugLogEnabled: &isCustomDebugLogEnabled,
			Use_64Bit:               &use64Bits,
			JobNamespace:            jobNamespace,
			JobType:                 jobType,
			EndpointAddress:         "",
			LegacyEngine:            pointerTo(""),
			JobConfiguration:        jobDescription,
			Timeout:                 durationpb.New(time.Second * 10),
			HungTimeout:             durationpb.New(time.Second * 10),
			ResultTtl:               durationpb.New(time.Second * 10),
			Credential: &pb.Credential{
				//CredentialType: "http://www.solarwinds.com/jobengine/2008/03/jobCredential#empty",
				CredentialType: "http://www.solarwinds.com/orion/2008/03/credentials#pollersjob",
				Data: &pb.XmlElementType{
					Value: credentialsXml,
				},
			},
			SupportedRoles: pb.PackageType_PACKAGE_TYPE_ALL_POLLERS.Enum(),
			Priority:       pb.JobPriorityType_JOB_PRIORITY_NORMAL.Enum(),
			ThrottleGroup:  pointerTo(""),
			WorkerType:     pb.WorkerType_WORKER_TYPE_NATIVE.Enum(),
			HostAddress: &pb.HostAddress{
				Address:     "127.0.0.1",
				AddressType: *pb.AddressType_ADDRESS_TYPE_IPV4.Enum(),
			},
			CustomWorkerCommandLine: pointerTo(""),
			CustomWorkerCommandArgs: pointerTo(""),
		},
	}
}

type JobEngineClient struct {
	conn   *grpc.ClientConn
	client pb.JobEngineClient
	ctx    context.Context
	cancel context.CancelFunc
	logger *zap.Logger
}

func NewGrpcClient(config *Config, logger *zap.Logger) (*JobEngineClient, error) {
	conn, err := CreateGrpcConnection(config, logger)

	if err != nil {
		logger.Error("did not connect", zap.Error(err))
	}

	client := pb.NewJobEngineClient(conn)

	logger.Info("Connected to JobEngine", zap.String("endpoint", config.JobEngineServiceEndpoint))

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return &JobEngineClient{
		conn:   conn,
		client: client,
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
	}, nil
}

func (c *JobEngineClient) Cancel() {
	c.conn.Close()
	c.cancel()
}

func (c *JobEngineClient) DeleteJobs() error {
	// clear jobs
	_, err := c.client.Clear(c.ctx, &pb.ClearRequest{
		ProductNamespaces: []string{"SolarWinds.PCU.Pollers", "SolarWinds.Orion.Core.Pollers"},
	})

	if err != nil {
		c.logger.Error("could not call Clear", zap.Error(err))
		return nil
	}

	c.logger.Debug("Jobs cleared")

	return nil
}

func (c *JobEngineClient) ListJobs() (int, error) {
	// enumerate scheduled jobs
	r, err := c.client.EnumerateScheduledJobs(c.ctx, &emptypb.Empty{})

	if err != nil {
		c.logger.Error("could not call EnumerateScheduledJobs", zap.Error(err))
		return 0, err
	}

	for _, job := range r.ScheduledJobs {
		c.logger.Info("JobEngine job",
			zap.String("id", job.Id.Value),
			zap.String("state", job.ScheduledJob.State))
	}

	return len(r.ScheduledJobs), nil
}

// add job
func (c *JobEngineClient) CreateJob_PCU(variables map[string]string) (string, error) {
	job := CreateScheduledJobs(
		models.JOB_NAMESPACE_PCU,
		models.JOB_TYPE_PCU,
		models.GetSnmpV2Credentials(variables),
		models.GetPCUJobDescription(variables))

	uid, err := c.client.AddJob(c.ctx, job)

	if err != nil {
		c.logger.Error("could not call CreateJob_PCU", zap.Error(err))
		return "", err
	}

	return uid.String(), nil
}

func (c *JobEngineClient) CreateJob_SNMP(variables map[string]string) (string, error) {
	// add job
	job := CreateScheduledJobs(
		models.JOB_NAMESPACE_CPU,
		models.JOB_TYPE_CPU,
		models.GetSnmpV2Credentials(variables),
		models.GetCoreSnmpJobDescription(variables))

	uid, err := c.client.AddJob(c.ctx, job)

	if err != nil {
		c.logger.Error("could not call CreateJob_SNMP", zap.Error(err))
		return "", err
	}

	return uid.String(), nil
}

func (c *JobEngineClient) CreateJob_ICMP(variables map[string]string) (string, error) {
	// add job
	job := CreateScheduledJobs(
		models.JOB_NAMESPACE_CPU,
		models.JOB_TYPE_CPU,
		models.GetSnmpV2Credentials(variables),
		models.GetCoreIcmpJobDescription(variables))

	uid, err := c.client.AddJob(c.ctx, job)

	if err != nil {
		c.logger.Error("could not call CreateJob_ICMP", zap.Error(err))
		return "", err
	}

	return uid.String(), nil
}

func (c *JobEngineClient) CreateJob_CoreInventory(variables map[string]string) (string, error) {
	// add job
	job := CreateScheduledJobs(
		models.JOB_NAMESPACE_CPU,
		models.JOB_TYPE_CPU,
		models.GetSnmpV2Credentials(variables),
		models.GetCoreInventoryJobDescription(variables))

	uid, err := c.client.AddJob(c.ctx, job)

	if err != nil {
		c.logger.Error("could not call CreateJob_CoreInventory", zap.Error(err))
		return "", err
	}

	return uid.String(), nil
}
