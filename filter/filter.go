package filter

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/masnax/canonical-bookmanager/book"
)

/*
filter structure: url/path?filter=KEY+OP+VALUE
where KEY    is a field
			OP     is an operator
			VALUE  is a series of words delimited by '+' corresponding to an entry
*/

type Filter struct {
	Key string
	Op  string
	Val string
}

var validOps = []string{"eq", "ne", "lt", "gt", "le", "ge"}

func FilterBooks(form string, book book.Book) (bool, error) {
	filter, err := parseFilter(form)
	if err != nil {
		return false, err
	}
	r := reflect.ValueOf(book)
	for i := 0; i < r.NumField(); i++ {
		field := r.Type().Field(i)
		if strings.ToLower(field.Name) == strings.ToLower(filter.Key) {
			switch field.Type.Kind() {
			case reflect.Int:
				return handleOpInt(r.Field(i), filter)
			case reflect.String:
				return handleOpString(field.Name, r.Field(i), filter)
			default:
				return false, errors.New(fmt.Sprintf("unexpected field for book: %s", field.Name))
			}
		}
	}
	return true, nil
}

func handleOpString(name string, value reflect.Value, filter Filter) (bool, error) {
	valueStr := strings.ToLower(value.String())
	if name == "Published" {
		valueDate, err := time.Parse("2006-01-02", valueStr)
		if err != nil {
			return false, errors.New(fmt.Sprintf("invalid date for book, got %s", valueStr))
		}
		filterDate, err := time.Parse("2006-01-02", filter.Val)
		if err != nil {
			return false, errors.New(fmt.Sprintf("expected date of form Y-M-D, got %s", filter.Val))
		}

		switch filter.Op {
		case "eq":
			return !valueDate.Before(filterDate) && !valueDate.After(filterDate), nil
		case "ne":
			return valueDate.Before(filterDate) || valueDate.After(filterDate), nil
		case "lt":
			return valueDate.Before(filterDate), nil
		case "gt":
			return valueDate.After(filterDate), nil
		case "le":
			return !valueDate.After(filterDate), nil
		case "ge":
			return !valueDate.Before(filterDate), nil
		}
	} else {
		filterVal := strings.ToLower(filter.Val)
		switch filter.Op {
		case "eq":
			return valueStr == filterVal, nil
		case "ne":
			return valueStr != filterVal, nil
		}
	}
	return false, errors.New("invalid filter")
}

func handleOpInt(value reflect.Value, filter Filter) (bool, error) {
	filterVal, err := strconv.Atoi(filter.Val)
	if err != nil {
		return false, errors.New("expected integer value in form")
	}
	valueInt := value.Int()
	switch filter.Op {
	case "eq":
		return valueInt == int64(filterVal), nil
	case "ne":
		return valueInt != int64(filterVal), nil
	case "lt":
		return valueInt < int64(filterVal), nil
	case "gt":
		return valueInt > int64(filterVal), nil
	case "le":
		return valueInt <= int64(filterVal), nil
	case "ge":
		return valueInt >= int64(filterVal), nil
	}
	return false, errors.New("invalid filter")
}

func parseFilter(form string) (Filter, error) {
	filter := Filter{}
	parts := strings.Split(form, " ")
	if len(parts) < 3 {
		return Filter{}, errors.New(fmt.Sprintf("invalid filter: %s", form))
	}
	filter.Key = parts[0]
	filter.Op = parts[1]
	for i, p := range parts[2:] {
		filter.Val += p
		if (i + 2) != len(parts)-1 {
			filter.Val += " "
		}
	}

	for _, o := range validOps {
		if filter.Op == o {
			return filter, nil
		}
	}
	return Filter{}, errors.New(fmt.Sprintf("invalid filter operator: %s", form))
}
