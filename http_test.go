package js

import "testing"

func TestLoad(t *testing.T) {

	vv, err := Load("https://api.binance.com/api/v3/ticker/price")

	tickers := map[string]float64{}
	for _, obj := range vv.Objects() {
		tickers[obj.GetStr("symbol")] = obj.GetNum("price")
	}

	require(t, err == nil)
	require(t, len(tickers) > 100)
	require(t, tickers["BTCUSDT"] > 10_000)
}

func TestLoadObject(t *testing.T) {

	obj, err := LoadObject("https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT")

	require(t, err == nil)
	require(t, obj.GetStr("symbol") == "BTCUSDT")
	require(t, obj.GetNum("price") > 10_000)
}
