package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	apiHost = "api.entropedia.net"
)

type Resource struct {
	Url      string
	Content  []byte
	Id       string `json:"id"`
	Sha256   string `json:"sha256"`
	DataSize int64  `json:"dataSize"`
}

type ResourcesResponse struct {
	Resources []Resource
}

func (res *Resource) Fetch() error {
	resp, err := http.Get("http://" + apiHost + ":8999/v1/resources?sha256=" + res.Sha256)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r := ResourcesResponse{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return errors.New("Sorry, but entropedia has no information about this file - just yet. Gobble it now!")
	}

	*res = r.Resources[0]
	return nil
}

func main() {
	var err error
	flag.Parse()
	path := flag.Arg(0)
	if len(path) == 0 {
		fmt.Println("Usage: extropy [path]")
		return
	}

	res := Resource{}

	fmt.Printf("- Reducing entropy for: %s\n", path)

	checksum := sha256.New()
	res.Content, err = ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("\t! Could not sha256sum:", path)
		return
	}
	checksum.Write(res.Content)
	res.Sha256 = hex.EncodeToString(checksum.Sum(nil))

	fmt.Printf("\t+ SHA256: %s\n", res.Sha256)
	err = res.Fetch()
	if err != nil {
		fmt.Printf("\t! Error: %s\n", err)
		return
	} else {
		fmt.Printf("\t+ Found: %s, size: %d\n", res.Id, res.DataSize)
	}

	fmt.Println()
	fmt.Println("- Reduced entropy successfully!")
}
