package futures

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/internal/logger"
	"net/http"
	"net/url"
)

type usdtFuturesMarket struct {
	*USDTFutures
}

func (f *usdtFuturesMarket) GetName() string {
	return "hbdm.com"
}

func (f *usdtFuturesMarket) DoNoAuthRequest(method, reqUrl string, params *url.Values) ([]byte, error) {
	cli := GetHttpCli()
	if method == http.MethodGet {
		reqUrl += "?" + params.Encode()
	}

	respBodyData, err := cli.DoRequest(method, reqUrl, "", map[string]string{
		"Content-Type": "application/json",
	})

	if err != nil {
		return nil, err
	}

	var baseResp BaseResponse
	err = json.Unmarshal(respBodyData, &baseResp)
	if err != nil {
		logger.Errorf("[DoNoAuthRequest] err=%s", err.Error())
		return nil, err
	}

	if baseResp.Status != "ok" {
		return nil, errors.New(string(respBodyData))
	}

	return respBodyData, nil
}

func (f *usdtFuturesMarket) GetDepth(pair CurrencyPair, limit int, opt ...OptionParameter) (*Depth, error) {
	//TODO implement me
	panic("implement me")
}

func (f *usdtFuturesMarket) GetTicker(pair CurrencyPair, opts ...OptionParameter) (*Ticker, error) {
	params := url.Values{}
	params.Set("contract_code", pair.Symbol)
	MergeOptionParams(&params, opts...)

	data, err := f.DoNoAuthRequest(http.MethodGet,
		fmt.Sprintf("%s%s", f.uriOpts.Endpoint, f.uriOpts.TickerUri), &params)
	if err != nil {
		return nil, err
	}

	tk, err := f.unmarshalerOpts.TickerUnmarshaler(data)
	if err != nil {
		return nil, err
	}

	tk.Pair = pair
	tk.Origin = data

	return tk, nil
}

func (f *usdtFuturesMarket) GetKline(pair CurrencyPair, period KlinePeriod, opts ...OptionParameter) ([]Kline, error) {
	params := url.Values{}
	params.Set("contract_code", pair.Symbol)
	params.Set("period", AdaptKlinePeriod(period))

	MergeOptionParams(&params, opts...)

	if params.Get("size") == "" && params.Get("from") == "" {
		params.Set("size", "100")
	}

	data, err := f.DoNoAuthRequest(http.MethodGet, fmt.Sprintf("%s%s", f.uriOpts.Endpoint, f.uriOpts.KlineUri), &params)
	if err != nil {
		return nil, err
	}
	logger.Debugf("[GetKline] data=%s", string(data))

	klines, err := f.unmarshalerOpts.KlineUnmarshaler(data)
	if err != nil {
		return nil, err
	}

	for i, _ := range klines {
		klines[i].Pair = pair
	}

	return klines, err
}