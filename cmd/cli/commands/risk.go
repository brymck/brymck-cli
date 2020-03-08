package commands

import (
	"strconv"
	"strings"

	pb "github.com/brymck/genproto/brymck/risk/v1"
	"github.com/brymck/helpers/services"
	"github.com/urfave/cli/v2"

	"github.com/brymck/brymck-cli/pkg"
)

func getRiskApi(addr string) pb.RiskAPIClient {
	if addr != "" {
		return pb.NewRiskAPIClient(services.MustConnectLocally(addr))
	} else {
		return pb.NewRiskAPIClient(services.MustConnect("risk-service"))
	}
}

func getFrequency(monthly bool) pb.Frequency {
	if monthly {
		return pb.Frequency_FREQUENCY_MONTHLY
	} else {
		return pb.Frequency_FREQUENCY_DAILY
	}
}

func GetRiskCommand() *cli.Command {
	var id uint64
	var idsString string
	var addr string
	var monthly bool
	flags := []cli.Flag{
		&cli.Uint64Flag{
			Name:        "id",
			Usage:       "security ID",
			Destination: &id,
		},
		&cli.BoolFlag{
			Name:        "monthly",
			Usage:       "use monthly returns",
			Destination: &monthly,
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
					freq := getFrequency(monthly)
					req := &pb.GetRiskRequest{SecurityId: id, Frequency: freq}
					api := getRiskApi(addr)
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

					freq := getFrequency(monthly)
					req := &pb.GetCovariancesRequest{SecurityIds: ids, Frequency: freq}

					api := getRiskApi(addr)
					resp, err := api.GetCovariances(makeContext(), req)
					if err != nil {
						return err
					}
					pkg.PrintAsJson(resp)
					return nil
				},
			},
			{
				Name:  "get-returns",
				Usage: "get return time series by security ID",
				Flags: flags,
				Action: func(c *cli.Context) error {
					freq := getFrequency(monthly)
					req := &pb.GetReturnTimeSeriesRequest{SecurityId: id, Frequency: freq}
					api := getRiskApi(addr)
					resp, err := api.GetReturnTimeSeries(makeContext(), req)
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
