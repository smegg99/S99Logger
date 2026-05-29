# S99Logger

`s99logger` is a small structured Go logger built around typed events,
human-readable console output, JSON logs, and optional localization.

Console logs are colored by default when writing to a terminal. JSON logs keep
stable machine fields such as `message_id` and `service`.

## Features

- Typed app-owned events.
- Console, JSON, and multi-sink in the core package.
- Integrated i18n and log rotation.
- Configurable console prefixes, service prefix, and colors.
- Optional context as the last log argument for per-request language.
- Translation coverage checks with `Translator.Verify`.

## Install

```sh
go get github.com/smegg99/s99logger
```

## Shape

Log calls are event-first. Context is optional and goes last when needed.

```go
event := s99logger.NewEvent("server_starting", s99logger.Int("port", 8080))
logger.Info(event)
logger.Info(event, s99logger.LanguageContext("pl"))
```

For named app events, wrap `NewEvent` in a small function. For plural
translations, use `s99logger.NewPluralEvent`.

Use `examples/` for complete, runnable code.

## Examples

Run any example with `go run ./examples/name`.

| Example | Shows |
| --- | --- |
| [`minimal`](./examples/minimal) | console logging with `NewEvent` |
| [`basic`](./examples/basic) | TOML translations with go-i18n |
| [`pluralization`](./examples/pluralization) | plural messages selected by event count |
| [`mixedformats`](./examples/mixedformats) | JSON, TOML, and YAML locale files together |
| [`localizedprefixes`](./examples/localizedprefixes) | language-specific console prefixes/colors |
| [`consoleprefixes`](./examples/consoleprefixes) | custom console prefixes, colors, and service style |
| [`jsonstdout`](./examples/jsonstdout) | newline-delimited JSON |
| [`multisink`](./examples/multisink) | console and JSON from one log call |
| [`rotatingfile`](./examples/rotatingfile) | optional `rotation` package for daily JSON files |
| [`scopedattrs`](./examples/scopedattrs) | `Logger.With` and `Options.MinLevel` |
| [`verifytranslations`](./examples/verifytranslations) | `Translator.Verify` coverage check |
| [`customsink`](./examples/customsink) | implementing `Sink` |

## License

MIT
