# Themes
You can customize the interface by creating a new theme. Themes are located in the `templates` directory. Each theme must have its own subdirectory. By default, the app uses the [default](/templates/default/) theme. The theme can be changed in the `config.yml` file by changing the `theme` configuration value.

If you are creating a new theme, it is probably easiest to use the `default` template as a base. Any static assets required by the theme must be in the `assets` subdirectory.


## Available template variables

| Variable          | Type          | Description
|------------------ | ------------- | -------------------------------------------------------
| .SiteName         | string        | Site name, defined in the config
| .Description      | string        | Site description, defined in the config
| .BaseURL          | string        | Site base URL, defined in the config
| .Title            | string        | Title of the current page
| .CurrentURL       | string        | URL requested to access the current view
| .IsAuthenticated  | bool          | Boolean indicating if user is authenticated
| .Bookmarks        | [[]Bookmark](#bookmark)    | List of bookmarks visible in the view
| .Tags             | []string      | List of all tags
| .TextFilter       | string        | Current search term
| .Pages            | []Page        | List of available pages for paginated content
| .BrokenBookmarksExist  | bool     | Boolean indicating if broken bookmarks exist

## Types

### Bookmark
| Field             | Type          | Description
|------------------ | ------------- | -------------------------------------------------------
| .ID               | int64         | Bookmark ID
| .URL              | string        | Bookmark URL
| .Title            | string        | Bookmark title
| .Description      | string        | Bookmark description
| .IsPrivate        | bool          | Is bookmark private
| .Created          | time.Time     | Bookmark creation datetime
| .Tags             | []string      | Bookmark tags
| .IsWorking        | bool          | Is bookmark working

## Components
Components render elements, that are necessary for the app to work. You should use these instead of building your own. Otherwise some features may not work.


| Components | Usage
|------------| --------
| bookmark_add_form | Form for adding new bookmarks
| bookmark_edit_form | Form for editing bookmarks
| bookmark_delete_form | Form for deleting bookmarks
| login_form | Form for logging in
| headers | Required headers used in `<head>` section

How to use:
```
{{ template "feed-links" . }}
```

## Functions

### feedUrl (feedType string)
Returns a feed URL. Supported `feedType` values:
- rss
- atom
- json

How to use:
```html
<link rel="alternate" type="application/rss+xml" title="RSS Feed" href='{{ feedUrl "rss" }}'>
```

### anchorUrl (id string)
Returns an URL pointing to a specific element on the page

How to use:
```html
<a href='{{ anchorUrl "q" }}' id="top">Go to Search</a>
...
<input type="text" name="q" id="q" aria-label="Search by keyword">
```

### paginationUrl (pageNumber int)
return an URL pointing to a specific page of the paginated content

How to use:
```html
<nav id="pagination" aria-label="Pagination Navigation">
    {{ range .Pages }}
        {{ if .IsActive }}
            {{ .Number }}
        {{ else }}
            <a href="{{ paginationUrl .Number }}" aria-label="Go to page {{ .Number }}">
                {{ .Number }}
            </a>
        {{ end }}
    {{ end }}
</nav>
```