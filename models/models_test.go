package models

import (
	"testing"
	"time"
)

var createSlugCases = []struct {
	post     Post
	expected string
}{
	{
		Post{Title: "A Title For A Post"},
		"a-title-for-a-post",
	},
	{
		Post{Title: "Some punctuation's in this one"},
		"some-punctuation-s-in-this-one",
	},
	{
		Post{Title: "  This will have   uneven   spaces   "},
		"this-will-have-uneven-spaces",
	},
	{
		Post{Title: "This Replaces & with and"},
		"this-replaces-and-with-and",
	},
}

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

func TestCreateSlug(t *testing.T) {
	t.Log("TestCreateSlug")
	for _, testCase := range createSlugCases {
		post := testCase.post
		actual := post.CreateSlug()
		expected := testCase.expected
		if expected != actual {
			t.Logf("Expected \"%s\" but got \"%s\"", expected, actual)
			t.Fail()
		}
	}
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
