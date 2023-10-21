package main

import (
	"brand/b"
	"brand/data"
	"fmt"
	"testing"

	"github.com/chiyoi/apricot/logs"
)

func TestData(t *testing.T) {
	data.Load()
	defer data.Save()
	id, title, err := b.GetLatest()
	if err != nil {
		logs.Panic(err)
	}
	if data.Data.LatestID != id {
		fmt.Println("Update.", id, title)
		return
	}
	fmt.Println("No update.", id, title)
}
