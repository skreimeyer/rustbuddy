// Package crates provides an interface to the crates.io api. The client focus
// is providing quick access to crates metadata
package crates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// FetchCrate will take the name of a crate and return a struct of the metadata
// from the crates.io api
func FetchCrate(crate string) (CrateData, error) {
	var data CrateData
	url := "https://crates.io/api/v1/crates/" + crate
	res, err := http.Get(url)
	if err != nil {
		return data, err
	}
	bod, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(bod, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// FetchDeps will return a list of dependencies given the name of a crate and
// its desired version
func FetchDeps(crate, ver string) ([]Dependencies, error) {
	var deps []Dependencies
	var data DepRoot
	url := fmt.Sprintf("https://crates.io/api/v1/crates/%s/%s/dependencies", crate, ver)
	res, err := http.Get(url)
	if err != nil {
		return deps, err
	}
	bod, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return deps, err
	}
	err = json.Unmarshal(bod, &data)
	if err != nil {
		return deps, err
	}
	return data.Dependencies, nil
}
