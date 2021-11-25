Release Notes
=============

## 1.2.0

- `fault.UserError` and `fault.SystemError` implement the `Stringer` interface.
- `fault.SystemError` implements the `fmt.Formatter` interface.
- `fault.SystemError` exposes a new `Cause()` method to retrieve the underlying `error` object.

## 1.1.0

Added `fault.Userf`, `fault.Systemf` and `fault.SystemWrapf` functions.

## 1.0.0

Initial release of various packages.