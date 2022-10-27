package app

import (
	"context"
	"net"
	"reflect"
	"strings"

	"github.com/thangchung/go-coffeeshop/cmd/product/config"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
	"google.golang.org/grpc"
)

type App struct {
	logger  *mylogger.Logger
	cfg     *config.Config
	network string
	address string
}

type ProductServiceServerImpl struct {
	gen.UnimplementedProductServiceServer
	logger *mylogger.Logger
}

var ItemTypes = map[string]gen.ItemTypeDto{
	"0": {
		Name:  "CAPPUCCINO",
		Type:  0,
		Price: 4.5,
	},
	"1": {Name: "COFFEE_BLACK",
		Type:  1,
		Price: 3,
	},
	"2": {
		Name:  "COFFEE_WITH_ROOM",
		Type:  2,
		Price: 3,
	},
	"3": {
		Name:  "ESPRESSO",
		Type:  3,
		Price: 3.5,
	},
	"4": {
		Name:  "ESPRESSO_DOUBLE",
		Type:  4,
		Price: 4.5,
	},
	"5": {
		Name:  "LATTE",
		Type:  5,
		Price: 4.5,
	},
	"6": {
		Name:  "CAKEPOP",
		Type:  6,
		Price: 2.5,
	},
	"7": {
		Name:  "CROISSANT",
		Type:  7,
		Price: 3.25,
	},
	"8": {
		Name:  "MUFFIN",
		Type:  8,
		Price: 3,
	},
	"9": {
		Name:  "CROISSANT_CHOCOLATE",
		Type:  9,
		Price: 3.5,
	},
}

func (g *ProductServiceServerImpl) GetItemTypes(ctx context.Context, request *gen.GetItemTypesRequest) (*gen.GetItemTypesResponse, error) {
	g.logger.Info("GET: GetItemTypes")

	res := gen.GetItemTypesResponse{}

	for _, v := range ItemTypes {
		res.ItemTypes = append(res.ItemTypes, &gen.ItemTypeDto{
			Name: v.Name,
			Type: v.Type,
		})
	}

	return &res, nil
}

func (g *ProductServiceServerImpl) GetItemsByType(ctx context.Context, request *gen.GetItemsByTypeRequest) (*gen.GetItemsByTypeResponse, error) {
	g.logger.Info("GET: GetItemsByType")

	res := gen.GetItemsByTypeResponse{}

	itemTypes := strings.Split(request.ItemTypes, ",")
	for _, itemType := range itemTypes {
		item := ItemTypes[itemType]
		if !reflect.DeepEqual(item, gen.ItemTypeDto{}) {
			res.Items = append(res.Items, &gen.ItemDto{
				Price: item.Price,
				Type:  item.Type,
			})
		}
	}

	return &res, nil
}

func New(log *mylogger.Logger, cfg *config.Config) *App {
	return &App{
		logger:  log,
		cfg:     cfg,
		network: "tcp",
		address: "0.0.0.0:5001",
	}
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Init %s %s\n", a.cfg.Name, a.cfg.Version)

	// Repository
	// ...

	// Use case
	// ...

	// gRPC Server
	l, err := net.Listen(a.network, a.address)
	if err != nil {
		return err
	}

	defer func() {
		if err := l.Close(); err != nil {
			a.logger.Error("Failed to close %s %s: %v", a.network, a.address, err)
		}
	}()

	s := grpc.NewServer()
	gen.RegisterProductServiceServer(s, &ProductServiceServerImpl{logger: a.logger})

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	a.logger.Info("Start server at " + a.address + " ...")

	return s.Serve(l)
}
