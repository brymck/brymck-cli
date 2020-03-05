package commands

import (
	"github.com/urfave/cli/v2"

	pb "github.com/brymck/brymck-cli/genproto/brymck/securities/v1"
	"github.com/brymck/brymck-cli/pkg"
	"github.com/brymck/brymck-cli/pkg/connections"
)

const securitiesServiceName = "securities-service"

type securitiesApi struct {
	client     pb.SecuritiesAPIClient
	connection *connections.Connection
}

func (api *securitiesApi) Close() {
	api.connection.Close()
}

func getSecuritiesApi(addr string) (*securitiesApi, error) {
	conn, err := connections.NewConnection(securitiesServiceName, addr)
	if err != nil {
		return nil, err
	}
	client := pb.NewSecuritiesAPIClient(conn.GrpcClientConnection)
	return &securitiesApi{client: client, connection: conn}, nil
}

func GetSecuritiesCommand() *cli.Command {
	var id uint64
	var addr string
	flags := []cli.Flag{
		&cli.Uint64Flag{
			Name:  "id",
			Usage: "security ID",
			Destination: &id,
		},
		&cli.StringFlag{
			Name:  "address",
			Usage: "address override",
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

					api, err := getSecuritiesApi(addr)
					if err != nil {
						return err
					}
					defer api.Close()

					resp, err := api.client.GetSecurity(api.connection.Context, req)
					if err != nil {
						return err
					}
					pkg.PrintAsJson(resp.Security)
					return nil
				},
			},
			{
				Name:  "update-prices",
				Usage: "update prices",
				Flags: flags,
				Action: func(c *cli.Context) error {
					req := &pb.UpdatePricesRequest{Id: id}

					api, err := getSecuritiesApi(addr)
					if err != nil {
						return err
					}
					defer api.Close()

					resp, err := api.client.UpdatePrices(api.connection.Context, req)
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
