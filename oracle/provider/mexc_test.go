package provider

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	"github.com/kiichain/price-feeder/config"
	"github.com/kiichain/price-feeder/oracle/types"
)

func TestMexcProvider_GetTickerPrices(t *testing.T) {
	t.Skip("skipping until mexc websocket endpoint is restored")
	p, err := NewMexcProvider(
		context.TODO(),
		zerolog.Nop(),
		config.ProviderEndpoint{},
		types.CurrencyPair{Base: "ATOM", Quote: "USDT"},
	)
	require.NoError(t, err)

	t.Run("valid_request_single_ticker", func(t *testing.T) {
		lastPrice := "34.69000000"
		volume := "2396974.02000000"

		tickerMap := map[string]MexcTicker{}
		tickerMap["ATOMUSDT"] = MexcTicker{
			Symbol:    "ATOMUSDT",
			LastPrice: lastPrice,
			Volume:    volume,
		}

		p.tickers = tickerMap

		prices, err := p.GetTickerPrices(types.CurrencyPair{Base: "ATOM", Quote: "USDT"})
		require.NoError(t, err)
		require.Len(t, prices, 1)
		require.Equal(t, math.LegacyMustNewDecFromStr(lastPrice), prices["ATOMUSDT"].Price)
		require.Equal(t, math.LegacyMustNewDecFromStr(volume), prices["ATOMUSDT"].Volume)
	})

	t.Run("valid_request_multi_ticker", func(t *testing.T) {
		lastPriceAtom := "34.69000000"
		lastPriceKii := "41.35000000"
		volume := "2396974.02000000"

		tickerMap := map[string]MexcTicker{}
		tickerMap["ATOMUSDT"] = MexcTicker{
			Symbol:    "ATOMUSDT",
			LastPrice: lastPriceAtom,
			Volume:    volume,
		}

		tickerMap["KIIUSDT"] = MexcTicker{
			Symbol:    "KIIUSDT",
			LastPrice: lastPriceKii,
			Volume:    volume,
		}

		p.tickers = tickerMap
		prices, err := p.GetTickerPrices(
			types.CurrencyPair{Base: "ATOM", Quote: "USDT"},
			types.CurrencyPair{Base: "KII", Quote: "USDT"},
		)
		require.NoError(t, err)
		require.Len(t, prices, 2)
		require.Equal(t, math.LegacyMustNewDecFromStr(lastPriceAtom), prices["ATOMUSDT"].Price)
		require.Equal(t, math.LegacyMustNewDecFromStr(volume), prices["ATOMUSDT"].Volume)
		require.Equal(t, math.LegacyMustNewDecFromStr(lastPriceKii), prices["KIIUSDT"].Price)
		require.Equal(t, math.LegacyMustNewDecFromStr(volume), prices["KIIUSDT"].Volume)
	})

	t.Run("invalid_request_invalid_ticker", func(t *testing.T) {
		prices, err := p.GetTickerPrices(types.CurrencyPair{Base: "FOO", Quote: "BAR"})
		require.NoError(t, err)
		require.Zero(t, len(prices))
	})
}

func TestMexcProvider_SubscribeCurrencyPairs(t *testing.T) {
	t.Skip("skipping until mexc websocket endpoint is restored")
	p, err := NewMexcProvider(
		context.TODO(),
		zerolog.Nop(),
		config.ProviderEndpoint{},
		types.CurrencyPair{Base: "ATOM", Quote: "USDT"},
	)
	require.NoError(t, err)

	t.Run("invalid_subscribe_channels_empty", func(t *testing.T) {
		err = p.SubscribeCurrencyPairs([]types.CurrencyPair{}...)
		require.ErrorContains(t, err, "currency pairs is empty")
	})
}

func TestMexcCurrencyPairToMexcPair(t *testing.T) {
	cp := types.CurrencyPair{Base: "ATOM", Quote: "USDT"}
	MexcSymbol := currencyPairToMexcPair(cp)
	require.Equal(t, MexcSymbol, "ATOM_USDT")
}
