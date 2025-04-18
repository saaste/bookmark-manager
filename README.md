# Bookmark Manager
A minimalistic self-hosted bookmark manager written in [Go](https://go.dev/).

The app aims for simplicity. That is why it is not filled with features. Bookmarks are
stored in SQLite database, so heavy database engines are not needed and makes database
backups easy to do.

Bookmarks can be marked as private. Private bookmarks are visible only for authenticated
users. Each bookmark can also have tags. RSS feed is available for recent bookmarks and
for each tag. The app can also do [scheduled bookmark checks](#scheduled-bookmark-check).

Bookmark Manager aims to be as accessible as possible. It should have a pretty a quite good
screen reader support, but if have an idea how to make it even better, feel free to file an issue!

## Requirements
- [Docker](https://www.docker.com)

## How to run
1. Copy the following example files:
   - bookmarks.db.example ➔ bookmarks.db
   - docker-compose.yml.example ➔ docker-compose.yml
   - config.yml.example ➔ config.yml
   - robots.txt.example ➔ robots.txt
2. Configure settings in `config.yml`. Each setting is described in the config file.
3. Configure `robots.txt` if you want to limit how search engines and bots scrape your site

### Running with Docker (recommended)
```
docker-compose up [-d]
```
By default, the forwarded port is 8000. You can change this in the `docker-compose.yaml`.

### Running with Go
If you know what you are doing, you can also run the application directly from the source code or build your own binaries.
```
go run .
```
Be default, the app is listening to port 8000. You can change this in the `config.yml`.

## Scheduled bookmark check
The app can check for broken bookmarks, but this feature is disabled by default.
To enable the check, set the `check_interval` setting in `config.yml` to `1` or more.
If you want the check to run when the app starts, set the `check_on_app_start` setting to `true`.

Broken bookmarks are indicated by an exclamation point icon. You will also see a *Broken Bookmarks* link in the top navigation bar, which will take you to a view listing all broken bookmarks. These are only visible if you are logged in.

Bookmark Manager can send a notification when it detects broken bookmarks. Currently, only
[Gotify](https://gotify.net/) notifications are supported. To enable notifications, set
the `gotify_enabled` setting to `true` and set `gotify_url` and `gotify_token` settings to match your
environment.

Sometimes the bookmark check fails because the server classifies the check as bot traffic and displays a [CAPTCHA](https://en.wikipedia.org/wiki/CAPTCHA). It may also block the request completely, which is known to happen with certain [Cloudflare](https://www.cloudflare.com) configuration.

If some URLs work fine in a browser but fail the bookmark check, you can simply disable the check for the problematic URLs by checking the `Ignore URL checking` checkbox in the bookmark edit view.

## Customizing the UI
Bookmark Manager supports themes. If you want to create your own theme, read the separate [Theme Documentation](/docs/THEMES.md).

## Screenshots
UI for guests

[![Screenshot showing the dark UI for guest users](docs/sshot-dark-guest-tn.jpg "Dark UI for guests")](docs/sshot-dark-guest.jpg)
[![Screenshot showing the light UI for guest users](docs/sshot-light-guest-tn.jpg "Light UI for guests")](docs/sshot-light-guest.jpg)

UI for authenticated users

[![Screenshot showing the dark UI for authenticated users](docs/sshot-dark-auth-tn.jpg "Dark UI for authenticated users")](docs/sshot-dark-auth.jpg)
[![Screenshot showing the light UI for authenticated users](docs/sshot-light-auth-tn.jpg "Light UI for authenticated users")](docs/sshot-light-auth.jpg)

Admin UI

[![Screenshot showing the dark admin UI](docs/sshot-dark-admin-tn.jpg "Dark Admin UI")](docs/sshot-dark-admin.jpg)
[![Screenshot showing the light admin UI](docs/sshot-light-admin-tn.jpg "Light Admin UI")](docs/sshot-light-admin.jpg)