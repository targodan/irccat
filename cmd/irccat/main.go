package main

import (
	"io"
	"irccat"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "irccat"
	app.Version = "0.1.0"

	config := &irccat.Config{}

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print only the version",
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "nick, n",
			Usage:       "nickname",
			Value:       "irccat",
			Destination: &config.Nick,
		},
		cli.StringFlag{
			Name:        "channel, c",
			Usage:       "the channel/user to send the message to, note that channels have to be prefixed by '#'",
			Value:       "#irccat",
			Destination: &config.Channel,
		},
		cli.BoolFlag{
			Name:        "tls",
			Usage:       "enables encryption",
			Destination: &config.UseTLS,
		},
		cli.IntFlag{
			Name:        "msgLimitPerSecond, l",
			Usage:       "maximum amount of messages sent per second",
			Value:       0,
			Destination: &config.MaxMessagesPerSecond,
		},
		cli.BoolFlag{
			Name:        "verbose, v",
			Usage:       "outputs more information",
			Destination: &config.Verbose,
		},
	}

	app.ArgsUsage = "<server> [filename]"

	app.Action = func(c *cli.Context) (err error) {
		if c.NArg() < 1 || c.NArg() > 2 {
			return cli.ShowAppHelp(c)
		}
		config.Server = c.Args().Get(0)

		var input io.ReadCloser
		needsClosing := false
		if c.NArg() < 2 {
			input = os.Stdin
		} else {
			infilename := c.Args().Get(1)
			if infilename == "-" {
				needsClosing = true
				input = os.Stdin
			} else {
				input, err = os.Open(infilename)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
			}
		}
		defer func() {
			if needsClosing {
				input.Close()
			}
		}()

		client := irccat.NewIRCClient(config)
		err = client.Connect()
		if err != nil {
			return cli.NewExitError(err, 2)
		}

		err = client.ConsumeAndSend(input)
		if err != nil {
			return cli.NewExitError(err, 3)
		}

		err = client.Close()
		if err != nil {
			return cli.NewExitError(err, 4)
		}

		return nil
	}

	app.Run(os.Args)
}
