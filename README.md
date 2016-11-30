# Slacker

Allows people to send messages as a bot in a number of Slack channels.

# Usage

```shell
go get github.com/porty/slacker
slacker
```

## Required Environment Variables

Environment variables are specified in `config.go`.

```
# Retrieved from Incoming Webhook settings
SLACKER_SLACKMESSAGEURL=https://hooks.slack.com/services/Txxxxx/Bxxxxx/xxxxx
# Channels to show in the select dropdown
SLACKER_SLACKCHANNELS=bots,general
# Username for basic auth
SLACKER_USERNAME=user
# Password for basic auth
SLACKER_PASSWORD=pass
```

## Docker

Example Docker workings are in `docker/`.
You will likely have to edit them to work to your liking.

```shell
# add some environment variables
vi .env
# ca-certificates.crt are required from somewhere
curl https://curl.haxx.se/ca/cacert.pem -o ca-certificates.crt
# add -d to daemonise
docker-compose up
```

As its set to listen on localhost, you might want to either proxy the connection or change the bind addrees by editing the port declaration in `docker-compose.yml`.

# Development

`go-bindata` is required to regenerate `bindata.go` from files in `templates/`.