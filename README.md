## Public Service Announcer

PSA can multiplex annoucements, important or otherwise.

It's a single binary that works as a CLI for sending messages from a
console, and can provide announcements as a service via e.g. a webhook.

## Usage

```shell
Usage of ./psa:
  -dryrun
    	validate and log but don't send anything
  -message string
    	the message to announce
  -v	verbose logging
```

## Configuration

For information on how to prepare your Discord server or Slack
workspace for receiving announcements, see the [Discord webhooks intro]
or [Slack webhooks API guide].

All configuration is optional, but without any there will be nowhere to
announce messages.

```shell
export PSA_DISCORD_WEBHOOK="https://discord.com/api/webhooks/.../..."
export PSA_SLACK_WEBHOOK="https://hooks.slack.com/services/.../.../..."
```

If you want to customize the way announcements are presented, you can
set `PSA_MSG_TEMPLATE` in your environment:

```shell
export PSA_MSG_TEMPLATE="📣 {{.Message}}"
```

`.Message` is (currently) the only field available to the template. If
your template does not use it, `psa` will not send anything.

[Discord webhooks intro]: https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks
[Slack webhooks API guide]: https://api.slack.com/messaging/webhooks

## TODO

- [x] CLI for sending
- [x] Discord announce
- [x] Slack announce
- [ ] incoming webhook (HMAC verification)
- [ ] Twitter announce
- [ ] configuration by simple file in `~/.config` (simplifying CLI use)
- [ ] more configurable announcement sinks (e.g. multiple Discord hooks)

## Maybe TODO

- Structured logging
- RBAC (users authenticated by e.g. Slack can be authorized to announce)
- Announce from Slack or Discord
- More integrations
  - IRC, email, text messages, Signal, Telegram, Cabal
