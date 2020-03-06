package commands

import (
	"github.com/brymck/helpers/services"
	"github.com/urfave/cli/v2"

	pb "github.com/brymck/brymck-cli/genproto/brymck/securities/v1"
	"github.com/brymck/brymck-cli/pkg"
)

func getSecuritiesApi() pb.SecuritiesAPIClient {
	return pb.NewSecuritiesAPIClient(services.MustConnect("securities-service"))
}

func GetSecuritiesCommand() *cli.Command {
	var id uint64
	var addr string
	flags := []cli.Flag{
		&cli.Uint64Flag{
			Name:        "id",
			Usage:       "security ID",
			Destination: &id,
		},
		&cli.StringFlag{
			Name:        "address",
			Usage:       "address override",
			Destination: &addr,
		},
	}

	return &cli.Command{
		Name:  "securities",
		Usage: "securities",
		Subcommands: []*cli.Command{
			{
				Name:  "get",
				Usage: "get a security by ID",
				Flags: flags,
				Action: func(c *cli.Context) error {
					req := &pb.GetSecurityRequest{Id: id}
					api := getSecuritiesApi()
					resp, err := api.GetSecurity(makeContext(), req)
					if err != nil {
						return err
					}
					pkg.PrintAsJson(resp.Security)
					return nil
				},
			},
			{
				Name:  "get-prices",
				Usage: "get prices",
				Flags: flags,
				Action: func(c *cli.Context) error {
					startDate := &pb.Date{Year: 2020, Month: 1, Day: 1}
					endDate := &pb.Date{Year: 2020, Month: 3, Day: 1}
					req := &pb.GetPricesRequest{Id: id, StartDate: startDate, EndDate: endDate}
					api := getSecuritiesApi()
					resp, err := api.GetPrices(makeContext(), req)
					if err != nil {
						return err
					}
					pkg.PrintAsJson(resp)
					return nil
				},
			},
			{
				Name:  "update-prices",
				Usage: "update prices",
				Flags: flags,
				Action: func(c *cli.Context) error {
					req := &pb.UpdatePricesRequest{Id: id}
					api := getSecuritiesApi()
					resp, err := api.UpdatePrices(makeContext(), req)
					if err != nil {
						return err
					}
					pkg.PrintAsJson(resp)
					return nil
				},
			},
		},
	}
}
