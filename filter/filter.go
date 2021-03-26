package filter

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/*

Specifier -- ID

BOOK FILTERS {
	Author -- String check
  Genre  -- String check
	Date   -- range of dates
}

BOOK LIMITERS {
	Block any columns
}

COLLECTION FILTERS {
	Author -- String check (return all collections that possess author)
	Genre  -- String check (...)
	Date   -- range (...)
}

COLLECTION LIMITERS {
	block all books
}

*/

type Filter struct {
	Specifier string
	Filters   string
	Limiters  string
}

/*
filters are of the format: .../id#filter=...#limiter=...
ex:
localhost:8080/collections/4?filter=author+eq+firstname+lastname&edition+eq+1#limiter=title,genre

this will return the titles and genres of all books from the collection with id 4,
whose author is firstname lastname and is of edition number 1

Filters and limiters do not have an order, but id must come first.


id									-> []
id?filter           -> [id, filter]
id?filter#limiter   -> [id, filter, limiter]
id#limiter          -> [id, limiter]
id#limiter?filter   -> [id, limiter, filter]

?filter             -> ["", filter]
?filter#limiter     -> ["", filter, limiter]
#limiter            -> ["", limiter]
#limiter?filter     -> ["", limiter, filter]

*/

func GetFilters(path string) (*Filter, error) {
	filter := &Filter{Limiters: "*"}
	re := regexp.MustCompile(`(#|\?)`)
	split := re.Split(path, -1)
	// No filters/limiters, only id
	if len(split) == 0 {
		err := filter.addSpecifier(path)
		if err != nil {
			return nil, err
		}
		return filter, nil
	}
	// Maximum 3 fields supported
	if len(split) > 3 {
		return nil, errors.New(fmt.Sprintf("invalid filter syntax"))
	}

	err := filter.addSpecifier(split[0])
	if err != nil {
		return nil, err
	}

	for _, p := range split[1:] {
		err := filter.addColumnLimiters(p)
		if err != nil {
			return nil, err
		}
		err = filter.addColumnFilters(p)
		if err != nil {
			return nil, err
		}
	}

	return filter, nil
}

func (f *Filter) addSpecifier(path string) error {
	if len(path) == 0 {
		return nil
	}
	_, err := strconv.Atoi(path)
	if err != nil {
		return err
	}
	f.Specifier = path
	return nil
}

func (f *Filter) addColumnFilters(path string) error {
	split := strings.Split(path, "filter=")
	if len(split) < 2 {
		return nil
	}
	if len(split) > 2 {
		return errors.New("invalid filter syntax")
	}
	filters := strings.Split(split[1], "&")
	for _, filter := range filters {
		if len(filter) == 0 {
			return errors.New("invalid filter syntax")
		}

	}
	return nil
}

func (f *Filter) addColumnLimiters(path string) error {
	split := strings.Split(path, "limiter=")
	if len(split) < 2 {
		return nil
	}
	if len(split) > 2 {
		return errors.New("invalid limiter syntax")
	}
	limiters := strings.Split(split[1], ",")
	for _, l := range limiters {
		if !isWord(l) {
			return errors.New("invalid limiter syntax")
		}
	}
	f.Limiters = split[1]
	return nil
}

func isWord(name string) bool {
	if len(name) == 0 {
		return false
	}
	for _, c := range name {
		if ((c < 'a' || c > 'z') && (c < 'A' || c > 'Z')) || c != '_' {
			return false
		}
	}
	return true
}

func validateFilter(filter string) (string, error) {
	//author+eq+firstname+lastname
	//&edition+eq+1
	split := strings.Split(filter, "+")

	for _, s := range split {
		if len(s) == 0 {
			return "", errors.New("invalid filter syntax")
		}
	}

	return "", nil
}
