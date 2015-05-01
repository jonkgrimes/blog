package models

import (
	"testing"
	"time"
)

var publishedAtCases = []struct {
	post     Post
	expected string
}{
	{
		Post{Title: "A Title of a Post", PublishedAt: time.Date(2015, time.May, 3, 8, 0, 0, 0, time.UTC)},
		"May 3, 2015 at 8:00am",
	},
	{
		Post{Title: "A Title of a Post", PublishedAt: time.Date(2015, time.May, 3, 8, 0, 0, 0, time.UTC)},
		"May 3, 2015 at 8:00am",
	},
}

func TestPrettyPublishedAt(t *testing.T) {
	t.Log("TestPublishedAt")
	for _, testCase := range publishedAtCases {
		post_view := PostView{Post: testCase.post}
		actual := post_view.PrettyPublishedAt()
		expected := testCase.expected
		if expected != actual {
			t.Logf("Expected \"%s\" but got \"%s\"", expected, actual)
			t.Fail()
		}
	}
}
