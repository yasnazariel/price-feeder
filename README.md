# Kiichain Oracle Price Feeder

This is the Kiichain version of the incredible work of SEI on the price feeder module. The original implementation is available at [SEI's price-feeder](https://github.com/sei-protocol/sei-chain/tree/main/oracle/price-feeder).

Further documentation for the project can be found at [Running the Price feeder](https://docs.kiiglobal.io/docs/validate-the-network/run-a-validator-full-node/running-the-price-feeder).

# Setup

If a cluster is running Oracle price-feeder, your validator is also required to run a price feeder or your validator will be jailed for missing votes.

## Bootstrap

We highly recommend bootstrapping the price-feeder installation using our installer:

```bash
wget https://raw.githubusercontent.com/KiiChain/testnets/main/testnet_oro/run_price_feeder.sh
chmod +x run_price_feeder.sh
./run_price_feeder.sh
```

## Create an account for Oracle Price Feeder Delegate

1. To avoid account sequence errors and security problems with the admin account, it's recommended to create a different account as an Oracle delegate. To do so, you'll need to create the account with:

```bash
kiichaind keys add price-feeder-delegate
```

Or use any arbitrary key name. This project may still cause account sequence errors for the delegate account but since it's only being used for the Oracle price feeder, it's not a concern

2. With the account address output run:

```bash
export PRICE_FEEDER_DELEGATE_ADDR=<output>
```

3. To define the feeder, you can execute:

```bash
kiichaind tx oracle set-feeder $PRICE_FEEDER_DELEGATE_ADDR --from <validator-wallet> --fees 10000000000000000akii -b block -y --chain-id {chain-id}
```

4. Make sure to send bank a tiny amount to the account in order for the account to be created:

```bash
kiichaind tx bank send [VALIDATOR_ACCOUNT] $PRICE_FEEDER_DELEGATE_ADDR --from [VALIDATOR_ACCOUNT] [AMOUNT] --fees 10000000000000000akii -b block -y
```

Then you need to export `PRICE_FEEDER_PASS` environment variable to set up the keyring password. That was entered during the account setup.

Ex:

```bash
export PRICE_FEEDER_PASS=keyringPassword
```

If this environment variable is not set, the price feeder will prompt the user for input.

## Build or install Price Feeder

To build the price feeder, run the following command from the root of the Git repository:

```bash
make build
```

To install run:

```bash
make install
```

## Run Price Feeder

You can run it as a separate binary but it's recommended to run it as a system daemon service, you can use the following as an example.

You need to setup the config.toml file (see [this for example](./config.example.toml)), you need to set the following fields in:

```bash
...
[account]
address = "<UPDATE ME>"  <-- $PRICE_FEEDER_DELEGATE_ADDR from above
validator = "<UPDATE ME>" <-- validator address
...
```

Finally make sure you are having the correct keyring type set on the config.toml

```bash
[keyring]
backend = "os" <-- check this is the same you are using
dir = "~/.kiichain3"
```

After finishing the config.toml, be sure you are on the root of this project and then type the following command to manually run the price-feeder:

```bash
price_feeder start oracle/price_feeder/config.toml
```

## HTTP server

A HTTP server can be enabled on the price feeder.
The server will expose the following endpoints:

- `/healthz`: A simple health check endpoint that returns a 200 OK response.
- `/prices`: Returns the current prices fetched from the oracle's set of exchange rate providers.
- `/metrics`: Returns the current metrics collected by the price feeder, including prices and their timestamps.

### HTTP server configuration

The HTTP server can be configured in the `config.toml` file under the `server` section. The following options are available:

```toml
# This is the main configuration for the price feeder module.
[main]
# Define if the price feeder should send votes to the chain
enable_voting = true
# Defines if the price feeder server is enabled
enable_server = true

# Defines the server configuration
[server]
# The address where the server will listen for HTTP requests
listen_addr = "0.0.0.0:7171"
# The timeout for read operations
read_timeout = "20s"
# The timeout for write operations
write_timeout = "20s"
# Define if cors is enabled
enable_cors = true
# The allowed origins for CORS requests
allowed_origins = ["*"]
```

If CORS is enabled, the server will allow requests from any origin. You can restrict this by modifying the `allowed_origins` field to include only specific origins.

## Providers

The list of current supported providers:

- [Binance](https://www.binance.com/en)
- [MEXC](https://www.mexc.com/)
- [Coinbase](https://www.coinbase.com/)
- [Gate](https://www.gate.io/)
- [Huobi](https://www.huobi.com/en-us/)
- [Kraken](https://www.kraken.com/en-us/)
- [Okx](https://www.okx.com/)

## Usage

The `price-feeder` tool runs off of a single configuration file. This configuration
file defines what exchange rates to fetch and what providers to get them from.
In addition, it defines the oracle's keyring and feeder account information.
The keyring's password is defined via environment variables or user input.
More information on the keyring can be found [here](#keyring)
Please see the [example configuration](./config.example.toml) for more details.

```shell
$ price-feeder start /path/to/price_feeder_config.toml
```

## Configuration

### telemetry

A set of options for the application's telemetry, which is disabled by default. An in-memory sink is the default, but Prometheus is also supported. We use the [cosmos sdk telemetry package](https://github.com/cosmos/cosmos-sdk/blob/main/docs/core/telemetry.md).

### deviation

Deviation allows validators to set a custom amount of standard deviations around the median which is helpful if any providers become faulty. It should be noted that the default for this option is 1 standard deviation.

### provider_endpoints

The provider_endpoints option enables validators to setup their own API endpoints for a given provider.

### currency_pairs

The `currency_pairs` sections contains one or more exchange rates along with the
providers from which to get market data from. It is important to note that the
providers supplied in each `currency_pairs` must support the given exchange rate.

For example, to get multiple price points on ATOM, you could define `currency_pairs`
as follows:

```toml
[[currency_pairs]]
base = "ATOM"
providers = [
  "binance",
]
quote = "USDT"

[[currency_pairs]]
base = "ATOM"
providers = [
  "kraken",
]
quote = "USD"
```

Providing multiple providers is beneficial in case any provider fails to return
market data. Prices per exchange rate are submitted on-chain via pre-vote and
vote messages using a time-weighted average price (TVWAP).

### account

The `account` section contains the oracle's feeder and validator account information.
These are used to sign and populate data in pre-vote and vote oracle messages.

### keyring

The `keyring` section contains Keyring related material used to fetch the key pair
associated with the oracle account that signs pre-vote and vote oracle messages.

### rpc

The `rpc` section contains the Tendermint and Cosmos application gRPC endpoints.
These endpoints are used to query for on-chain data that pertain to oracle
functionality and for broadcasting signed pre-vote and vote oracle messages.

### healthchecks

The `healthchecks` section defines optional healthcheck endpoints to ping on successful
oracle votes. This provides a simple alerting solution which can integrate with a service
like [healthchecks.io](https://healthchecks.io). It's recommended to configure additional
monitoring since third-party services can be unreliable.

## Keyring

Our keyring must be set up to sign transactions before running the price feeder.
Additional info on the different keyring modes is available [here](https://docs.cosmos.network/master/run-node/keyring.html).
**Please note that the `test` and `memory` modes are only for testing purposes.**
**Do not use these modes for running the price feeder against mainnet.**
