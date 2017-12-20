package irccat

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	irc "gopkg.in/irc.v2"
)

type IRCClient struct {
	cfg       *Config
	conn      net.Conn
	client    *irc.Client
	ready     bool
	connected bool
	done      <-chan error
}

func NewIRCClient(cfg *Config) *IRCClient {
	return &IRCClient{
		cfg:       cfg,
		ready:     false,
		connected: false,
	}
}

func (c *IRCClient) Connect() error {
	if c.connected {
		return nil
	}
	c.connected = true

	clientCfg := irc.ClientConfig{
		Nick: c.cfg.Nick,
		Pass: c.cfg.Password,
		User: "IRCCAT",
		Name: "IRCCAT Bot",
		Handler: irc.HandlerFunc(func(ic *irc.Client, m *irc.Message) {
			if m.Command == "001" {
				if c.cfg.Channel[0] == '#' {
					ic.Write(fmt.Sprintf("JOIN %s", c.cfg.Channel))
				}
				c.ready = true
			} else if m.Command == "ERROR" {
				c.conn.Close()
			}
		}),
	}

	var err error

	if c.cfg.UseTLS {
		c.conn, err = tls.Dial("tcp", c.cfg.Server, &tls.Config{})
	} else {
		c.conn, err = net.Dial("tcp", c.cfg.Server)
	}
	if err != nil {
		return err
	}

	c.client = irc.NewClient(c.conn, clientCfg)

	if c.cfg.Verbose {
		c.client.Writer.DebugCallback = func(msg string) { fmt.Println("<- " + strings.TrimRight(msg, "\r\n")) }
		c.client.Reader.DebugCallback = func(msg string) { fmt.Println("-> " + strings.TrimRight(msg, "\r\n")) }
	}

	doneChan := make(chan error)
	c.done = doneChan
	go func(done chan<- error) {
		done <- c.client.Run()
	}(doneChan)

	for !c.ready {
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (c *IRCClient) ConsumeAndSend(input io.Reader) error {
	if !c.connected {
		return errors.New("the client needs to be connected first")
	}

	lastWaitTime := time.Now()
	numMessagesSinceLastWait := 0
	buf := bufio.NewReader(input)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		line = strings.Trim(line, " \n\r")
		if len(line) == 0 {
			line = " "
		}

		if c.cfg.MaxMessagesPerSecond > 0 && numMessagesSinceLastWait >= c.cfg.MaxMessagesPerSecond {
			waitTime := 1*time.Second - time.Since(lastWaitTime)
			time.Sleep(waitTime)
			lastWaitTime = time.Now()
			numMessagesSinceLastWait = 0
		}
		err = c.Send(line)
		numMessagesSinceLastWait++
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *IRCClient) Send(msg string) error {
	if !c.connected {
		return errors.New("the client needs to be connected first")
	}

	return c.client.WriteMessage(&irc.Message{
		Command: "PRIVMSG",
		Params: []string{
			c.cfg.Channel,
			msg,
		},
	})
}

func (c *IRCClient) Close() (err error) {
	c.client.WriteMessage(&irc.Message{
		Command: "QUIT",
		Params:  []string{"EOF"},
	})
	// Conn is closed upon receiving the ERROR answer to our QUIT
	select {
	case err = <-c.done:
	case <-time.After(500 * time.Millisecond):
	}
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
