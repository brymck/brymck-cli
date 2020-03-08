package commands

import (
	"github.com/brymck/helpers/services"
	"github.com/urfave/cli/v2"

	dt "github.com/brymck/genproto/brymck/dates/v1"
	sec "github.com/brymck/genproto/brymck/securities/v1"

	"github.com/brymck/brymck-cli/pkg"
)

func getSecuritiesApi() sec.SecuritiesAPIClient {
	return sec.NewSecuritiesAPIClient(services.MustConnect("securities-service"))
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
					req := &sec.GetSecurityRequest{Id: id}
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
					startDate := &dt.Date{Year: 2020, Month: 1, Day: 1}
					endDate := &dt.Date{Year: 2020, Month: 3, Day: 1}
					req := &sec.GetPricesRequest{Id: id, StartDate: startDate, EndDate: endDate}
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
					req := &sec.UpdatePricesRequest{Id: id}
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
