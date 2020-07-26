/*
	В первый раз стартует фуннкция и отравлчяет полное значение
	Слудующее отправляет дельту
*/

package kiddylineprocessor

import (
	"log"
	"math"
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
	delta := false

	go s.subscribeOnSportLinesSender(&stream, &sportNames, &duration, &delta)

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
		delta = false
		duration = data.Time
	}
}

func (s *GRPSServer) subscribeOnSportLinesSender(stream *api.Processor_SubscribeOnSportLinesServer, sports *[]func(*api.Coefficients), t *int32, delta *bool) {
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
		s.loger.Trace("delta :", delta, *delta, " coefficients.Coefficients : ", coefficients.Coefficients, " firstCoefficients :", firstCoefficients)
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
		(*stream).Send(&coefficients)
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
