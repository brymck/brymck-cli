package commands

import (
	"time"

	"github.com/urfave/cli/v2"

	cal "github.com/brymck/brymck-cli/genproto/brymck/calendar/v1"
	dt "github.com/brymck/brymck-cli/genproto/brymck/dates/v1"
	"github.com/brymck/brymck-cli/pkg"
	"github.com/brymck/brymck-cli/pkg/connections"
)

const calendarServiceName = "calendar-service"

type calendarApi struct {
	client     cal.CalendarAPIClient
	connection *connections.Connection
}

func (api *calendarApi) Close() {
	api.connection.Close()
}

func getCalendarApi(addr string) (*calendarApi, error) {
	conn, err := connections.NewConnection(calendarServiceName, addr)
	if err != nil {
		return nil, err
	}
	client := cal.NewCalendarAPIClient(conn.GrpcClientConnection)
	return &calendarApi{client: client, connection: conn}, nil
}

func toProtoDate(text string) (*dt.Date, error) {
	date, err := time.Parse("2006-01-02", text)
	if err != nil {
		return nil, nil
	}
	year, month, day := date.Date()
	return &dt.Date{Year: int32(year), Month: int32(month), Day: int32(day)}, nil
}

func GetCalendarCommand() *cli.Command {
	var startDateText string
	var endDateText string
	var addr string

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "start-date",
			Usage:       "start date",
			Destination: &startDateText,
		},
		&cli.StringFlag{
			Name:        "end-date",
			Usage:       "end date",
			Destination: &endDateText,
		},
		&cli.StringFlag{
			Name:        "address",
			Usage:       "address override",
			Destination: &addr,
		},
	}

	return &cli.Command{
		Name:  "calendar",
		Usage: "calendar",
		Subcommands: []*cli.Command{
			{
				Name:  "get-dates",
				Usage: "get dates",
				Flags: flags,
				Action: func(c *cli.Context) error {
					start, err := toProtoDate(startDateText)
					if err != nil {
						return err
					}
					end, err := toProtoDate(endDateText)
					if err != nil {
						return err
					}
					req := &cal.GetDatesRequest{StartDate: start, EndDate: end}

					api, err := getCalendarApi(addr)
					if err != nil {
						return err
					}
					defer api.Close()

					resp, err := api.client.GetDates(api.connection.Context, req)
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
