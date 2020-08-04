/*
	В первый раз стартует фуннкция и отравлчяет полное значение
	Слудующее отправляет дельту
*/

package kiddylineprocessor

import (
	"context"
	"math"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vlasovoleg/kiddyLineProcessor/internal/store"
	"github.com/vlasovoleg/kiddyLineProcessor/pkg/api"
	"google.golang.org/grpc"
)

// GRPCServer ...
type gRPCServer struct {
	store *store.Store
	loger *logrus.Logger
}

func (kp *Kiddylineprocessor) runGRPC(ctx context.Context, wg *sync.WaitGroup, store *store.Store, loger *logrus.Logger) {
	kp.loger.Trace("Kiddylineprocessor : runGRPC")
	lis, err := net.Listen("tcp", kp.config.GRPC.Address)
	if err != nil {
		kp.loger.Panic("Kiddylineprocessor : gRPC : failed to listen : ", err)
	}
	grpcServer := grpc.NewServer()
	g := &gRPCServer{
		store: store,
		loger: loger,
	}
	api.RegisterProcessorServer(grpcServer, g)

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			kp.loger.Error("Kiddylineprocessor : runGRPC : grpcServer.Serve : ", err)
		}
	}()

	<-ctx.Done()
	grpcServer.Stop()
	kp.loger.Info("Kiddylineprocessor : gRPC : Server closed")
	wg.Done()
}

// SubscribeOnSportLines grpc method
func (s *gRPCServer) SubscribeOnSportLines(stream api.Processor_SubscribeOnSportLinesServer) error {
	ctx, stop := context.WithCancel(context.Background())
	for {
		data, err := stream.Recv()
		if err != nil {
			stop()
			s.loger.Error("gRPCServer : SubscribeOnSportLines : stream.Recv : ", err)
			return err
		}
		stop()
		ctx, stop = context.WithCancel(context.Background())
		s.loger.Debug("SubscribeOnSportLines : gRPC SubscribeOnSportLines : sports:", data.Sports, "time:", data.Time)
		go s.subscribeOnSportLinesSender(ctx, data.Sports, stream, data.Time)
	}
}

func (s *gRPCServer) subscribeOnSportLinesSender(ctx context.Context, sports []string, stream api.Processor_SubscribeOnSportLinesServer, t int32) {
	sportsFunc := []func(*api.Coefficients){}
	firstCoefficients := make(map[string]float32)

	for _, sport := range sports {
		if sport == "baseball" {
			sportsFunc = append(sportsFunc, s.getBaseball)
			s.getBaseball(&api.Coefficients{Coefficients: firstCoefficients})
		}
		if sport == "football" {
			sportsFunc = append(sportsFunc, s.getFootball)
			s.getFootball(&api.Coefficients{Coefficients: firstCoefficients})
		}
		if sport == "soccer" {
			sportsFunc = append(sportsFunc, s.getSoccer)
			s.getSoccer(&api.Coefficients{Coefficients: firstCoefficients})
		}
	}

	err := stream.Send(&api.Coefficients{Coefficients: firstCoefficients})
	if err != nil {
		s.loger.Error("GRPCServer : subscribeOnSportLinesSender : stream.Send : ", err)
	}

	for {
		time.Sleep(time.Second * time.Duration(t))
		select {
		case <-ctx.Done():
			s.loger.Trace("exit from loop")
			return
		default:
			coefficients := api.Coefficients{}
			coefficients.Coefficients = make(map[string]float32)
			for i := 0; i < len(sportsFunc); i++ {
				(sportsFunc)[i](&coefficients)
			}
			for s, c := range firstCoefficients {
				coefficients.Coefficients[s] -= c
				coefficients.Coefficients[s] = float32(math.Ceil(float64(coefficients.Coefficients[s])*1000) / 1000)
			}
			err := stream.Send(&coefficients)
			if err != nil {
				s.loger.Error("GRPCServer : subscribeOnSportLinesSender : stream.Send : ", err)
			}
		}
	}
}

func (s *gRPCServer) getBaseball(c *api.Coefficients) {
	dbc, err := (*s.store).BaseballRepository().GetCoefficient()
	if err != nil {
		s.loger.Error("kiddylineprocessor : gRPC API : getBaseball : BaseballRepository.GetCoefficient : " + err.Error())
	}
	c.Coefficients["baseball"] = dbc
}
func (s *gRPCServer) getFootball(c *api.Coefficients) {
	dbc, err := (*s.store).FootballRepository().GetCoefficient()
	if err != nil {
		s.loger.Error("kiddylineprocessor : gRPC API : getFootball : FootballRepository.GetCoefficient : " + err.Error())
	}
	c.Coefficients["football"] = dbc
}
func (s *gRPCServer) getSoccer(c *api.Coefficients) {
	dbc, err := (*s.store).SoccerRepository().GetCoefficient()
	if err != nil {
		s.loger.Error("kiddylineprocessor : gRPC API : getSoccer : SoccerRepository.GetCoefficient : " + err.Error())
	}
	c.Coefficients["soccer"] = dbc
}
