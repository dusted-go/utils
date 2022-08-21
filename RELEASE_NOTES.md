Release Notes
=============

## 1.5.0

- `db.Get` will set the entity to `nil` if it cannot be found.

## 1.4.0

### Added several new packages

- `maps` for map related helper functions
- `stack` for capturing a human friendly stack trace
- `storage` for Google Cloud Storage related operations
- `webfile` for `multipart.File` related helper functions

### Changes

- The `fault` package doesn't have a dependency on `pkg/errors` anymore
- `fault.SystemError` doesn't implement `Cause` anymore (use `Unwrap`)
- `fault.SystemError` used the new `stack.Trace` object to capture the stack trace
- `db.NewService` requires a `context.Context` object now

## 1.3.0

The fault.SystemError implements the Unwrap method now.

## 1.2.0

- `fault.UserError` and `fault.SystemError` implement the `Stringer` interface.
- `fault.SystemError` implements the `fmt.Formatter` interface.
- `fault.SystemError` exposes a new `Cause()` method to retrieve the underlying `error` object.

## 1.1.0

Added `fault.Userf`, `fault.Systemf` and `fault.SystemWrapf` functions.

## 1.0.0

Initial release of various packages.