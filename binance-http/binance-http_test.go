package binancehttp

import (
	"fmt"
	"testing"
)

func TestKlineService(t *testing.T) {
	res, err := KlinesService(
		"ETHUSDT",
		"1m",
		500,
	)
	if err != nil {
		t.Log(err)
	}

	for i := range res {
		fmt.Println(res[i])
	}
}
