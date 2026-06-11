# Curd

A cli application to stream anime with [Anilist](https://anilist.co/) integration and Discord RPC written in golang.
Works on Linux, MacOS and Windows.

## Disclaimer

This fork does not add functionality, only graphical improvements to fit my own preferences. For any feature requests, errors, donations and support please turn to the upstream repository at [Wraient/curd](https://github.com/Wraient/curd).

## Join the original discord server

https://discord.gg/rrpBfu2gHq

## Join the Matrix server

https://matrix.to/#/#curd:matrix.org

## Demo Video
Normal mode:


https://github.com/user-attachments/assets/376e7580-b1af-40ee-82c3-154191f75b79

Rofi with Image preview - from upstream (wrong rasi configs)


https://github.com/user-attachments/assets/cbf799bc-9fdd-4402-ab61-b4e31f1e264d



## Features
- Multiple Content Providers (AllAnime and Animepahe) with ordered fallback and up to 1080p support
- Built-in headless browser to bypass Cloudflare/DDoS-Guard protections
- Stream anime online
- Track anime locally, on AniList, or on MyAnimeList
- Browser-based AniList and MyAnimeList login flows
- Skip anime Intro and Outro
- Skip Filler and Recap episodes
- Discord RPC about the anime
- Rofi support
- Image preview in rofi
- Local anime history to continue from where you left off last time
- Save mpv speed for next episode
- Configurable through config file


## Installing and Setup
> **Note**: `Curd` requires `mpv`, `rofi`, and `ueberzugpp` for Rofi support and image preview. These are included in the installation instructions below for each distribution.

### Linux
<details>
<summary>Arch Linux / Manjaro (AUR-based systems)</summary>

Using Yay:

```bash
yay -Sy curd
```

or using Paru:

```bash
paru -Sy curd
```

Or, to manually clone and install:

```bash
git clone https://aur.archlinux.org/curd.git
cd curd
makepkg -si
sudo pacman -S rofi ueberzugpp
```
</details>

<details>
<summary>Debian / Ubuntu (and derivatives)</summary>

```bash
sudo apt update
sudo apt install mpv curl rofi ueberzugpp

# For x86_64 systems:
curl -Lo curd https://github.com/Pol-Jak-295/curd-polland/releases/latest/download/curd-linux-x86_64

# For ARM64 systems:
curl -Lo curd https://github.com/Pol-Jak-295/curd-polland/releases/latest/download/curd-linux-arm64

chmod +x curd
sudo mv curd /usr/bin/
curd
```
</details>

<details>
<summary>Fedora Installation</summary>

```bash
sudo dnf update
sudo dnf install mpv curl rofi ueberzugpp

# For x86_64 systems:
curl -Lo curd https://github.com/Pol-Jak-295/curd-polland/releases/latest/download/curd-linux-x86_64

# For ARM64 systems:
curl -Lo curd https://github.com/Pol-Jak-295/curd-polland/releases/latest/download/curd-linux-arm64

chmod +x curd
sudo mv curd /usr/bin/
curd
```
</details>

<details>
<summary>openSUSE Installation</summary>

```bash
sudo zypper refresh
sudo zypper install mpv curl rofi ueberzugpp

# For x86_64 systems:
curl -Lo curd https://github.com/Pol-Jak-295/curd-polland/releases/latest/download/curd-linux-x86_64

# For ARM64 systems:
curl -Lo curd https://github.com/Pol-Jak-295/curd-polland/releases/latest/download/curd-linux-arm64

chmod +x curd
sudo mv curd /usr/bin/
curd
```
</details>

<details>
<summary>NixOS Installation</summary>

1. Add curd as a flake input, for example:
```nix
{
    inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
        curd = {
            url = "github:Pol-Jak-295/curd";
            inputs.nixpkgs.follows = "nixpkgs";
        };
    }
}
```
2. Install the package, for example:
```nix
{inputs, pkgs, ...}: {
  environment.systemPackages = [
    inputs.curd.packages.${pkgs.system}.default
  ];
}
```

</details>

<details>
<summary>Generic Installation</summary>

Choose the appropriate binary for your system:
```bash
# For Linux x86_64:
curl -Lo curd https://github.com/Pol-Jak-295/curd-polland/releases/latest/download/curd-linux-x86_64

# For Linux ARM64:
curl -Lo curd https://github.com/Pol-Jak-295/curd-polland/releases/latest/download/curd-linux-arm64

chmod +x curd
sudo mv curd /usr/bin/
curd
```
</details>

<details>
<summary>Uninstallation</summary>

```bash
sudo rm /usr/bin/curd
```

For AUR-based distributions:

```bash
yay -R curd-polland-git
OR
yay -S curd-polland-bin
```
</details>


<details>
<summary>Uninstallation</summary>

```bash
sudo rm /usr/local/bin/curd
```

</details>


## Data Storage

</details>

<details>
<summary>Linux/Unix</summary>
Stroage: (Token, Timestamps, debug.log, etc)

```bash
$USER/.local/share/curd
```

Config : 

```bash
$USER/.config/curd
```

</details>

## Usage

Run `curd` with the following options:

```bash
curd [options]
```

### Arguments would always take precedence over configuration

> **Note**:
> - To use rofi you need rofi and ueberzug installed.
> - Rofi .rasi files are at default `~/.local/share/curd/`
> - You can edit them as you like.
> - If there are no rasi files with specific names, they would be downloaded from this repo.


### Options

| Flag                      | Description                                                             | Default       |
|---------------------------|-------------------------------------------------------------------------|---------------|
| `-c`                      | Continue the last episode                                              | -             |
| `-change-token`           | Change your authentication token                                       | -             |
| `-dub`                    | Watch the dubbed version of the anime                                  | -             |
| `-sub`                    | Watch the subbed version of the anime                                  | -             |
| `-new`                    | Add a new anime to your list                                           | -             |
| `-e`                      | Edit the configuration file                                            | -             |
| `-skip-op`                | Automatically skip the opening section of each episode                 | `true`        |
| `-skip-ed`                | Automatically skip the ending section of each episode                  | `true`        |
| `-skip-filler`            | Automatically skip filler episodes                                     | `true`        |
| `-skip-recap`             | Automatically skip recap sections                                      | `true`        |
| `-discord-presence`       | Enable or disable Discord presence                                     | `true`        |
| `-image-preview`          | Show an image preview of the anime                                     | -             |
| `-no-image-preview`       | Disable image preview                                                  | -             |
| `-next-episode-prompt`    | Prompt for the next episode after completing one                       | -             |
| `-rofi`                   | Open anime selection in the rofi interface                             | -             |
| `-no-rofi`                | Disable rofi interface                                                 | -             |
| `-percentage-to-mark-complete` | Set the percentage watched to mark an episode as complete       | `85`          |
| `-player`                 | Specify the player to use for playback                                 | `"mpv"`       |
| `-save-mpv-speed`         | Save the current MPV speed setting for future sessions                 | `true`        |
| `-score-on-completion`    | Prompt to score the episode on completion                              | `true`        |
| `-storage-path`           | Path to the storage directory                                          | `"$HOME/.local/share/curd"` |
| `-subs-lang`              | Set the language for subtitles                                         | `"english"`   |
| `-u`                      | Update the script                                                      | -             |
| `-v`                      | Show curd version                                                      | -             |

### Examples

- **Continue the Last Episode**:
  ```bash
  curd -c
  ```

- **Add a New Anime**:
  ```bash
  curd -percentage-to-mark-complete=90
  ```

- **Play with Rofi and Image Preview**:
  ```bash
  curd -rofi -image-preview
  ```

## Configuration

All configurations are stored in a file you can edit with the `-e` option.

```bash
curd -e
```

Script is made in a way that you use it for one session of watching.

You can quit it anytime and the resume time would be saved in the history file

more settings can be found at config file.
config file is located at ```~/.config/curd/curd.conf```

On first start, curd asks which tracking mode you want to use:

- local
- anilist
- myanimelist
- anilist + myanimelist

Local history stays enabled in all modes. `anilist` means local history + AniList sync, `myanimelist` means local history + MyAnimeList sync, and `anilist + myanimelist` updates both platforms. In dual-sync mode, curd compares the latest remote update time for each anime and pushes the newest status/progress/category back to the older tracker so both sides converge automatically. Legacy installs are migrated to `anilist` automatically. If you switch trackers later, curd exposes a **Change Tracker** menu entry and can either merge both lists, replace MyAnimeList with AniList, or replace AniList with MyAnimeList.

MyAnimeList OAuth requires your own MAL application credentials. You can set them in the config file or with environment variables:

```bash
CURD_MAL_CLIENT_ID=your_client_id
CURD_MAL_CLIENT_SECRET=your_client_secret
```

If the browser reaches the localhost callback page but curd does not continue automatically, rerun the command and paste the full callback URL when prompted. Curd now keeps the pending MyAnimeList PKCE state so the login can be completed manually.

| **Option**               | **Type**   | **Valid Values**                           | **Description**                                                                                   |
|---------------------------|------------|-------------------------------------------|---------------------------------------------------------------------------------------------------|
| `DiscordPresence`         | Boolean    | `true`, `false`                           | Enables or disables Discord Rich Presence integration.                                            |
| `AnimeNameLanguage`       | Enum       | `english`, `romaji`                       | Sets the preferred language for anime names.                                                      |
| `MpvArgs`                 | List       | all mpv args eg ["--fullscreen=yes", "--mute=yes"]     | Add args to mpv player                                                               | 
| `AddMissingOptions`       | Boolean    | `true`, `false`                           | Automatically adds missing configuration options with default values to the config file.          |
| `AlternateScreen`         | Boolean    | `true`, `false`                           | Toggles the use of an alternate screen buffer for cleaner UI.                                     |
| `RofiSelection`           | Boolean    | `true`, `false`                           | Enables or disables anime selection via Rofi.                                                     |
| `PercentageToMarkComplete`| Integer    | `0` to `100`                              | Sets the percentage of an episode watched to consider it as completed.                            |
| `StoragePath`             | String     | Any valid path (Environment variables accepted)  | Specifies the directory where Curd stores its data.                                        |
| `SubOrDub`                | Enum       | `sub`, `dub`                              | Sets the preferred format for anime audio.                                                        |
| `NextEpisodePrompt`       | Boolean    | `true`, `false`                           | Prompts the user before automatically playing the next episode.                                   |
| `SubsLanguage`            | String     | `english` (redundant rn)                  | Sets the preferred subtitle language.                                                             |
| `ScoreOnCompletion`       | Boolean    | `true`, `false`                           | Automatically prompts the user to rate the anime upon completion.                                 |
| `SkipOp`                  | Boolean    | `true`, `false`                           | Automatically skips the opening of episodes when supported.                                       |
| `SkipEd`                  | Boolean    | `true`, `false`                           | Automatically skips the ending of episodes when supported.                                        |
| `SkipRecap`               | Boolean    | `true`, `false`                           | Skips recap sections in episodes when supported.                                                  |
| `ImagePreview`            | Boolean    | `true`, `false`                           | Enables or disables image previews during anime selection (only for rofi).                        |
| `Player`                  | String     | any mpv-compatible binary (e.g. `mpv`, `iina`) | Player binary used for playback. If not found, Curd falls back to `mpv`.                          |
| `SaveMpvSpeed`            | Boolean    | `true`, `false`                           | Retains the playback speed set in MPV for next episode.                                           |
| `SkipFiller`              | Boolean    | `true`, `false`                           | Skips filler episodes when supported.                                                             |
| `MenuOrder`               | String     | Comma-separated list                      | Controls which menu items appear and their order. Available options: `CURRENT`, `ALL`, `UNTRACKED`, `UPDATE`, `CONTINUE_LAST`, `PLANNING`, `COMPLETED`, `PAUSED`, `DROPPED`, `REWATCHING`, `TRACKER`, `PROVIDER`. Only listed items will be shown. Default: `CURRENT,ALL,UNTRACKED,UPDATE,CONTINUE_LAST,TRACKER,PROVIDER` |
| `Provider`                | List       | `["allanime"]`, `["animepahe"]`, `["allanime","animepahe"]`, `["allanime","no-animepahe"]` | Sets the ordered content-provider fallback list. Curd tries providers from left to right for the requested stream. If AllAnime has no search results or the requested episode stream is unavailable and Animepahe is not configured, Curd asks before enabling Animepahe and warns that DDoS-Guard verification may require a ~500 MB chromium download. Declining stores `no-animepahe` so Curd does not ask again. Legacy single values like `allanime` are migrated to list form. Default: `["allanime"]` |
| `TrackingLocal`           | Boolean    | `true`                                    | Legacy compatibility flag. Local playback history is always enabled.                              |
| `TrackingRemote`          | Enum       | `none`, `anilist`, `myanimelist`, `anilist+myanimelist` | Selects which remote tracker curd syncs with.                                           |
| `TrackingConfigured`      | Boolean    | `true`, `false`                           | Internal flag used to remember that the startup tracking prompt has already been completed.       |
| `MyAnimeListClientID`     | String     | MAL OAuth client ID                       | Client ID used for MyAnimeList browser login.                                                     |
| `MyAnimeListClientSecret` | String     | MAL OAuth client secret                   | Optional secret used for MyAnimeList browser login and token refresh.                             |
| `MyAnimeListImported`     | Boolean    | `true`, `false`                           | Tracks whether the one-time AniList-to-MyAnimeList import prompt has already been handled.        |


## Dependencies
- mpv - Video player (required fallback)
- rofi - Selection menu
- ueberzug - Display images in rofi
- chromium - Required for Animepahe (auto-downloaded by default, but Termux users must install manually via `pkg install chromium`)

## API Used
- [Anilist API](https://anilist.gitbook.io/anilist-apiv2-docs) - Update user data and download user data
- [MyAnimeList API](https://myanimelist.net/apiconfig/references/api/v2) - MyAnimeList OAuth and tracking sync
- [AniSkip API](https://api.aniskip.com/api-docs) - Get anime intro and outro timings
- [AllAnime Content](https://allanime.to/) - Fetch anime url
- [Animepahe Content](https://animepahe.pw/) - Alternative provider for 1080p streams
- [Jikan](https://jikan.moe/) - Get filler episode number

## Credits
- [curd](https://github.com/Wraient/curd) - the original, upstream version
- [ani-cli](https://github.com/pystardust/ani-cli) - Code for fetching anime url
- [jerry](https://github.com/justchokingaround/jerry) - For the inspiration
