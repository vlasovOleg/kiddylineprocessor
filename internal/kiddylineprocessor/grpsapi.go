/*
	В первый раз стартует фуннкция и отравлчяет полное значение
	Слудующее отправляет дельту
*/

package kiddylineprocessor

import (
	"math"
	"net"
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

// NewGRPC ..
func (kp *Kiddylineprocessor) NewGRPC(store *store.Store, loger *logrus.Logger) {
	kp.loger.Trace("Kiddylineprocessor : NewGRPC...")
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

	err = grpcServer.Serve(lis)
	if err != nil {
		kp.loger.Error("Kiddylineprocessor : NewGRPC : grpcServer.Serve : ", err)
	}
}

// SubscribeOnSportLines grpc method
func (s *gRPCServer) SubscribeOnSportLines(stream api.Processor_SubscribeOnSportLinesServer) error {
	sportNames := []func(*api.Coefficients){}
	duration := int32(0)
	delta := false

	go s.subscribeOnSportLinesSender(stream, &sportNames, &duration, &delta)

	for {
		data, err := stream.Recv()
		if err != nil {
			return err
		}
		s.loger.Debug("SubscribeOnSportLines : gRPC SubscribeOnSportLines : sports:", data.Sports, "time:", data.Time)
		sportNames = []func(*api.Coefficients){}
		for _, sports := range data.Sports {
			if sports == "baseball" {
				sportNames = append(sportNames, s.getBaseball)
			}
			if sports == "football" {
				sportNames = append(sportNames, s.getFootball)
			}
			if sports == "soccer" {
				sportNames = append(sportNames, s.getSoccer)
			}
		}
		delta = false
		duration = data.Time
	}
}

// grpc stream for send - func`s for getting info - time for delay - send delta
func (s *gRPCServer) subscribeOnSportLinesSender(stream api.Processor_SubscribeOnSportLinesServer, sports *[]func(*api.Coefficients), t *int32, delta *bool) {
	firstCoefficients := make(map[string]float32)

	for {
		coefficients := api.Coefficients{}
		coefficients.Coefficients = make(map[string]float32)
		if len(*sports) == 0 {
			continue
		}
		for i := 0; i < len(*sports); i++ {
			(*sports)[i](&coefficients)
		}
		if *delta {
			for s, c := range firstCoefficients {
				coefficients.Coefficients[s] -= c
				coefficients.Coefficients[s] = float32(math.Ceil(float64(coefficients.Coefficients[s])*1000) / 1000)
			}
		} else {
			firstCoefficients = make(map[string]float32)
			for s, c := range coefficients.Coefficients {
				firstCoefficients[s] = c
			}
			*delta = true
		}
		err := stream.Send(&coefficients)
		if err != nil {
			s.loger.Error("GRPCServer : subscribeOnSportLinesSender : stream.Send : ", err)
		}
		time.Sleep(time.Second * time.Duration(*t))
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
