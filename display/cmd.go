package display

import (
	"errors"
	"time"

	ui "github.com/gizak/termui"
	"github.com/puigfp/observer/util"
	"github.com/urfave/cli"
)

// Command is the main function of this package
var Command = cli.Command{
	Name:    "display",
	Aliases: []string{"d"},
	Usage:   "Dislpaying the metrics stored in influxDB",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "path to an observer configuration file",
			Value: "config.json",
		},
	},
	Action: func(c *cli.Context) error {
		// read config
		config, err := util.ReadConfigJSON(c.String("config"))
		if err != nil {
			return err
		}
		if len(config.Websites) == 0 {
			return errors.New("no websites to monitor")
		}

		// init influxDB client
		influxdbClient, err := util.NewInfluxDBClient(config.InfluxDB)
		if err != nil {
			return err
		}

		// init UI
		if err := ui.Init(); err != nil {
			return err
		}
		defer ui.Close()

		// init state
		st := &state{}
		initState(config, st)

		// init widgets
		w := initBody(st)

		// register event listeners
		registerListeners(&w, st)

		// first UI render
		render()

		// first metrics poll
		updateState(influxdbClient, w, st)

		// start metrics polling loop
		go func() {
			for range time.Tick(time.Duration(10) * time.Second) {
				updateState(influxdbClient, w, st)
			}
		}()

		// keep process runnning
		ui.Loop()

		return nil
	},
}
