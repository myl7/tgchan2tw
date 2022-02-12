package test

import (
	"fmt"
	"github.com/myl7/tgchan2tw/pkg/cfg"
	"github.com/myl7/tgchan2tw/pkg/tg"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func TestFilterText(t *testing.T) {
	err := cfg.LoadConfig()
	if err != nil {
		t.Error(err)
	}

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
		info := tg.FilterText(item.Body, "https://t.me/myl7s/543")
		if !item.Info.EqItemBody(info) {
			t.Errorf("item %s failed: required:\n%s\nvs got:\n%s\n", item.Title, fmt.Sprint(item.Info), fmt.Sprint(info))
		}
	}
}
