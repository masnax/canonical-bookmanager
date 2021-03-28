package parser

import (
	"fmt"
	"net/url"
	"testing"
)

func TestValidURL(t *testing.T) {
	basic, _ := url.Parse("http://www.google.com")
	basicWithPath, _ := url.Parse("http://www.google.com/")
	basicWithSinglePath, _ := url.Parse("http://www.google.com/a")
	longPath, _ := url.Parse("http://www.google.com/aabc/def/ge")
	longPathWithEnd, _ := url.Parse("http://www.google.com/aabc/def/ge/")

	testCases := []struct {
		desc string
		url  *url.URL
		out  []string
	}{
		{
			desc: "basic url",
			url:  basic,
			out:  nil,
		},
		{
			desc: "empty end path",
			url:  basicWithPath,
			out:  nil,
		},
		{
			desc: "single path",
			url:  basicWithSinglePath,
			out:  nil,
		},
		{
			desc: "long title",
			url:  longPath,
			out:  []string{"", "aabc", "def", "ge"},
		},
		{
			desc: "long path with end",
			url:  longPathWithEnd,
			out:  []string{"", "aabc", "def", "ge", ""},
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.desc), func(t *testing.T) {
			res := URLParser(tc.url)
			if tc.out == nil && res != nil {
				t.Fatalf("expected nil response, got [%v]", res)
			}
			if res == nil && tc.out != nil {
				t.Fatalf("expected [%v] as response, got nil", tc.out)
			}
			if len(res) != len(tc.out) {
				t.Fatalf("expected result [%v], got [%v]", tc.out, res)
			}
			for i, p := range tc.out {
				if res[i] != p {
					t.Fatalf("expected result [%v], got [%v]", tc.out, res)
				}
			}

		})
	}
}
