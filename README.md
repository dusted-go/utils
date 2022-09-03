# Utils

Various utility packages and functions used in personal Go projects.

It's not recommended to be used by other people. The code is very likely to be unhelpful to other projects than my own.

## Packages

### Clients

- `db`: Generic repository for Google Cloud Datastore.
- `hcaptcha`: Verify a submitted captcha response with hCaptcha.
- `mailman`: Helper pkg to interact with a privately hosted email service.
- `storage`: Helper client to interact with Google Cloud Storage (mostly tailored to my own needs).

### Types

- `typ`: Package with additional rich types like emails and URLs (different to `url.URL`).

### Misc

- `mapsort`: Generic sort functions for maps.
- `webfile`: Utility functions to deal with user submitted `multipart.file` objects.