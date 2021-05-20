## Public Service Announcer

PSA can multiplex annoucements, important or otherwise.

It's a single binary that works as a CLI for sending messages from a
console, and can provide announcements as a service via e.g. a webhook.

## Usage

```shell
Usage of ./psa:
  -message string
    	the message to announce
  -template string
    	template to apply when announcing (default "📣 {{.Message}}")
  -v	verbose logging
```

## Configuration

For information on how to prepare your Discord server or Slack
workspace for receiving announcements, see the [Discord webhooks intro]
or [Slack webhooks API guide].

All configuration is optional, but without any there will be nowhere to
announce messages.

```shell
export DISCORD_WEBHOOK="https://discord.com/api/webhooks/.../..."
export SLACK_WEBHOOK="https://hooks.slack.com/services/.../.../..."
```

[Discord webhooks intro]: https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks
[Slack webhooks API guide]: https://api.slack.com/messaging/webhooks

## TODO

[x] CLI for sending
[x] Discord announce
[x] Slack announce
[ ] incoming webhook (HMAC verification)
[ ] Twitter announce
[ ] configuration by simple file in `~/.config` (simplifying CLI use)
[ ] more configurable announcement sinks (e.g. multiple Discord hooks)

## Maybe TODO

- RBAC (users authenticated by e.g. Slack can be authorized to announce)

- Announce from Slack or Discord

- More integrations
  - IRC, email, text messages, Signal, Telegram, Cabal
