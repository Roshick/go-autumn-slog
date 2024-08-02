# go-autumn-slog

Inspired by the go-autumn library, go-autumn-slog offers an implementation of the [go-autumn-logging](https://github.com/StephanHCB/go-autumn-logging) interface. It utilizes the standard Go library for structured logging, log/slog, ensuring seamless integration and compatibility.

## Features

- **Context-Aware Logging**: Enhance logs with contextual information.
- **Seamless Integration with slog**: Utilizes slog.Logger and slog.Handler.
- **Respects slog.Default**: Adheres to the default logger for simplicity.
- **Extended Logging Levels**: Provides finer granularity with additional log levels.
- **Callback Support**: Allows custom handling via slog.Handler with callbacks.
- **ConfigLoader Compatibility**: Instantiates slog.HandlerOptions compatible with [go-autumn-configloader](https://github.com/Roshick/go-autumn-configloader).

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
    - [Logging](#logging)
        - [Resources via ConfigLoader](#resources-via-configloader)
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

#### Resources via ConfigLoader

The package provides configuration-based instantiation of various slog resources, compatible with [go-autumn-configloader](https://github.com/Roshick/go-autumn-configloader).

Supported resources include `slog.HandlerOptions`, which are used by `slog.TextHandler`, `slog.JSONHandler`, and third-party handlers. This allows users to define log levels and manipulate record attributes, useful in standardized logging scenarios, like those defined by the [Elastic Common Schema](https://www.elastic.co/guide/en/ecs/current/index.html). The `slog.HandlerOptions` also map new log levels to their respective string values, preventing incorrect level mappings.

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

## Contributing

Contributions are welcome! Please fork the repository, submit a pull request, or open an issue for any bugs or feature requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
