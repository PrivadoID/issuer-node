# Formal verification with Gobra

[Gobra](https://github.com/viperproject/gobra) is a verifier for Go from ETH Zurich. It proves that a function
satisfies its specification for every possible input, not just the inputs a test happens to try. Specifications
are written next to the code, as comments starting with `// @`.

## What is verified

`internal/basicauth/basicauth.go`: `Allowed` decides whether a request may reach an API endpoint protected by
HTTP basic auth. `BasicAuthMiddleware` in `internal/api` calls it for every request. Gobra proves that:

1. Endpoints without basic auth are always reached.
2. Leaving both `ISSUER_API_AUTH_USER` and `ISSUER_API_AUTH_PASSWORD` empty disables basic auth.
3. **No bypass:** if either value is configured, a protected endpoint is reached only with credentials that
   match a fully configured user and password. A half-configured pair rejects every request.
4. So a request without credentials reaches a protected endpoint only when basic auth is disabled.
5. Requests with the configured credentials are accepted.

Gobra also proves that `Allowed` terminates and cannot panic.

## Running it

```shell
make verify                          # with docker
make verify CONTAINER_ENGINE=podman  # with podman
```

It takes a few seconds. Success ends with `Gobra found 0 errors.` and exit code 0. If a property doesn't hold,
Gobra names it and exits with code 1.

CI runs `make verify` on every push, in `.github/workflows/formal-verification.yml`.

## What is assumed, not proved

- `tools/gobra/stubs/crypto/subtle/subtle.gobra` gives `subtle.ConstantTimeCompare` the contract from its Go
  documentation: it returns 1 if the two slices have equal contents, and 0 otherwise.
- `constantTimeEqual` assumes that converting a string to `[]byte` and back yields the same string, as the Go
  specification defines.

## What is not covered

- The glue in `BasicAuthMiddleware`: reading the basic auth header with `r.BasicAuth()` and deciding whether an
  endpoint is protected. The generated router marks an endpoint as protected when `api/api.yaml` gives it
  `security: basicAuth`, so the proof is only as good as that list. The middleware's unit tests cover this glue.
- The startup check that rejects a half-configured pair (`checkEnvVars` in `internal/config`). It is tested,
  not proved.
- Timing. Gobra proves which requests are accepted, not how long the checks take.
