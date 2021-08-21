<h1 align="center">
  <img src="static/notify-logo.png" alt="notify" width="200px"></a>
  <br>
</h1>


<p align="center">
<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/license-MIT-_red.svg"></a>
<a href="https://github.com/projectdiscovery/notify/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat"></a>
<a href="https://goreportcard.com/badge/github.com/projectdiscovery/notify"><img src="https://goreportcard.com/badge/github.com/projectdiscovery/notify"></a>
<a href="https://github.com/projectdiscovery/notify/releases"><img src="https://img.shields.io/github/release/projectdiscovery/notify"></a>
<a href="https://hub.docker.com/r/projectdiscovery/notify"><img src="https://img.shields.io/docker/pulls/projectdiscovery/notify.svg"></a>
<a href="https://twitter.com/pdiscoveryio"><img src="https://img.shields.io/twitter/follow/pdiscoveryio.svg?logo=twitter"></a>
<a href="https://discord.gg/projectdiscovery"><img src="https://img.shields.io/discord/695645237418131507.svg?logo=discord"></a>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#notify-installation">Installation</a> •
  <a href="#provider-config">Providers</a> •
  <a href="#usage">Usage</a> •
  <a href="#running-notify">Running Notify</a> •
  <a href="#notes">Notes</a> •
  <a href="https://discord.gg/projectdiscovery">Join Discord</a>
</p>


Notify is a Go-based assistance package that enables you to stream the output of several tools (or read from a file) and publish it to a variety of supported platforms.

<h1 align="left">
  <img src="static/notify-httpx.png" alt="notify-httpx" width="700px"></a>
  <br>
</h1>

# Features

- Supports for Slack / Discord / Telegram
- Supports for Pushover / Email / Teams
- Supports for File / Pipe output
- Supports Line by Line / Bulk Post
- Supports using Single / Multiple providers
- Supports Custom Web-hooks
- Supports Custom data formatting 


# Usage

```sh
notify -h
```

This will display help for the tool. Here are all the switches it supports.

| Flag             | Description                                     | Example                           |
| ---------------- | ----------------------------------------------- | --------------------------------- |
| -config          | Notify configuration file                       | notify -config config.yaml        |
| -silent          | Don't print the banner                          | notify -silent                    |
| -version         | Show version of notify                          | notify -version                   |
| -v               | Show Verbose output                             | notify -v                         |
| -no-color        | Don't Use colors in output                      | notify -no-color                  |
| -data            | File path to read data from                     | notify -data test.txt             |
| -bulk            | Read and send data in bulk from file.           | notify -bulk                      |
| -char-limit      | Character limit for message (default 4000)      | notify -char-limit 2000           |
| -provider-config | provider config path                            | notify -provider-config conf.yaml |
| -provider        | provider to send the notification to (optional) | notify -provider slack,telegram   |
| -id              | id to send the notification to (optional)       | notify -id recon,scans            |


# Notify Installation

```sh
GO111MODULE=on go get -v github.com/projectdiscovery/notify/cmd/notify
```

### Provider Config

The default provider config file can be created at `$HOME/.config/notify/provider-config.yaml` and can have the following contents:

```yaml
slack:
  - id: "slack"
    slack_channel: "recon"
    slack_username: "test"
    slack_format: "{{data}}"
    slack_webhook_url: "https://hooks.slack.com/services/XXXXXX"

  - id: "vulns"
    slack_channel: "vulns"
    slack_username: "test"
    slack_format: "{{data}}"
    slack_webhook_url: "https://hooks.slack.com/services/XXXXXX"

discord:
  - id: "crawl"
    discord_channel: "crawl"
    discord_username: "test"
    discord_format: "{{data}}"
    discord_webhook_url: "https://discord.com/api/webhooks/XXXXXXXX"

  - id: "subs"
    discord_channel: "subs"
    discord_username: "test"
    discord_format: "{{data}}"
    discord_webhook_url: "https://discord.com/api/webhooks/XXXXXXXX"

telegram:
  - id: "tel"
    telegram_api_key: "XXXXXXXXXXXX"
    telegram_chat_id: "XXXXXXXX"
    telegram_format: "{{data}}"

pushover:
  - id: "push"
    pushover_user_key: "XXXX"
    pushover_api_token: "YYYY"
    pushover_format: "{{data}}"
    pushover_devices:
      - "iphone"

smtp:
  - id: email
    smtp_server: mail.example.com
    smtp_username: test@example.com
    smtp_password: password
    from_address: from@email.com
    smtp_cc:
      - to@email.com
    smtp_format: "{{data}}"

custom:
  - id: webhook
    custom_webook_url: http://host/api/webhook
    custom_method: GET
    custom_format: '{{data}}'
    custom_headers:
      Content-Type: application/json
      X-Api-Key: XXXXX
``` 

# Running Notify

Notify supports piping output of any tool or output file and send it to configured provider/s (e.g, discord, slack channel) as notification.

### Send notification using piped(stdin) output

```sh
subfinder -d hackerone.com | notify
```

<h1 align="left">
<img width="365" alt="notify-subfinder" src="https://user-images.githubusercontent.com/8293321/130240854-e3031bc6-ecc8-47f8-9654-4c58e09cc622.png">


### Send notification using output file


```sh
subfinder -d hackerone.com -o h1.txt; notify -data h1.txt
```

### Send notification using output file in bulk mode


```sh
subfinder -d hackerone.com -o h1.txt; notify -data h1.txt -bulk
```

### Send notification using output file to specific provider's


```sh
subfinder -d hackerone.com -o h1.txt; notify -data h1.txt -bulk -provider discord,slack
```

### Send notification using output file to specific ID's


```sh
subfinder -d hackerone.com -o h1.txt; notify -data h1.txt -bulk -id recon,vulns,scan
```

### Example Uses

Following command will enumerate subdomains using [SubFinder](https://github.com/projectdiscovery/subfinder) and probe alive URLs using [httpx](https://github.com/projectdiscovery/httpx), runs [Nuclei](https://github.com/projectdiscovery/nuclei) templates and send the nuclei results as a notifications to configured provider/s.


```sh
subfinder -d intigriti.com | httpx | nuclei -tags exposure -o output.txt; notify -bulk -data output.txt
```


### Provider Config


The tool tries to use the default provider config (`$HOME/.config/notify/provider-config.yaml`), it can also be specified via CLI by using **provider-config** flag.

To run the tool with custom providers config, just use the following command.

```sh
notify -provider-config providers.yaml
```

### Notify Config

Notify flags can be configured at default config (`$HOME/.config/notify/config.yaml`) or custom config can be also provided using `config` flag.

## Notes
- As default notify sends notification line by line
- **bulk** flag is supported with data flag.
- stdin/pipe input doesn't support **bulk** posting.

## References

- [Creating Slack webhook](https://slack.com/intl/en-it/help/articles/115005265063-Incoming-webhooks-for-Slack)
- [Creating Discord webhook](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks)
- [Creating Telegram bot](https://core.telegram.org/bots#3-how-do-i-create-a-bot)
- [Creating Pushover Token](https://github.com/containrrr/shoutrrr/blob/main/docs/services/pushover.md)

Notify is made with 🖤 by the [projectdiscovery](https://projectdiscovery.io) team.