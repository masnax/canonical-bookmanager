package filter

import (
	"fmt"
	"testing"

	"github.com/masnax/canonical-bookmanager/book"
)

func TestValidFilter(t *testing.T) {
	testCases := []struct {
		desc      string
		formValue string
	}{
		{
			desc:      "simple title",
			formValue: "title eq a",
		},
		{
			desc:      "uppercase title",
			formValue: "title eq A",
		},
		{
			desc:      "empty title",
			formValue: "title eq ",
		},
		{
			desc:      "long title",
			formValue: "title eq abc def g hgi",
		},
		{
			desc:      "valid date",
			formValue: "published eq 2020-04-01",
		},
		{
			desc:      "valid edition",
			formValue: "edition gt 1",
		},
		{
			desc:      "non-present key",
			formValue: "thing gt 1",
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.desc), func(t *testing.T) {
			_, err := FilterBooks(tc.formValue, book.Book{Published: "2020-01-01"})
			if err != nil {
				t.Fatalf("expected no error, got [%v]", err)
			}
		})
	}
}

func TestInvalidFilter(t *testing.T) {
	testCases := []struct {
		desc      string
		formValue string
	}{
		{
			desc:      "broken",
			formValue: "fladfdkf",
		},
		{
			desc:      "invalid operator",
			formValue: "title stuff A",
		},
		{
			desc:      "invalid edition type",
			formValue: "edition eq abc",
		},
		{
			desc:      "string gt",
			formValue: "title gt 4",
		},
		{
			desc:      "string lt",
			formValue: "title lt 4",
		},
		{
			desc:      "string le",
			formValue: "title le 4",
		},
		{
			desc:      "string ge",
			formValue: "title ge 4",
		},
		{
			desc:      "invalid date",
			formValue: "published eq 2",
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.desc), func(t *testing.T) {
			_, err := FilterBooks(tc.formValue, book.Book{Published: "2020-01-01"})
			if err == nil {
				t.Fatalf("expected an error, got none")
			}
		})
	}
}
