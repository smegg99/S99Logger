# S99Logger

`s99logger` is a small structured Go logger built around typed events,
human-readable console output, JSON logs, and optional localization.

Console logs are colored by default when writing to a terminal. JSON logs keep
stable machine fields such as `message_id` and `service`.

## Features

- Typed app-owned events.
- Console, JSON, and multi-sink in the core package.
- Process-wide default logger with a one-call `rotation.Configure` setup.
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

## Global logger

For apps that would rather not thread a `*Logger` through every call, the
package keeps a process-wide default. `rotation.Configure` builds a console
logger (plus a rotating JSON file sink when `EnableFiles` is set) and installs
it; `Close` releases the file sink at shutdown.

```go
if err := rotation.Configure(rotation.Config{
    Service:     "api",
    Level:       "info",
    EnableFiles: true,
    Directory:   "logs",
}); err != nil {
    panic(err)
}
defer s99logger.Close()

s99logger.Info(s99logger.NewEvent("server_starting", s99logger.Int("port", 8080)))
```

`s99logger.Default()` returns the installed logger when you need a `*Logger`
value (e.g. to derive a request-scoped child with `With`), and `SetDefault`
swaps it directly when you build sinks yourself. The default is usable before
configuration: it writes to stderr at debug level.

Use `examples/` for complete, runnable code.

## Examples

Run any example with `go run ./examples/name`.

| Example | Shows |
| --- | --- |
| [`minimal`](./examples/minimal) | console logging with `NewEvent` |
| [`global`](./examples/global) | process-wide default logger via `rotation.Configure` |
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
