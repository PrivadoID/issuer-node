// Package basicauth decides whether a request may reach an API endpoint protected by HTTP basic auth.
//
// The "// @" comments are Gobra specifications. Running `make verify` makes Gobra prove that Allowed
// satisfies them for every possible input. See tools/gobra/README.md.
package basicauth

import "crypto/subtle"

// Allowed reports whether a request may reach the handler of an endpoint.
//
//   - protected: the endpoint requires basic auth in the API spec.
//   - user, pass: the credentials configured for the API. Leaving both empty disables basic auth.
//   - sent, sentUser, sentPass: whether the request carries basic auth credentials, and their values.
//
// Endpoints without basic auth are always reached:
// @ ensures !protected ==> res
//
// Leaving both configured values empty disables basic auth:
// @ ensures protected && user == "" && pass == "" ==> res
//
// No bypass: if either value is configured, a protected endpoint is reached only with credentials that
// match a fully configured user and password:
// @ ensures protected && (user != "" || pass != "") && res ==> user != "" && pass != "" && sent && sentUser == user && sentPass == pass
//
// So a request without credentials reaches a protected endpoint only when basic auth is disabled:
// @ ensures protected && !sent && res ==> user == "" && pass == ""
//
// Requests with the configured credentials are accepted:
// @ ensures protected && user != "" && pass != "" && sent && sentUser == user && sentPass == pass ==> res
// @ decreases
func Allowed(protected bool, user, pass string, sent bool, sentUser, sentPass string) (res bool) {
	if !protected || (user == "" && pass == "") {
		return true
	}
	if user == "" || pass == "" {
		// Only one value is configured: fail closed instead of skipping authentication.
		return false
	}
	if !sent {
		return false
	}
	return constantTimeEqual(user, sentUser) && constantTimeEqual(pass, sentPass)
}

// constantTimeEqual reports whether a and b are equal. Its running time depends on their lengths but not
// on their contents.
// @ ensures res == (a == b)
// @ decreases
func constantTimeEqual(a, b string) (res bool) {
	x, y := []byte(a), []byte(b)
	// Assumed by Gobra, not proved: by the Go spec, converting a string to []byte and back yields the
	// same string.
	// @ assume string(x) == a && string(y) == b
	return subtle.ConstantTimeCompare(x, y) == 1
}
