package hql

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	DEFAULT_RESULT_SIZE int = 250
	DEFAULT_SORT_ORDER      = "asc"
	DEFAULT_SORT_BY         = "Id"
)

type FilterOptions struct {
	Offset    int
	Limit     int
	SortBy    string
	SortOrder string
	Filter    map[string][]string
}

func (fo *FilterOptions) AddFilter(key string, value interface{}) {
	if fo.Filter == nil {
		fo.Filter = map[string][]string{}
	}

	if _, ok := fo.Filter[key]; ok {
		fo.Filter[key] = append(fo.Filter[key], value)
		return
	}
	fo.Filter[key] = []interface{}{value}
}

func ParseFilterOptionsFromHttpRequest(r *http.Request) (FilterOptions, error) {
	return ParseFilterOptions(r.URL.Query())
}

func ParseFilterOptions(r map[string][]string) (filterOpts FilterOptions, err error) {

	filterOpts = DefaultFilterOptions()

	for k, v := range r {
		switch k {
		case "offset":
			var offset int64
			if offset, err = strconv.ParseInt(v[0], 10, 32); err == nil {
				filterOpts.Offset = int(offset)
			}
		case "limit":
			var limit int64
			if limit, err = strconv.ParseInt(v[0], 10, 32); err == nil {
				filterOpts.Limit = int(limit)
			}
		case "sort":
			filterOpts.SortBy, filterOpts.SortOrder, err = parseSortOptions(v[0])
		default:
			for _, iv := range v {
				filterOpts.AddFilter(k, iv)
			}
		}

		if err != nil {
			return
		}
	}

	return
}

func DefaultFilterOptions() FilterOptions {
	return FilterOptions{
		Offset:    0,
		Limit:     DEFAULT_RESULT_SIZE,
		SortBy:    DEFAULT_SORT_BY,
		SortOrder: DEFAULT_SORT_ORDER,
		Filter:    map[string][]interface{}{},
	}
}

func parseSortOptions(sortOpts string) (by, order string, err error) {

	sbo := strings.Split(sortOpts, ":")
	if len(sbo) == 2 {
		so := strings.ToLower(sbo[1])
		if sbo[0] != "" && (so == "asc" || so == "desc") {
			by = sbo[0]
			order = so
		} else {
			err = fmt.Errorf("Invalid sort options: %s", sortOpts)
		}
	} else if len(sbo) == 1 {
		so := strings.ToLower(sbo[0])
		if so == "asc" || so == "desc" {
			order = so
			by = DEFAULT_SORT_BY
		} else if sbo[0] != "" {
			by = sbo[0]
			order = DEFAULT_SORT_ORDER
		}
	} else {
		by = DEFAULT_SORT_BY
		order = DEFAULT_SORT_ORDER
	}

	return
}
