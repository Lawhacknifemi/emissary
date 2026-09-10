package model

import (
	"encoding/base64"
	"net/url"
	"strings"
)

// remoteAccountIDPrefix marks a Mastodon account ID that encodes a remote actor's
// URL rather than a local hex ObjectID. "u_" can never begin a 24-char hex
// ObjectID, so a prefixed value is unambiguously a remote ID.
const remoteAccountIDPrefix = "u_"

// EncodeRemoteAccountID turns a remote actor's URL into an opaque, URL-path-safe
// Mastodon account ID. The ID is a pure function of the URL, so it stays stable
// even though the ascache row backing the actor is reminted on every refetch.
func EncodeRemoteAccountID(actorURL string) string {
	return remoteAccountIDPrefix + base64.RawURLEncoding.EncodeToString([]byte(actorURL))
}

// DecodeRemoteAccountID reverses EncodeRemoteAccountID. ok is false when s is not
// one of our encoded remote IDs (no prefix, not valid base64url, or the decoded
// value is not an absolute URL), so callers fall through to local-hex / bare-URL
// handling.
func DecodeRemoteAccountID(s string) (actorURL string, ok bool) {

	if !strings.HasPrefix(s, remoteAccountIDPrefix) {
		return "", false
	}

	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(s, remoteAccountIDPrefix))

	if err != nil {
		return "", false
	}

	if parsed, err := url.Parse(string(raw)); err != nil || !parsed.IsAbs() {
		return "", false
	}

	return string(raw), true
}
