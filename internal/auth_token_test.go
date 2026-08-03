package internal

import (
	"path/filepath"
	"testing"
	"time"
)

func writeAnilistTokenFileForTest(t *testing.T, dir, token string) {
	t.Helper()
	tok := &AnilistToken{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   31536000,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}
	if err := saveToken(filepath.Join(dir, "anilist_token.json"), tok); err != nil {
		t.Fatalf("saveToken: %v", err)
	}
}

func withGlobalConfig(t *testing.T, config *CurdConfig) {
	t.Helper()
	previous := GetGlobalConfig()
	SetGlobalConfig(config)
	t.Cleanup(func() { SetGlobalConfig(previous) })
}

func TestAnilistTokenForAPIPrefersStoredToken(t *testing.T) {
	dir := t.TempDir()
	writeAnilistTokenFileForTest(t, dir, "stored-anilist-token")
	withGlobalConfig(t, &CurdConfig{StoragePath: dir})

	// The tracker token (user.Token) may be a MyAnimeList token; the AniList
	// call must use the real AniList token from disk instead.
	if got := anilistTokenForAPI(nil, "mal-jwt-token"); got != "stored-anilist-token" {
		t.Fatalf("anilistTokenForAPI = %q, want stored-anilist-token", got)
	}
}

func TestAnilistTokenForAPIFallsBackToCurrentWhenNoStoredToken(t *testing.T) {
	dir := t.TempDir()
	withGlobalConfig(t, &CurdConfig{StoragePath: dir})

	if got := anilistTokenForAPI(nil, "some-token"); got != "some-token" {
		t.Fatalf("anilistTokenForAPI = %q, want some-token", got)
	}
}

func TestAnilistTokenForAPIFallsBackWhenNoGlobalConfig(t *testing.T) {
	withGlobalConfig(t, nil)

	if got := anilistTokenForAPI(nil, "some-token"); got != "some-token" {
		t.Fatalf("anilistTokenForAPI = %q, want some-token", got)
	}
}

func TestTryRenewUsesStoredTokenInsteadOfBrowserReauth(t *testing.T) {
	dir := t.TempDir()
	writeAnilistTokenFileForTest(t, dir, "stored-anilist-token")
	withGlobalConfig(t, &CurdConfig{StoragePath: dir})

	// A MyAnimeList tracker token must not trigger a browser re-auth: the real
	// AniList token is simply swapped in and the request retried.
	newToken, renewed, err := tryRenewAniListTokenForAPI("mal-jwt-token", "invalid token")
	if err != nil {
		t.Fatalf("tryRenewAniListTokenForAPI error: %v", err)
	}
	if !renewed {
		t.Fatal("expected renewed=true")
	}
	if newToken != "stored-anilist-token" {
		t.Fatalf("renewed token = %q, want stored-anilist-token", newToken)
	}
}

func TestTryRenewEmptyTokenUsesStoredToken(t *testing.T) {
	dir := t.TempDir()
	writeAnilistTokenFileForTest(t, dir, "stored-anilist-token")
	withGlobalConfig(t, &CurdConfig{StoragePath: dir})

	newToken, renewed, err := tryRenewAniListTokenForAPI("", "invalid token")
	if err != nil {
		t.Fatalf("tryRenewAniListTokenForAPI error: %v", err)
	}
	if !renewed {
		t.Fatal("expected renewed=true")
	}
	if newToken != "stored-anilist-token" {
		t.Fatalf("renewed token = %q, want stored-anilist-token", newToken)
	}
}

func TestTryRenewNoGlobalConfigFallsThrough(t *testing.T) {
	withGlobalConfig(t, nil)

	// Without global config there is no stored token to compare against, so the
	// function must not panic and simply report no renewal.
	newToken, renewed, err := tryRenewAniListTokenForAPI("token", "invalid token")
	if err != nil {
		t.Fatalf("tryRenewAniListTokenForAPI error: %v", err)
	}
	if renewed {
		t.Fatal("expected renewed=false when no global config")
	}
	if newToken != "token" {
		t.Fatalf("renewed token = %q, want original token", newToken)
	}
}
