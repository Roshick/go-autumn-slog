# go-autumn-slog

Built upon the inspiration from go-autumn, this library offers an implementation of
the [go-autumn-logging](https://github.com/StephanHCB/go-autumn-logging) interface. It leverages the standard Go library
for structured logging, log/slog.

## About go-autumn-slog

This library seamlessly integrates with the log/slog library, ensuring 100% compatibility and enabling users to utilize
every existing [handler](https://pkg.go.dev/log/slog#hdr-Writing_a_handler).

### Features

* provides context-aware logging
* consumes slog.Logger and therefore slog.Handler
* respects slog.Default
* extends slog.Level for finer granularity
* provides slog.Handler with callback support
* provides [go-autumn-configloader](https://github.com/Roshick/go-autumn-configloader) compatible instantiation of
  slog.HandlerOptions

## Usage

The library is divided into four largely independent areas: 'logging', 'level', 'handlers' and 'handleroptions'.

### Logging

A simple invocation of `myLogging := logging.New()` is all that's required to generate a functional instance of the
logging system.

Given a slog.Handler 'myHandler' and a *slog.Logger `myLogger := slog.New(myHandler)`, the preferred logger utilized by
the system can be established in the following order of priority:

1. By adding the logger to a context using `mySubCtx := logging.ContextWithLogger(ctx, myLogger)`.
2. By setting the logger as the default for a new instance through `mySubLogging := myLogging.WithLogger(myLogger)`.
3. By setting the logger as slog's default logger via `slog.SetDefault(myLogger)`.

### Level

This library expands the default slog.Levels to include the following severity levels, arranged in ascending order:

1. Trace
2. Debug
3. Info
4. Warn
5. Error
6. Fatal
7. Panic
8. Silent

### Handlers

The custom handlers provided by this library are designed to be used as standalone components and are entirely
independent of the need to utilize the entire logging implementation.

#### Callback

This handler functions as a wrapper for another handler, enabling the registration of callback
functions of the type `func(ctx context.Context, record *slog.Record) error`. These
callbacks have the ability to access the current context and manipulate the current slog.Record. This handler is
particularly useful for adding context values to every record, such as tracing information.

#### Noop

This handler performs no operations and is employed in situations where neither the context, the logging system, nor
slog has any logger configured.

### HandlerOptions

This package provides configuration-based instantiation (compatible with but not limited
to [go-autumn-configloader](https://github.com/Roshick/go-autumn-configloader)) of slog.HandlerOptions, a feature
utilized by slog.TextHandler,
slog.JSONHandler, and various third-party handlers. Users can leverage these options to define log levels for handlers
and manipulate the attributes of each passing record.

This flexibility is especially beneficial in scenarios demanding standardized log fields, as demonstrated by
the [Elastic Common Schema](https://www.elastic.co/guide/en/ecs/current/index.html). Additionally, the supplied
slog.HandlerOptions map the new log levels to their respective correct string
values. This proves crucial, given that log/slog defaults to mapping them in relation to the default values —
illustrated
by, for instance, 'PANIC' being emitted as 'ERROR+8'.

## Examples

This section delves into practical use cases to aid the integration of the library into various applications.

### Tracing

Consider a scenario where we aim to append tracing information to every log record generated during a request to one of
our services. Assuming we have incorporated a middleware (either third-party or self-written) that appends a trace-id
and span-id to our context, the next step is to propagate this information to our logger.

With this library, there are essentially two recommended approaches to achieve our goal.

#### The Middleware Approach

If we can register our middleware to execute after the one responsible for adding tracing information, and if this
information remains constant throughout the context's lifetime, we can directly include this data into a sub-logger
within our middleware and attach it to our context:

```
myLogger := logging.FromContext(ctx)
if myLogger != nil {
    // if context has no logger attached, obtain logger (e.g. from slog.Default, logging instance or simply create a new one)
    myHandler := slog.NewJSONHandler(os.Stdout, nil)
    myLogger = slog.New(myHandler)
}
myLogger = myLogger.With("trace-id", ctx.Value("trace-id"), "span-id", ctx.Value("span-id"))
ctx = logging.ContextWithLogger(ctx, myLogger)
```

#### The Callback Handler Approach

If the information within the context might change during its lifetime, opting for the callback handler, despite being
slightly slower, offers a safer solution.

```
myHandler := slog.NewJSONHandler(os.Stdout, nil)
myCallbackHandler := callbackhandler.New(myHandler)
err := myCallbackHandler.RegisterContextCallback(func(ctx context.Context, record *slog.Record) error {
    record.Add("span-id", ctx.Value("span-id"), "request-id", ctx.Value("request-id"))
    return nil
}, "add-tracing-attributes")
if err != nil {
    return err
}
myLogger := slog.New(myCallbackHandler)
// use logger as slog.Default, add it to our logging instance or attach it to a context
slog.SetDefault(myLogger)
```
