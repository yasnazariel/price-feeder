package provider

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	"github.com/kiichain/price-feeder/config"
	"github.com/kiichain/price-feeder/oracle/types"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestKrakenProvider_GetTickerPrices(t *testing.T) {
	p, err := NewKrakenProvider(
		context.TODO(),
		zerolog.Nop(),
		config.ProviderEndpoint{},
		types.CurrencyPair{Base: "BTC", Quote: "USDT"},
	)
	require.NoError(t, err)

	t.Run("valid_request_single_ticker", func(t *testing.T) {
		lastPrice := math.LegacyMustNewDecFromStr("34.69000000")
		volume := math.LegacyMustNewDecFromStr("2396974.02000000")

		tickerMap := map[string]TickerPrice{}
		tickerMap["ATOMUSDT"] = TickerPrice{
			Price:  lastPrice,
			Volume: volume,
		}

		p.tickers = tickerMap

		prices, err := p.GetTickerPrices(types.CurrencyPair{Base: "ATOM", Quote: "USDT"})
		require.NoError(t, err)
		require.Len(t, prices, 1)
		require.Equal(t, lastPrice, prices["ATOMUSDT"].Price)
		require.Equal(t, volume, prices["ATOMUSDT"].Volume)
	})

	t.Run("valid_request_multi_ticker", func(t *testing.T) {
		lastPriceAtom := math.LegacyMustNewDecFromStr("34.69000000")
		lastPriceKii := math.LegacyMustNewDecFromStr("41.35000000")
		volume := math.LegacyMustNewDecFromStr("2396974.02000000")

		tickerMap := map[string]TickerPrice{}
		tickerMap["ATOMUSDT"] = TickerPrice{
			Price:  lastPriceAtom,
			Volume: volume,
		}

		tickerMap["KIIUSDT"] = TickerPrice{
			Price:  lastPriceKii,
			Volume: volume,
		}

		p.tickers = tickerMap
		prices, err := p.GetTickerPrices(
			types.CurrencyPair{Base: "ATOM", Quote: "USDT"},
			types.CurrencyPair{Base: "KII", Quote: "USDT"},
		)
		require.NoError(t, err)
		require.Len(t, prices, 2)
		require.Equal(t, lastPriceAtom, prices["ATOMUSDT"].Price)
		require.Equal(t, volume, prices["ATOMUSDT"].Volume)
		require.Equal(t, lastPriceKii, prices["KIIUSDT"].Price)
		require.Equal(t, volume, prices["KIIUSDT"].Volume)
	})

	t.Run("invalid_request_invalid_ticker", func(t *testing.T) {
		prices, err := p.GetTickerPrices(types.CurrencyPair{Base: "FOO", Quote: "BAR"})
		require.NoError(t, err)
		require.Zero(t, len(prices))
	})
}

func TestKrakenProvider_SubscribeCurrencyPairs(t *testing.T) {
	p, err := NewKrakenProvider(
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

func TestKrakenPairToCurrencyPairSymbol(t *testing.T) {
	cp := types.CurrencyPair{Base: "ATOM", Quote: "USDT"}
	currencyPairSymbol := krakenPairToCurrencyPairSymbol("ATOM/USDT")
	require.Equal(t, cp.String(), currencyPairSymbol)
}

func TestKrakenCurrencyPairToKrakenPair(t *testing.T) {
	cp := types.CurrencyPair{Base: "ATOM", Quote: "USDT"}
	krakenSymbol := currencyPairToKrakenPair(cp)
	require.Equal(t, krakenSymbol, "ATOM/USDT")
}

func TestNormalizeKrakenBTCPair(t *testing.T) {
	btcSymbol := normalizeKrakenBTCPair("XBT/USDT")
	require.Equal(t, btcSymbol, "BTC/USDT")

	atomSymbol := normalizeKrakenBTCPair("ATOM/USDT")
	require.Equal(t, atomSymbol, "ATOM/USDT")
}
