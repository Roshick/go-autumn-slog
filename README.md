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

Under construction
