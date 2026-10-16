# wax

CLI and TUI for browsing your Discogs vinyl collection and marketplace.

## Installation

```bash
git clone https://github.com/phlawski/wax
cd wax
go build -o wax .
```

## Configuration

On first run, wax creates `~/.config/wax/config.yaml`:

```yaml
discogs:
  token: <your_discogs_personal_access_token>
  username: <your_discogs_username>
```

Get your personal access token at [discogs.com/settings/developers](https://www.discogs.com/settings/developers). Your username is the display name shown on your Discogs profile.

## Usage

### Interactive TUI

```bash
wax tui
```

| Key | Action |
|-----|--------|
| `↑/↓` `k/j` | Move cursor |
| `n/p` | Next / previous page |
| `Enter` | Open release detail |
| `/` | Search Discogs database |
| `Esc` | Go back |
| `q` | Quit |

From search results, `Enter` opens marketplace listings for the selected release.

### CLI commands

```bash
# List your collection
wax collection list
wax collection list --page 2 --per-page 25

# Search your collection
wax collection search "Angelic Upstarts"

# Search Discogs database
wax search "Anti-Nowhere League"
wax search "Discharge" --type release

# Marketplace listings for a release ID
wax marketplace listings 12345
wax marketplace search "4 Skins"
```

## License

MIT
