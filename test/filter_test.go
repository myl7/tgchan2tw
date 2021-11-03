package test

import (
	"fmt"
	"github.com/myl7/tgchan2tw/pkg/fetch"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func TestFilterText(t *testing.T) {
	f, err := os.ReadFile("data/items.yaml")
	if err != nil {
		t.Error(err)
	}

	var items Items
	err = yaml.UnmarshalStrict(f, &items)
	if err != nil {
		t.Error(err)
	}

	for i := range items.Items {
		item := items.Items[i]
		info, err := fetch.FilterText(item.Body, "https://t.me/myl7s/543")
		if err != nil {
			t.Error(err)
		}

		if !item.Info.EqItemBody(info) {
			t.Errorf("item %s failed: required:\n%s\nvs got:\n%s\n", item.Title, fmt.Sprint(item.Info), fmt.Sprint(info))
		}
	}
}
