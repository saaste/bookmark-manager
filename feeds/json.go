package feeds

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/saaste/bookmark-manager/bookmarks"
)

type jsonFeed struct {
	Version     string       `json:"version"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	HomePageUrl string       `json:"home_page_url"`
	FeedUrl     string       `json:"feed_url,omitempty"`
	Authors     []jsonAuthor `json:"authors,omitempty"`
	Items       []jsonItem   `json:"items"`
}

type jsonAuthor struct {
	Name string `json:"name"`
}

type jsonItem struct {
	Id            string   `json:"id"`
	Title         string   `json:"title"`
	ContentText   string   `json:"content_text"`
	Url           string   `json:"url"`
	DatePublished string   `json:"date_published"`
	Tags          []string `json:"tags,omitempty"`
}

func ToJSON(feedInfo FeedInfo, bookmarks []*bookmarks.Bookmark) (string, error) {

	feed := jsonFeed{
		Version:     "https://jsonfeed.org/version/1.1",
		Title:       feedInfo.SiteName,
		Description: feedInfo.Description,
		HomePageUrl: feedInfo.BaseURL,
		FeedUrl:     feedInfo.CurrentURL,
	}

	if feedInfo.AuthorName != "" {
		feed.Authors = append(feed.Authors, jsonAuthor{Name: feedInfo.AuthorName})
	}

	for _, bm := range bookmarks {
		item := jsonItem{
			Id:            fmt.Sprintf("%s:%s", feedInfo.BaseURL, bm.URL),
			Title:         bm.Title,
			ContentText:   bm.Description,
			Url:           bm.URL,
			DatePublished: bm.Created.Format(time.RFC3339),
			Tags:          bm.Tags,
		}
		feed.Items = append(feed.Items, item)
	}

	byt, err := json.MarshalIndent(feed, "", "  ")
	return string(byt), err

}
