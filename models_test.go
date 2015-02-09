package main

import "testing"

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
