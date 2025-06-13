package oracle

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/kiichain/price-feeder/config"
	"github.com/kiichain/price-feeder/oracle/provider"

	"github.com/stretchr/testify/require"
)

func TestComputeVWAP(t *testing.T) {
	testCases := map[string]struct {
		prices   map[string]map[string]provider.TickerPrice
		expected map[string]math.LegacyDec
	}{
		"empty prices": {
			prices:   make(map[string]map[string]provider.TickerPrice),
			expected: make(map[string]math.LegacyDec),
		},
		"nil prices": {
			prices:   nil,
			expected: make(map[string]math.LegacyDec),
		},
		"non empty prices": {
			prices: map[string]map[string]provider.TickerPrice{
				config.ProviderBinance: {
					"ATOM": provider.TickerPrice{
						Price:  math.LegacyMustNewDecFromStr("28.21000000"),
						Volume: math.LegacyMustNewDecFromStr("2749102.78000000"),
					},
					"UMEE": provider.TickerPrice{
						Price:  math.LegacyMustNewDecFromStr("1.13000000"),
						Volume: math.LegacyMustNewDecFromStr("249102.38000000"),
					},
					"KII": provider.TickerPrice{
						Price:  math.LegacyMustNewDecFromStr("64.87000000"),
						Volume: math.LegacyMustNewDecFromStr("7854934.69000000"),
					},
				},
				config.ProviderKraken: {
					"ATOM": provider.TickerPrice{
						Price:  math.LegacyMustNewDecFromStr("28.268700"),
						Volume: math.LegacyMustNewDecFromStr("178277.53314385"),
					},
					"KII": provider.TickerPrice{
						Price:  math.LegacyMustNewDecFromStr("64.87853000"),
						Volume: math.LegacyMustNewDecFromStr("458917.46353577"),
					},
				},
				"FOO": {
					"ATOM": provider.TickerPrice{
						Price:  math.LegacyMustNewDecFromStr("28.168700"),
						Volume: math.LegacyMustNewDecFromStr("4749102.53314385"),
					},
				},
			},
			expected: map[string]math.LegacyDec{
				"ATOM": math.LegacyMustNewDecFromStr("28.185812745610043621"),
				"UMEE": math.LegacyMustNewDecFromStr("1.13000000"),
				"KII":  math.LegacyMustNewDecFromStr("64.870470848638112395"),
			},
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			vwap, err := ComputeVWAP(tc.prices)
			require.NoError(t, err)
			require.Len(t, vwap, len(tc.expected))

			for k, v := range tc.expected {
				require.Equalf(t, v, vwap[k], "unexpected VWAP for %s", k)
			}
		})
	}
}

func TestComputeTVWAP(t *testing.T) {
	now := provider.PastUnixTime(0)
	testCases := []struct {
		name     string
		prices   provider.AggregatedProviderCandles
		expected map[string]math.LegacyDec
		now      int64
	}{
		{
			name:     "empty candles",
			prices:   make(provider.AggregatedProviderCandles),
			expected: make(map[string]math.LegacyDec),
		},
		{
			name:     "nil candles",
			prices:   nil,
			expected: make(map[string]math.LegacyDec),
		},
		{
			name: "non-empty candles",
			prices: provider.AggregatedProviderCandles{
				config.ProviderBinance: {
					"ATOM": []provider.CandlePrice{
						{TimeStamp: provider.PastUnixTime(600), Price: math.LegacyMustNewDecFromStr("28.21000000"), Volume: math.LegacyMustNewDecFromStr("1000")},
						{TimeStamp: provider.PastUnixTime(300), Price: math.LegacyMustNewDecFromStr("28.23000000"), Volume: math.LegacyMustNewDecFromStr("1000")},
					},
					"UMEE": []provider.CandlePrice{
						{TimeStamp: provider.PastUnixTime(600), Price: math.LegacyMustNewDecFromStr("1.13000000"), Volume: math.LegacyMustNewDecFromStr("500")},
					},
				},
			},
			expected: map[string]math.LegacyDec{
				"ATOM": math.LegacyMustNewDecFromStr("28.22000000"), // Adjust this after checking the expected result of TVWAP
				"UMEE": math.LegacyMustNewDecFromStr("1.13000000"),
			},
		},
		{
			name: "non-empty candles",
			prices: provider.AggregatedProviderCandles{
				config.ProviderBinance: {
					"ATOM": []provider.CandlePrice{
						{TimeStamp: provider.PastUnixTime(600), Price: math.LegacyMustNewDecFromStr("28.21000000"), Volume: math.LegacyMustNewDecFromStr("1000")},
						{TimeStamp: provider.PastUnixTime(300), Price: math.LegacyMustNewDecFromStr("28.23000000"), Volume: math.LegacyMustNewDecFromStr("1000")},
					},
					"UMEE": []provider.CandlePrice{
						{TimeStamp: provider.PastUnixTime(600), Price: math.LegacyMustNewDecFromStr("1.13000000"), Volume: math.LegacyMustNewDecFromStr("500")},
					},
				},
			},
			expected: map[string]math.LegacyDec{
				"ATOM": math.LegacyMustNewDecFromStr("28.22000000"),
				"UMEE": math.LegacyMustNewDecFromStr("1.13000000"),
			},
		},
		{
			name: "candles with same time as now",
			now:  now,
			prices: provider.AggregatedProviderCandles{
				config.ProviderBinance: {
					"ATOM": []provider.CandlePrice{
						{TimeStamp: now, Price: math.LegacyMustNewDecFromStr("28.50000000"), Volume: math.LegacyMustNewDecFromStr("1000")},
					},
					"UMEE": []provider.CandlePrice{
						{TimeStamp: now, Price: math.LegacyMustNewDecFromStr("1.25000000"), Volume: math.LegacyMustNewDecFromStr("500")},
					},
				},
			},
			expected: map[string]math.LegacyDec{
				"ATOM": math.LegacyMustNewDecFromStr("28.50000000"),
				"UMEE": math.LegacyMustNewDecFromStr("1.25000000"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.now > 0 {
				mockNow = tc.now
			} else {
				mockNow = 0
			}
			tvwap, err := ComputeTVWAP(tc.prices)
			require.NoError(t, err)
			require.Len(t, tvwap, len(tc.expected))

			for k, v := range tc.expected {
				require.Equalf(t, v, tvwap[k], "unexpected TVWAP for %s", k)
			}
		})
	}
}

func TestStandardDeviation(t *testing.T) {
	type deviation struct {
		mean      math.LegacyDec
		deviation math.LegacyDec
	}
	testCases := map[string]struct {
		prices   map[string]map[string]math.LegacyDec
		expected map[string]deviation
	}{
		"empty prices": {
			prices:   make(map[string]map[string]math.LegacyDec),
			expected: map[string]deviation{},
		},
		"nil prices": {
			prices:   nil,
			expected: map[string]deviation{},
		},
		"not enough prices": {
			prices: map[string]map[string]math.LegacyDec{
				config.ProviderBinance: {
					"ATOM": math.LegacyMustNewDecFromStr("28.21000000"),
					"UMEE": math.LegacyMustNewDecFromStr("1.13000000"),
					"KII":  math.LegacyMustNewDecFromStr("64.87000000"),
				},
				config.ProviderKraken: {
					"ATOM": math.LegacyMustNewDecFromStr("28.23000000"),
					"UMEE": math.LegacyMustNewDecFromStr("1.13050000"),
					"KII":  math.LegacyMustNewDecFromStr("64.85000000"),
				},
			},
			expected: map[string]deviation{},
		},
		"some prices": {
			prices: map[string]map[string]math.LegacyDec{
				config.ProviderBinance: {
					"ATOM": math.LegacyMustNewDecFromStr("28.21000000"),
					"UMEE": math.LegacyMustNewDecFromStr("1.13000000"),
					"KII":  math.LegacyMustNewDecFromStr("64.87000000"),
				},
				config.ProviderKraken: {
					"ATOM": math.LegacyMustNewDecFromStr("28.23000000"),
					"UMEE": math.LegacyMustNewDecFromStr("1.13050000"),
				},
				config.ProviderCoinbase: {
					"ATOM": math.LegacyMustNewDecFromStr("28.40000000"),
					"UMEE": math.LegacyMustNewDecFromStr("1.14000000"),
					"KII":  math.LegacyMustNewDecFromStr("64.10000000"),
				},
			},
			expected: map[string]deviation{
				"ATOM": {
					mean:      math.LegacyMustNewDecFromStr("28.28"),
					deviation: math.LegacyMustNewDecFromStr("0.085244745683629475"),
				},
				"UMEE": {
					mean:      math.LegacyMustNewDecFromStr("1.1335"),
					deviation: math.LegacyMustNewDecFromStr("0.004600724580614015"),
				},
			},
		},

		"non empty prices": {
			prices: map[string]map[string]math.LegacyDec{
				config.ProviderBinance: {
					"ATOM": math.LegacyMustNewDecFromStr("28.21000000"),

					"UMEE": math.LegacyMustNewDecFromStr("1.13000000"),
					"KII":  math.LegacyMustNewDecFromStr("64.87000000"),
				},
				config.ProviderKraken: {
					"ATOM": math.LegacyMustNewDecFromStr("28.23000000"),
					"UMEE": math.LegacyMustNewDecFromStr("1.13050000"),
					"KII":  math.LegacyMustNewDecFromStr("64.85000000"),
				},
				config.ProviderCoinbase: {
					"ATOM": math.LegacyMustNewDecFromStr("28.40000000"),
					"UMEE": math.LegacyMustNewDecFromStr("1.14000000"),
					"KII":  math.LegacyMustNewDecFromStr("64.10000000"),
				},
			},
			expected: map[string]deviation{
				"ATOM": {
					mean:      math.LegacyMustNewDecFromStr("28.28"),
					deviation: math.LegacyMustNewDecFromStr("0.085244745683629475"),
				},
				"UMEE": {
					mean:      math.LegacyMustNewDecFromStr("1.1335"),
					deviation: math.LegacyMustNewDecFromStr("0.004600724580614015"),
				},
				"KII": {
					mean:      math.LegacyMustNewDecFromStr("64.606666666666666666"),
					deviation: math.LegacyMustNewDecFromStr("0.358360464089193609"),
				},
			},
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			deviation, mean, err := StandardDeviation(tc.prices)
			require.NoError(t, err)
			require.Len(t, deviation, len(tc.expected))
			require.Len(t, mean, len(tc.expected))

			for k, v := range tc.expected {
				require.Equalf(t, v.deviation, deviation[k], "unexpected deviation for %s", k)
				require.Equalf(t, v.mean, mean[k], "unexpected mean for %s", k)
			}
		})
	}
}
