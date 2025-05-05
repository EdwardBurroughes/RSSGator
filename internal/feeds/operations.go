package feeds

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"net/http"

	"github.com/EdBurroughes/rss-blog-aggregator/internal/database"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to generate new request: %v", err)
	}

	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("bad response: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return &RSSFeed{}, fmt.Errorf("non-200 status code: %d", res.StatusCode)
	}

	var feed RSSFeed
	if err := xml.NewDecoder(res.Body).Decode(&feed); err != nil {
		return &RSSFeed{}, err
	}
	feed.unescapeFields()
	return &feed, nil
}

func insertPosts(ctx context.Context, db *database.Queries, rssItems []RSSItem, feed_id int) error {
	// this is all a bit smelly I should really be updating records if they already exist
	// however currently just ignoring
	for _, item := range rssItems {
		_, err := db.CreatePost(ctx, database.CreatePostParams{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: item.Title, Valid: true},
			Url:         sql.NullString{String: item.Link, Valid: true},
			PublishedAt: sql.NullString{String: item.PubDate, Valid: true},
			FeedID:      sql.NullInt32{Int32: int32(feed_id), Valid: true},
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			fmt.Printf("Couldn't create posts: %v\n, %t", err, err)
			continue
		}

	}
	return nil
}

func ScrapeFeeds(ctx context.Context, db *database.Queries) error {
	feedToFetch, err := db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Fetching posts from: %s\n", feedToFetch.Name)
	rss, err := fetchFeed(ctx, feedToFetch.Url)
	if err != nil {
		return err
	}
	insertPosts(ctx, db, rss.Channel.Item, int(feedToFetch.ID))
	db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID:            feedToFetch.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	return nil
}
