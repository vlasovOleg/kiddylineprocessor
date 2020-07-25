/*
	В первый раз стартует фуннкция и отравлчяет полное значение
	Слудующее отправляет дельту
*/

package kiddylineprocessor

import (
	"log"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vlasovoleg/kiddyLineProcessor/internal/store"
	"github.com/vlasovoleg/kiddyLineProcessor/pkg/api"
	"google.golang.org/grpc"
)

// GRPSServer ...
type GRPSServer struct {
	store *store.Store
	loger *logrus.Logger
}

// NewGRPS ..
func (kp *Kiddylineprocessor) NewGRPS(store *store.Store, loger *logrus.Logger) {
	kp.loger.Trace("Kiddylineprocessor : NewGRPS...")
	lis, err := net.Listen("tcp", kp.config.GRPCserverAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	g := &GRPSServer{
		store: store,
		loger: loger,
	}
	api.RegisterProcessorServer(grpcServer, g)

	grpcServer.Serve(lis)
}

// SubscribeOnSportLines ...
func (s *GRPSServer) SubscribeOnSportLines(stream api.Processor_SubscribeOnSportLinesServer) error {
	sportNames := []func(*api.Coefficients){}
	duration := int32(0)

	go s.subscribeOnSportLinesSender(&stream, &sportNames, &duration)

	for {
		data, err := stream.Recv()
		if err != nil {
			return err
		}
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
		duration = data.Time
	}
}

func (s *GRPSServer) subscribeOnSportLinesSender(stream *api.Processor_SubscribeOnSportLinesServer, sports *[]func(*api.Coefficients), t *int32) {
	for {
		c := api.Coefficients{}
		c.Coefficients = make(map[string]float32)
		if len(*sports) != 0 {
			for i := 0; i < len(*sports); i++ {
				(*sports)[i](&c)
			}
			(*stream).Send(&c)
		}
		time.Sleep(time.Second * time.Duration(*t))
	}
}

func (s *GRPSServer) getBaseball(c *api.Coefficients) {
	dbc, err := (*s.store).BaseballRepository().GetCoefficient()
	if err != nil {
		s.loger.Error("kiddylineprocessor : gRPS API : getBaseball : BaseballRepository.GetCoefficient : " + err.Error())
	}
	c.Coefficients["baseball"] = dbc
}

func (s *GRPSServer) getFootball(c *api.Coefficients) {
	dbc, err := (*s.store).FootballRepository().GetCoefficient()
	if err != nil {
		s.loger.Error("kiddylineprocessor : gRPS API : getFootball : FootballRepository.GetCoefficient : " + err.Error())
	}
	c.Coefficients["football"] = dbc
}

func (s *GRPSServer) getSoccer(c *api.Coefficients) {
	dbc, err := (*s.store).SoccerRepository().GetCoefficient()
	if err != nil {
		s.loger.Error("kiddylineprocessor : gRPS API : getSoccer : SoccerRepository.GetCoefficient : " + err.Error())
	}
	c.Coefficients["soccer"] = dbc
}
