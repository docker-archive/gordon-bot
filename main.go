package main

import (
	"flag"
	"fmt"

	log "github.com/Sirupsen/logrus"
	leeroy "github.com/docker/gordon-bot/leeroy"
	bot "github.com/fabioxgn/go-bot"
	_ "github.com/fabioxgn/go-bot/commands/gif"
	_ "github.com/fabioxgn/go-bot/commands/godoc"
)

const (
	VERSION = "v0.1.0"
)

var (
	server  string
	channel string
	user    string
	nick    string
	pass    string
	debug   bool
	version bool
)

func init() {
	// parse flags
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&debug, "d", false, "run in debug mode")
	flag.StringVar(&server, "server", "irc.freenode.net:6697", "irc server")
	flag.StringVar(&channel, "channel", "#docker-maintainers", "irc channel")
	flag.StringVar(&user, "user", "GordonTheTurtle", "irc user")
	flag.StringVar(&nick, "nick", "GordonTheTurtle", "irc nick")
	flag.StringVar(&pass, "pass", "", "irc pass")
	flag.Parse()
}

func main() {
	// set log level
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if version {
		fmt.Println(VERSION)
		return
	}

	// Leeroy plugin needs the channel information so it can exclude any
	// rebuilding requests from undesired channels (e.g., #general when used
	// over a Slack bridge).
	leeroy.Channel = channel

	bot.Run(&bot.Config{
		Server:   server,
		Channels: []string{channel},
		User:     user,
		Nick:     nick,
		Password: pass,
		UseTLS:   true,
		Debug:    debug,
	})
}
