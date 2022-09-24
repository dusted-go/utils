Release Notes
=============

## 1.13.0

- Upgraded dependencies

## 1.12.0

- Improved the output from `(u URL) Pretty()`

## 1.11.1

- Fixed a bug in `webfile.MimeType`

## 1.11.0

- Improved `webfile.MimeType` to get the mime type correctly for SVGs and other files as well.

## 1.10.0

- Improved `mapsort.KeyByValue` to be more generic

## 1.9.0

- Renamed `db.Client` to `db.Repo` to better reflect the latest re-factorings
- Added more tests for db stuff

## 1.8.0

- Moved `stack` to `github.com/dusteg-go/fault/stack`
- Moved `fault` to `github.com/dusteg-go/fault/fault`

## 1.7.0

- Renamed `db.Service` to `db.Client`
- Renamed `storage.Service` to `storage.Client`
- Renamed `mailer` pkg to `mailman`
- Renamed `mailer.Mailer` to `mailman.Client`
- Re-factored `db.Service`/`db.Client` to provide better UX through generics
- Renamed `maps` pkg to `mapsort`
- Renamed `email.Address` to `typ.Email`
- Added `typ.URL` type

## 1.6.0

- Changed `db.Get` interface to include `bool` flag indicating if the entity could be found.

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