package commands

import (
	"strconv"
	"strings"

	"github.com/brymck/helpers/services"
	"github.com/urfave/cli/v2"

	pb "github.com/brymck/brymck-cli/genproto/brymck/risk/v1"
	"github.com/brymck/brymck-cli/pkg"
)

func getRiskApi() pb.RiskAPIClient {
	return pb.NewRiskAPIClient(services.MustConnect("risk-service"))
}

func GetRiskCommand() *cli.Command {
	var id uint64
	var idsString string
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
		Name:  "risk",
		Usage: "risk",
		Subcommands: []*cli.Command{
			{
				Name:  "get-risk",
				Usage: "get security risk by ID",
				Flags: flags,
				Action: func(c *cli.Context) error {
					req := &pb.GetRiskRequest{SecurityId: id}
					api := getRiskApi()
					resp, err := api.GetRisk(makeContext(), req)
					if err != nil {
						return err
					}
					pkg.PrintAsJson(resp)
					return nil
				},
			},
			{
				Name:  "get-covariances",
				Usage: "get covariances by security IDs",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "ids",
						Usage:       "comma-delimited list of security IDs",
						Destination: &idsString,
					},
					&cli.StringFlag{
						Name:        "address",
						Usage:       "address override",
						Destination: &addr,
					},
				},
				Action: func(c *cli.Context) error {
					parts := strings.Split(idsString, ",")
					ids := make([]uint64, len(parts))
					for i, part := range parts {
						n, err := strconv.Atoi(part)
						if err != nil {
							return err
						}
						ids[i] = uint64(n)
					}

					req := &pb.GetCovariancesRequest{SecurityIds: ids}

					api := getRiskApi()
					resp, err := api.GetCovariances(makeContext(), req)
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
