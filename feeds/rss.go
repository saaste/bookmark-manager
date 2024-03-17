package feeds

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/saaste/bookmark-manager/bookmarks"
)

func ToRSS(feedInfo FeedInfo, bookmarks []*bookmarks.Bookmark) string {
	pubDate := time.Now()
	if len(bookmarks) > 0 {
		pubDate = bookmarks[0].Created
	}

	output := make([]string, 0)

	output = append(output, `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">`)
	output = append(output, "<channel>")

	// Channel
	output = append(output, fmt.Sprintf("\t<title>%s</title>", html.EscapeString(feedInfo.SiteName)))
	output = append(output, fmt.Sprintf("\t<link>%s</link>", feedInfo.BaseURL))
	output = append(output, fmt.Sprintf("\t<description>%s</description>", html.EscapeString(feedInfo.Description)))
	output = append(output, fmt.Sprintf("\t<pubDate>%s</pubDate>", pubDate.Format(time.RFC1123Z)))
	output = append(output, fmt.Sprintf("\t<lastBuildDate>%s</lastBuildDate>", pubDate.Format(time.RFC1123Z)))
	output = append(output, "\t<generator>Bookmark Manager (https://github.com/saaste/bookmark-manager)</generator>")
	output = append(output, fmt.Sprintf(`%s<atom:link href="%s" rel="self" type="application/rss+xml"></atom:link>`, "\t", feedInfo.CurrentURL))

	// Items
	for _, bm := range bookmarks {
		output = append(output, "\t<item>")

		output = append(output, fmt.Sprintf("\t\t<title>%s</title>", html.EscapeString(bm.Title)))
		output = append(output, fmt.Sprintf("\t\t<link>%s</link>", bm.URL))

		if bm.Description != "" {
			output = append(output, fmt.Sprintf("\t\t<description>%s</description>", html.EscapeString(bm.Title)))
		}

		output = append(output, fmt.Sprintf("\t\t<pubDate>%s</pubDate>", bm.Created.Format(time.RFC1123Z)))
		output = append(output, fmt.Sprintf("\t\t<guid>%s:%s</guid>", feedInfo.BaseURL, bm.URL))
		output = append(output, "\t</item>")
	}

	output = append(output, "</channel>")
	output = append(output, "</rss>")

	return strings.Join(output, "\n")
}
