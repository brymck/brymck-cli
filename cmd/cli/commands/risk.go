package commands

import (
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"

	pb "github.com/brymck/brymck-cli/genproto/brymck/risk/v1"
	"github.com/brymck/brymck-cli/pkg"
	"github.com/brymck/brymck-cli/pkg/connections"
)

const riskServiceName = "risk-service"

type riskApi struct {
	client     pb.RiskAPIClient
	connection *connections.Connection
}

func (api *riskApi) Close() {
	api.connection.Close()
}

func getRiskApi(addr string) (*riskApi, error) {
	conn, err := connections.NewConnection(riskServiceName, addr)
	if err != nil {
		return nil, err
	}
	client := pb.NewRiskAPIClient(conn.GrpcClientConnection)
	return &riskApi{client: client, connection: conn}, nil
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

					api, err := getRiskApi(addr)
					if err != nil {
						return err
					}
					defer api.Close()

					resp, err := api.client.GetRisk(api.connection.Context, req)
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

					api, err := getRiskApi(addr)
					if err != nil {
						return err
					}
					defer api.Close()

					resp, err := api.client.GetCovariances(api.connection.Context, req)
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
