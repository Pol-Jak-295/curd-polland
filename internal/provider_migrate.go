package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const stackedProviderConfigValue = "stacked"

func isStackedProviderConfig(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "stacked", "stack", "auto", "all":
		return true
	default:
		return false
	}
}

func isFactoryDefaultProvider(raw string) bool {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return true
	}
	if isStackedProviderConfig(raw) {
		return true
	}

	names, declined := parseProviderConfig(raw)
	if declined {
		return false
	}
	switch len(names) {
	case 0:
		return true
	case 1:
		switch names[0] {
		case "senshi", "allanime":
			return true
		}
	}
	return false
}

func providerListsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func migrateProviderConfig(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if isFactoryDefaultProvider(raw) {
		if raw == stackedProviderConfigValue {
			return raw, false
		}
		return stackedProviderConfigValue, true
	}

	if isStackedProviderConfig(raw) {
		if raw != stackedProviderConfigValue {
			return stackedProviderConfigValue, true
		}
		return raw, false
	}

	names, declined := parseProviderConfig(raw)
	if len(names) > 1 {
		return stackedProviderConfigValue, true
	}

	canonical := canonicalProviderConfigValue(raw)
	if canonical != raw {
		return canonical, true
	}
	_ = declined
	return raw, false
}

func readStoredCurdVersion(storagePath string) string {
	path := storageVersionFilePath(storagePath)
	raw, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(raw))
}

func writeStoredCurdVersion(storagePath, version string) error {
	storagePath = strings.TrimSpace(storagePath)
	version = strings.TrimSpace(version)
	if storagePath == "" || version == "" {
		return nil
	}
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return err
	}
	return os.WriteFile(storageVersionFilePath(storagePath), []byte(version+"\n"), 0644)
}

// configOptionsIntroducedInVersion maps config keys to the first curd release that
// introduced them. Only these keys are eligible for automatic append on upgrade.
// Baseline options (Player, StoragePath, …) are NOT listed — they only appear via
// createDefaultConfig for brand-new installs, never re-appended into sparse configs.
//
// When adding a new option:
//  1. Add it to CurdConfig + defaultConfigMap()
//  2. Register it here under the release version that ships it
//  3. Bump VERSION.txt so MigrateOnVersionUpgrade runs for existing users
func configOptionsIntroducedInVersion() map[string]string {
	return map[string]string{
		// 2.0.3 — playback fallback timeout + vim selection motions
		"MpvPlaybackStartTimeout": "2.0.3",
		"VimKeys":                 "2.0.3",
		// 2.0.4 — idle background update checks
		"CheckUpdates": "2.0.4",
		// 2.0.5 — MPV episode playlist + alternate audio entries
		"MpvEpisodePlaylist": "2.0.5",
	}
}

// injectConfigOptionsSince appends defaults for options introduced in versions
// (fromVersion, toVersion] that are still missing from configMap.
// Example: from 2.0.2 → 2.0.3 injects only keys introduced in 2.0.3, not every
// historical default.
func injectConfigOptionsSince(configMap map[string]string, fromVersion, toVersion string) []string {
	if configMap == nil {
		return nil
	}
	defaults := defaultConfigMap()
	introduced := configOptionsIntroducedInVersion()

	keys := make([]string, 0, len(introduced))
	for key := range introduced {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	added := make([]string, 0)
	for _, key := range keys {
		since := introduced[key]
		// Only options born after the user's previous version and at/before this release.
		if !versionLess(fromVersion, since) {
			continue // introduced at or before fromVersion — user already "passed" that release
		}
		if !versionLessOrEqual(since, toVersion) {
			continue // not shipped yet in toVersion
		}
		if _, exists := configMap[key]; exists {
			continue
		}
		value, ok := defaults[key]
		if !ok {
			continue
		}
		configMap[key] = value
		added = append(added, key)
	}
	return added
}

// appendConfigKeys appends only the given keys to the config file so existing
// user options and ordering are left untouched.
func appendConfigKeys(configPath string, configMap map[string]string, keys []string) error {
	if strings.TrimSpace(configPath) == "" || len(keys) == 0 {
		return nil
	}
	file, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, key := range keys {
		value, ok := configMap[key]
		if !ok {
			continue
		}
		if _, err := fmt.Fprintf(file, "%s=%s\n", key, value); err != nil {
			return err
		}
	}
	return nil
}

// MigrateOnVersionUpgrade updates stored state and config when curd is upgraded.
// On version change it appends only config options registered as introduced in
// versions (storedVersion, appVersion], then runs provider migrations.
// Same-version launches do not rewrite the config file.
// Returns whether the config file was updated.
func MigrateOnVersionUpgrade(configPath string, config *CurdConfig, appVersion string) (bool, error) {
	if config == nil {
		return false, nil
	}

	appVersion = strings.TrimSpace(appVersion)
	if appVersion == "" {
		appVersion = CurdVersion()
	}

	storagePath := os.ExpandEnv(config.StoragePath)
	if storagePath == "" {
		storagePath = filepath.Join(os.ExpandEnv("$HOME"), ".local", "share", "curd")
	}

	storedVersion := readStoredCurdVersion(storagePath)
	configUpdated := false
	versionChanged := storedVersion != appVersion

	if versionChanged && strings.TrimSpace(configPath) != "" {
		configMap, err := LoadConfigFromFile(configPath)
		if err != nil {
			return false, err
		}

		// Respect AddMissingOptions=false as a hard opt-out of writing new keys.
		addMissing := true
		if config != nil {
			addMissing = config.AddMissingOptions
		}
		if val, exists := configMap["AddMissingOptions"]; exists {
			if parsed, parseErr := parseConfigBool(val); parseErr == nil {
				addMissing = parsed
			}
		}

		if addMissing {
			if added := injectConfigOptionsSince(configMap, storedVersion, appVersion); len(added) > 0 {
				if err := appendConfigKeys(configPath, configMap, added); err != nil {
					return false, fmt.Errorf("append new config options: %w", err)
				}
				configUpdated = true
				Log(fmt.Sprintf("Injected config options for upgrade %s → %s: %s",
					storedVersion, appVersion, strings.Join(added, ", ")))
			}
		}

		if nextProvider, changed := migrateProviderConfig(configMap["Provider"]); changed {
			configMap["Provider"] = nextProvider
			// Provider value already exists in the file — rewrite the full map once.
			if err := SaveConfigToFile(configPath, configMap); err != nil {
				return configUpdated, err
			}
			configUpdated = true
		}

		// Refresh in-memory config so new keys (e.g. VimKeys) apply immediately.
		next := PopulateConfig(configMap)
		normalizeTrackingConfig(&next)
		*config = next
	} else if versionChanged {
		// No config path, still run in-memory provider migration.
		if nextProvider, changed := migrateProviderConfig(config.Provider); changed {
			config.Provider = nextProvider
			configUpdated = true
		}
	}

	if err := writeStoredCurdVersion(storagePath, appVersion); err != nil {
		return configUpdated, fmt.Errorf("write curd version file: %w", err)
	}

	return configUpdated, nil
}

func parseConfigBool(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1", "yes", "y", "on":
		return true, nil
	case "false", "0", "no", "n", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid bool %q", value)
	}
}

func providerConfigDisplayLabel(raw string) string {
	if isStackedProviderConfig(raw) {
		names := defaultEnabledProviderStack()
		if len(names) == 0 {
			return "Default with fallback"
		}
		return fmt.Sprintf("Default with fallback (%s)", strings.Join(names, " → "))
	}
	names, _ := parseProviderConfig(raw)
	if len(names) == 1 {
		return names[0]
	}
	if len(names) > 1 {
		return strings.Join(names, " → ")
	}
	return raw
}
