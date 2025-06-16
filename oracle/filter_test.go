package oracle

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	"github.com/kiichain/price-feeder/config"
	"github.com/kiichain/price-feeder/oracle/provider"
	"github.com/kiichain/price-feeder/oracle/types"
)

func TestSuccessFilterCandleDeviations(t *testing.T) {
	providerCandles := make(provider.AggregatedProviderCandles, 4)
	pair := types.CurrencyPair{
		Base:  "ATOM",
		Quote: "USDT",
	}

	atomPrice := math.LegacyMustNewDecFromStr("29.93")
	atomVolume := math.LegacyMustNewDecFromStr("1994674.34000000")

	atomCandlePrice := []provider.CandlePrice{
		{
			Price:     atomPrice,
			Volume:    atomVolume,
			TimeStamp: provider.PastUnixTime(1 * time.Minute),
		},
	}

	providerCandles[config.ProviderBinance] = map[string][]provider.CandlePrice{
		pair.Base: atomCandlePrice,
	}
	providerCandles[config.ProviderHuobi] = map[string][]provider.CandlePrice{
		pair.Base: atomCandlePrice,
	}
	providerCandles[config.ProviderKraken] = map[string][]provider.CandlePrice{
		pair.Base: atomCandlePrice,
	}
	providerCandles[config.ProviderCoinbase] = map[string][]provider.CandlePrice{
		pair.Base: {
			{
				Price:     math.LegacyMustNewDecFromStr("27.1"),
				Volume:    atomVolume,
				TimeStamp: provider.PastUnixTime(1 * time.Minute),
			},
		},
	}

	pricesFiltered, err := FilterCandleDeviations(
		zerolog.Nop(),
		providerCandles,
		make(map[string]math.LegacyDec),
	)

	_, ok := pricesFiltered[config.ProviderCoinbase]
	require.NoError(t, err, "It should successfully filter out the provider using candles")
	require.False(t, ok, "The filtered candle deviation price at coinbase should be empty")

	customDeviations := make(map[string]math.LegacyDec, 1)
	customDeviations[pair.Base] = math.LegacyNewDec(2)

	pricesFilteredCustom, err := FilterCandleDeviations(
		zerolog.Nop(),
		providerCandles,
		customDeviations,
	)

	_, ok = pricesFilteredCustom[config.ProviderCoinbase]
	require.NoError(t, err, "It should successfully not filter out coinbase")
	require.True(t, ok, "The filtered candle deviation price of coinbase should remain")
}

func TestSuccessFilterTickerDeviations(t *testing.T) {
	providerTickers := make(provider.AggregatedProviderPrices, 4)
	pair := types.CurrencyPair{
		Base:  "ATOM",
		Quote: "USDT",
	}

	atomPrice := math.LegacyMustNewDecFromStr("29.93")
	atomVolume := math.LegacyMustNewDecFromStr("1994674.34000000")

	atomTickerPrice := provider.TickerPrice{
		Price:  atomPrice,
		Volume: atomVolume,
	}

	providerTickers[config.ProviderBinance] = map[string]provider.TickerPrice{
		pair.Base: atomTickerPrice,
	}
	providerTickers[config.ProviderHuobi] = map[string]provider.TickerPrice{
		pair.Base: atomTickerPrice,
	}
	providerTickers[config.ProviderKraken] = map[string]provider.TickerPrice{
		pair.Base: atomTickerPrice,
	}
	providerTickers[config.ProviderCoinbase] = map[string]provider.TickerPrice{
		pair.Base: {
			Price:  math.LegacyMustNewDecFromStr("27.1"),
			Volume: atomVolume,
		},
	}

	pricesFiltered, err := FilterTickerDeviations(
		zerolog.Nop(),
		providerTickers,
		make(map[string]math.LegacyDec),
	)

	_, ok := pricesFiltered[config.ProviderCoinbase]
	require.NoError(t, err, "It should successfully filter out the provider using tickers")
	require.False(t, ok, "The filtered ticker deviation price at coinbase should be empty")

	customDeviations := make(map[string]math.LegacyDec, 1)
	customDeviations[pair.Base] = math.LegacyNewDec(2)

	pricesFilteredCustom, err := FilterTickerDeviations(
		zerolog.Nop(),
		providerTickers,
		customDeviations,
	)

	_, ok = pricesFilteredCustom[config.ProviderCoinbase]
	require.NoError(t, err, "It should successfully not filter out coinbase")
	require.True(t, ok, "The filtered candle deviation price of coinbase should remain")
}
