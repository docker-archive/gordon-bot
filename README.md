# Gordon IRC Bot

Gordon the turtle IRC bot to restart jenkins PR builds.

```console
$ gordon-bot -h
Usage of gordon-bot:
  -channel="#docker-maintainers": irc channel
  -d=false: run in debug mode
  -nick="GordonTheTurtle": irc nick
  -pass="": irc pass
  -server="chat.freenode.net:6697": irc server
  -user="GordonTheTurtle": irc user
  -v=false: print version and exit (shorthand)
  -version=false: print version and exit

```

Example docker run command:

```bash
$ docker run -d --restart always \
    --name gordon-bot \
    -e LEEROY_USERNAME \
    -e LEEROY_PASS \
    dockercore/gordon-bot -d \
    -pass="YOUR_IRCPASS" 
```
