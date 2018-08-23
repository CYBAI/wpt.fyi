// Copyright 2018 The WPT Dashboard Project. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package query

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"github.com/web-platform-tests/wpt.fyi/shared"
	"google.golang.org/appengine"
)

var (
	autocompleteDefaultLimit = 1
	autocompleteMinLimit     = 1
	autocompleteMaxLimit     = 50
)

// AutocompleteResult contains a single autocomplete suggestion.
type AutocompleteResult struct {
	// QueryString represents the most basic form of an autocomplete result. It is
	// the complete query string being recommended. The client is responsibile for
	// fitting it to any substring(s) already appearing in a search UI.
	QueryString string `json:"query_string"`
}

// AutocompleteResponse contains a response to autocmplete API calls.
type AutocompleteResponse struct {
	// Suggestions is the collection of autocomplete suggestions.
	Suggestions []AutocompleteResult `json:"results"`
}

// byQueryIndex sorts by strings.Index(r.rs[-].QueryString, r.q). If the
// substring index of q is the same in both QueryStrings, then sorting falls
// back on Less(i, j) = r.rs[i].QueryString < r.rs[j].QueryString.
type byQueryIndex struct {
	q  string
	rs []AutocompleteResult
}

func (r byQueryIndex) Len() int      { return len(r.rs) }
func (r byQueryIndex) Swap(i, j int) { r.rs[i], r.rs[j] = r.rs[j], r.rs[i] }
func (r byQueryIndex) Less(i, j int) bool {
	a := strings.Index(r.rs[i].QueryString, r.q)
	b := strings.Index(r.rs[j].QueryString, r.q)
	if a == b {
		return r.rs[i].QueryString < r.rs[j].QueryString
	}
	return strings.Index(r.rs[i].QueryString, r.q) < strings.Index(r.rs[j].QueryString, r.q)
}

type autocompleteHandler struct {
	queryHandler
}

func apiAutocompleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	sh := autocompleteHandler{queryHandler{
		sharedImpl: defaultShared{ctx},
		dataSource: cachedStore{
			cache: gzipReadWritable{memcacheReadWritable{ctx}},
			store: httpReadable{ctx},
		},
	}}
	sh.ServeHTTP(w, r)
}

func (ah autocompleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	limit, filters, testRuns, summaries, err := ah.processInput(w, r)
	// processInput handles writing any error to w.
	if err != nil {
		return
	}

	resp := prepareAutocompleteResponse(limit, filters, testRuns, summaries)

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(data)
}

func (ah autocompleteHandler) processInput(w http.ResponseWriter, r *http.Request) (int, *shared.QueryFilter, []shared.TestRun, []summary, error) {
	limit, err := ah.parseLimit(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 0, nil, nil, nil, err
	}

	filter, testRuns, summaries, err := ah.queryHandler.processInput(w, r)
	return limit, filter, testRuns, summaries, err
}

func (ah autocompleteHandler) parseLimit(r *http.Request) (int, error) {
	limit, err := ah.sharedImpl.ParseQueryParamInt(r, "limit")
	if err == shared.ErrMissing {
		return autocompleteDefaultLimit, nil
	}
	if err != nil {
		return 0, err
	}

	if limit < autocompleteMinLimit {
		return autocompleteMinLimit, nil
	}
	if limit > autocompleteMaxLimit {
		return autocompleteMaxLimit, nil
	}
	return limit, nil
}

func prepareAutocompleteResponse(limit int, filters *shared.QueryFilter, testRuns []shared.TestRun, summaries []summary) AutocompleteResponse {
	fileSet := mapset.NewSet()
	for _, smry := range summaries {
		for file := range smry {
			fileSet.Add(file)
		}
	}

	files := make([]AutocompleteResult, 0, fileSet.Cardinality()/len(testRuns))
	for fileInterface := range fileSet.Iter() {
		file := fileInterface.(string)
		if strings.Contains(file, filters.Q) {
			files = append(files, AutocompleteResult{file})
		}
	}

	sortable := byQueryIndex{
		q:  filters.Q,
		rs: files,
	}
	sort.Sort(sortable)
	if len(sortable.rs) > limit {
		sortable.rs = sortable.rs[:limit]
	}
	return AutocompleteResponse{sortable.rs}
}
