package list

import (
	"fmt"
	"log"
	"reflect"

	"github.com/masnax/rest-api/cli/cmd/rest"
	"github.com/masnax/rest-api/collection"
	"github.com/mitchellh/mapstructure"
)

func GetCollectionStatList(sourceUrl string, path string, argPath string) ([]string, [][]string) {
	url := sourceUrl + path + argPath
	data := []collection.BookCollection{}
	res, err := rest.MakeRequest(url, "GET", nil)
	mapstructure.Decode(res, &data)
	if err != nil {
		log.Fatal(err)
	}

	out := [][]string{}
	keys := []string{}
	r := reflect.ValueOf(collection.BookCollection{})
	for i := 0; i < r.NumField(); i++ {
		keys = append(keys, r.Type().Field(i).Name)
	}
	for _, b := range data {
		r := reflect.ValueOf(b)
		row := []string{}
		for i := 0; i < r.NumField(); i++ {
			row = append(row, fmt.Sprint(r.Field(i).Interface()))
		}
		out = append(out, row)
	}
	return keys, out
}

func GetCollectionList(sourceUrl string, path string, argPath string) ([]string, [][]string) {
	url := sourceUrl + path + argPath
	data := []collection.Collection{}
	res, err := rest.MakeRequest(url, "GET", nil)
	mapstructure.Decode(res, &data)
	if err != nil {
		log.Fatal(err)
	}

	out := [][]string{}
	keys := []string{}
	r := reflect.ValueOf(collection.Collection{})
	for i := 0; i < r.NumField(); i++ {
		keys = append(keys, r.Type().Field(i).Name)
	}
	for _, b := range data {
		r := reflect.ValueOf(b)
		row := []string{}
		for i := 0; i < r.NumField(); i++ {
			row = append(row, fmt.Sprint(r.Field(i).Interface()))
		}
		out = append(out, row)
	}
	return keys, out
}
