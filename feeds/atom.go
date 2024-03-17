package feeds

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/saaste/bookmark-manager/bookmarks"
)

func ToAtom(feedInfo FeedInfo, bookmarks []*bookmarks.Bookmark) string {
	pubDate := time.Now()
	if len(bookmarks) > 0 {
		pubDate = bookmarks[0].Created
	}

	output := make([]string, 0)

	output = append(output, `<?xml version="1.0" encoding="utf-8"?>`)

	// Feed
	output = append(output, `<feed xmlns="http://www.w3.org/2005/Atom">`)
	output = append(output, fmt.Sprintf("\t<title>%s</title>", html.EscapeString(feedInfo.SiteName)))
	output = append(output, fmt.Sprintf("\t<subtitle>%s</subtitle>", html.EscapeString(feedInfo.Description)))
	output = append(output, fmt.Sprintf(`%s<link href="%s" />`, "\t", feedInfo.BaseURL))
	output = append(output, fmt.Sprintf(`%s<link href="%s" rel="self" />`, "\t", feedInfo.CurrentURL))
	output = append(output, fmt.Sprintf("\t<updated>%s</updated>", pubDate.Format(time.RFC3339)))

	output = append(output, "\t<author>")
	output = append(output, fmt.Sprintf("\t\t<name>%s</name>", html.EscapeString(feedInfo.AuthorName)))
	if feedInfo.AuthorEmail != "" {
		output = append(output, fmt.Sprintf("\t\t<email>%s</email>", html.EscapeString(feedInfo.AuthorEmail)))
	}
	output = append(output, "\t</author>")
	output = append(output, fmt.Sprintf("\t<id>%s</id>", feedInfo.BaseURL))
	output = append(output, "\t<generator>Bookmark Manager (https://github.com/saaste/bookmark-manager)</generator>")

	// Entries
	for _, bm := range bookmarks {
		output = append(output, "\t<entry>")
		output = append(output, fmt.Sprintf("\t\t<id>%s:%s</id>", feedInfo.BaseURL, bm.URL))
		output = append(output, fmt.Sprintf("\t\t<title>%s</title>", html.EscapeString(bm.Title)))
		output = append(output, fmt.Sprintf("\t\t<updated>%s</updated>", bm.Created.Format(time.RFC3339)))
		output = append(output, fmt.Sprintf("\t\t<content>%s</content>", html.EscapeString(bm.Description)))
		output = append(output, "\t</entry>")
	}

	output = append(output, "</feed>")

	return strings.Join(output, "\n")
}
