package internal

import (
	"os"
	"testing"
)

// Exercises determineProviderTotalEpisodes against the real provider stack for
// "One Piece" using the app's own Go HTTP client (no browser), to confirm the
// total is resolved even when individual providers are down or rate-limited.
func TestLiveDetermineProviderTotalEpisodes(t *testing.T) {
	if os.Getenv("CURD_LIVE_DETERMINE_TOTAL") == "" {
		t.Skip("set CURD_LIVE_DETERMINE_TOTAL=1 to run live")
	}
	cfg := &CurdConfig{
		StoragePath: os.ExpandEnv("$HOME/.local/share/curd"),
		SubOrDub:    "sub",
		Provider:    "stacked",
		TrackingRemote: "myanimelist",
	}
	withGlobalConfig(t, cfg)

	anime := &Anime{
		AnilistId:    21,
		MalId:        21,
		TotalEpisodes: 1100,
		Title:         AnimeTitle{English: "One Piece", Romaji: "One Piece", Japanese: "ONE PIECE"},
		ProviderName:  "",
		ProviderId:    "",
		IsAiring:      true,
	}

	for _, pn := range configuredProviderNames(cfg) {
		p, err := ProviderByName(pn)
		if err != nil {
			t.Logf("  [%s] not configured: %v", pn, err)
			continue
		}
		opts, derr := p.SearchAnime("One Piece", "sub")
		if derr != nil {
			t.Logf("  [%s] search ERROR: %v", pn, derr)
			continue
		}
		if len(opts) == 0 {
			t.Logf("  [%s] search returned no results", pn)
			continue
		}
		best, _ := selectBestProviderSearchResult(opts, anime, "One Piece")
		total, terr := getProviderTotalEpisodes(p, best.Key, "sub")
		t.Logf("  [%s] matched Key=%q total=%d err=%v", pn, best.Key, total, terr)
	}

	total, err := determineProviderTotalEpisodes(cfg, "One Piece", anime, "sub")
	t.Logf("determineProviderTotalEpisodes -> total=%d err=%v", total, err)
	if err != nil {
		t.Fatalf("could not determine total episodes from providers: %v", err)
	}
	if total <= 0 {
		t.Fatalf("got non-positive total %d", total)
	}
	if total < 1000 {
		t.Fatalf("expected a realistic One Piece total >1000, got %d (wrong match or broken ep list)", total)
	}
	t.Logf("OK: One Piece total episodes = %d", total)
}