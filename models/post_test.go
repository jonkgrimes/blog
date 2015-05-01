package models

import (
	"testing"
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
		"some-punctuations-in-this-one",
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

func TestCreateSlug(t *testing.T) {
	t.Log("TestCreateSlug")
	for _, testCase := range createSlugCases {
		post := testCase.post
		post.CreateSlug()
		expected := testCase.expected
		actual := post.Slug

		if expected != actual {
			t.Logf("Expected \"%s\" but got \"%s\"", expected, actual)
			t.Fail()
		}
	}
}
