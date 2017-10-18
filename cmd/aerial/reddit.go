package main

import "github.com/jzelinskie/geddit"

func redditSearch(subreddit string) ([]*geddit.Submission, error) {
	listOptions := geddit.ListingOptions{
		Limit: 50,
	}
	reddit := geddit.NewSession("discordbot")
	results, err := reddit.SubredditSubmissions(subreddit, geddit.NewSubmissions, listOptions)
	if err != nil {
		return nil, err
	}
	return results, nil
}
