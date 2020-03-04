package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"

	pb "github.com/brymck/brymck-cli/genproto"
)

const ServiceAddressTemplate = "%s-4tt23pryoq-an.a.run.app:443"

type thing struct {
	context              *context.Context
	grpcClientConnection *grpc.ClientConn
	cancel               *context.CancelFunc
}

func (t *thing) Close() {
	t.grpcClientConnection.Close()
	(*t.cancel)()
}

func getThing(serviceName string) (*thing, error) {
	addr := fmt.Sprintf(ServiceAddressTemplate, serviceName)
	conn, err := getConnection(addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return &thing{context: &ctx, grpcClientConnection: conn, cancel: &cancel}, nil
}

func printAsJson(message proto.Message) {
	m := jsonpb.Marshaler{}
	result, _ := m.MarshalToString(message)
	fmt.Println(result)
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "securities",
				Usage: "securities",
				Subcommands: []*cli.Command{
					{
						Name:  "get",
						Usage: "get a security by ID",
						Action: func(c *cli.Context) error {
							id, err := strconv.Atoi(c.Args().First())
							if err != nil {
								return err
							}
							req := &pb.GetSecurityRequest{Id: int32(id)}

							x, err :=  getThing("securities-service")
							if err != nil {
								return err
							}
							defer x.Close()
							client := pb.NewSecuritiesAPIClient(x.grpcClientConnection)

							resp, err := client.GetSecurity(*x.context, req)
							if err != nil {
								return err
							}
							printAsJson(resp.Security)
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
