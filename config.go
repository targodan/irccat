package irccat

type Config struct {
	Nick                 string
	Server               string
	Password             string
	Channel              string
	UseTLS               bool
	MaxMessagesPerSecond int
	Verbose              bool
}
