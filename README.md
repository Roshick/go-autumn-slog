# go-autumn-slog

Inspired by the go-autumn library, go-autumn-slog offers an implementation of the [go-autumn-logging](https://github.com/StephanHCB/go-autumn-logging) interface. It utilizes the standard Go library for structured logging, log/slog, ensuring seamless integration and compatibility.

## Features

- **Context-Aware Logging**: Enhance logs with contextual information.
- **Seamless Integration with slog**: Utilizes slog.Logger and slog.Handler.
- **Respects slog.Default**: Adheres to the default logger for simplicity.
- **Extended Logging Levels**: Provides finer granularity with additional log levels.
- **Callback Support**: Allows custom handling via slog.Handler with callbacks.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
    - [Logging](#logging)
        - [Configuration via Environment Variables](#configuration-via-environment-variables)
    - [Level](#level)
    - [Handlers](#handlers)
        - [Callback](#callback)
        - [Noop](#noop)
- [Examples](#examples)
    - [Simple Plaintext Logging](#simple-plaintext-logging)
    - [Structured JSON Logging and Context Awareness](#structured-json-logging-and-context-awareness)
    - [Tracing](#tracing)
        - [The Middleware Approach](#the-middleware-approach)
        - [The Callback Handler Approach](#the-callback-handler-approach)
    - [Environment-Based Logger Creation](#environment-based-logger-creation)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install go-autumn-slog, use the following command:

```sh
go get github.com/Roshick/go-autumn-slog
```

## Usage

go-autumn-slog is divided into three primary areas: Loggers, Levels, and Handlers.

### Logging

Creating a logging instance is simple:

```go
myLogging := logging.New()
```

Given a `slog.Handler` named `myHandler` and a `*slog.Logger` named `myLogger`:

```go
myLogger := slog.New(myHandler)
```

Setting the preferred logger can be prioritized in three ways:

1. Attach the logger to a context:
   ```go
   mySubCtx := logging.ContextWithLogger(ctx, myLogger)
   ```

2. Set the logger as the default for a new instance:
   ```go
   mySubLogging := myLogging.WithLogger(myLogger)
   ```

3. Set the logger as slog's default logger:
   ```go
   slog.SetDefault(myLogger)
   ```

### Level

Expands the default slog levels to include:

1. **Trace**
2. **Debug**
3. **Info**
4. **Warn**
5. **Error**
6. **Fatal**
7. **Panic**
8. **Silent**

### Handlers

Custom handlers provided by this library can be used independently of the entire logging implementation.

#### Callback

A handler that wraps another handler, allowing registration of callback functions of the type `func(ctx context.Context, record *slog.Record) error`. These callbacks can modify the slog record and add context values, useful for adding consistent tracing information.

#### Noop

A no-op handler that performs no operations, useful in cases where there is no configured logger in the context, the logging system, or slog.

### Configuration

Configuration for go-autumn-slog can be done through environment variables, allowing for flexible and dynamic log management.

#### Configuration via Environment Variables

The package supports configuration through environment variables using the `caarlos0/env` library.

Set the following environment variables to configure logging:

- `LOG_STYLE`: The log output style ("PLAIN" or "JSON").
- `LOG_LEVEL`: The minimum log level (e.g., "INFO", "DEBUG", "FATAL").
- `LOG_TIME_TRANSFORMER`: How to transform timestamps ("UTC" or "ZERO").
- `LOG_ATTRIBUTE_KEY_MAPPINGS`: A JSON string mapping attribute keys (e.g., `{"time":"@timestamp","level":"log.level","msg":"message","error":"error.message"}`).

Example:

```bash
export LOG_STYLE="JSON"
export LOG_LEVEL="INFO"
export LOG_TIME_TRANSFORMER="UTC"
export LOG_ATTRIBUTE_KEY_MAPPINGS='{"time":"@timestamp","level":"log.level","msg":"message","error":"error.message"}'
```

To load the configuration and create a logger:

```go
logger, err := logging.NewLoggerFromEnv()
if err != nil {
    panic("failed to create logger: " + err.Error())
}
slog.SetDefault(logger)
```

## Examples

Below are practical examples to help integrate the library into various applications.

### Simple Plaintext Logging

```go
// Create a simple plaintext logger and set the autumn global default logger to it
aulogging.Logger = logging.New()

// Use it
aulogging.Logger.NoCtx().Info().Print("hello")
```

### Structured JSON Logging and Context Awareness

```go
// Build a structured logger using slog.NewJSONHandler
structuredLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

// Set the autumn global default logger to it
aulogging.Logger = logging.New().WithLogger(structuredLogger)

// Use it (chain style or convenience style)
aulogging.Logger.NoCtx().Info().Print("hello")
aulogging.Info(context.Background(), "hello, too")

// Augment the logger with an extra field
augmentedLogger := structuredLogger.With("some-field", "some-value")

// Place it in a context
ctx := context.Background()
ctx = logging.ContextWithLogger(ctx, augmentedLogger)

// Use it from the context (chain style or convenience style)
aulogging.Logger.Ctx(ctx).Info().Print("hi")
aulogging.Info(ctx, "hi, again")
```

### Tracing

Add tracing information to every log record generated during a request to a service:

#### The Middleware Approach

Add constant tracing information within a middleware:

```go
myLogger := logging.FromContext(ctx)
if myLogger == nil {
    myHandler := slog.NewJSONHandler(os.Stdout, nil)
    myLogger = slog.New(myHandler)
}
myLogger = myLogger.With("trace-id", ctx.Value("trace-id"), "span-id", ctx.Value("span-id"))
ctx = logging.ContextWithLogger(ctx, myLogger)
```

#### The Callback Handler Approach

For dynamic context information, use the callback handler:

```go
myHandler := slog.NewJSONHandler(os.Stdout, nil)
myCallbackHandler := callbackhandler.New(myHandler)
err := myCallbackHandler.RegisterContextCallback(func(ctx context.Context, record *slog.Record) error {
    record.Add("span-id", ctx.Value("span-id"), "request-id", ctx.Value("request-id"))
    return nil
}, "add-tracing-attributes")
if err != nil {
    log.Fatalf("Failed to register callback: %v", err)
}
myLogger := slog.New(myCallbackHandler)
slog.SetDefault(myLogger)
```

### Environment-Based Logger Creation

Create a logger directly from environment variables:

```go
logger, err := logging.NewLoggerFromEnv()
if err != nil {
    panic("failed to create logger: " + err.Error())
}
slog.SetDefault(logger)

// Now use slog as usual
slog.Info("This is an info message")
```

## Contributing

Contributions are welcome! Please fork the repository, submit a pull request, or open an issue for any bugs or feature requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
