#!/bin/bash
set -ux

# always returns true so set -e doesn't exit if it is not running.
killall oraid || true
rm -rf $HOME/.oraid/
killall screen

# make four orai directories
mkdir $HOME/.oraid
mkdir $HOME/.oraid/validator1

# init all three validators
oraid init --chain-id=testing validator1 --home=$HOME/.oraid/validator1

# create keys for all three validators
oraid keys add validator1 --keyring-backend test --home=$HOME/.oraid/validator1

update_genesis () {    
    cat $HOME/.oraid/validator1/config/genesis.json | jq "$1" > $HOME/.oraid/validator1/config/tmp_genesis.json && mv $HOME/.oraid/validator1/config/tmp_genesis.json $HOME/.oraid/validator1/config/genesis.json
}

# change staking denom to orai
update_genesis '.app_state["staking"]["params"]["bond_denom"]="orai"'

# create validator node 1
oraid genesis add-genesis-account $(oraid keys show validator1 -a --keyring-backend=test --home=$HOME/.oraid/validator1 | grep orai) 1000000000000orai,1000000000000stake --home=$HOME/.oraid/validator1
oraid genesis gentx validator1 500000000orai --keyring-backend=test --home=$HOME/.oraid/validator1 --chain-id=testing
oraid genesis collect-gentxs --home=$HOME/.oraid/validator1
oraid genesis validate-genesis --home=$HOME/.oraid/validator1

# update staking genesis
update_genesis '.app_state["staking"]["params"]["unbonding_time"]="240s"'

# update crisis variable to orai
update_genesis '.app_state["crisis"]["constant_fee"]["denom"]="orai"'

# udpate gov genesis
update_genesis '.app_state["gov"]["voting_params"]["voting_period"]="15s"'
update_genesis '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="orai"'

# update mint genesis
update_genesis '.app_state["mint"]["params"]["mint_denom"]="orai"'

# port key (validator1 uses default ports)
# validator1 1317, 9090, 9091, 26658, 26657, 26656, 6060

# change app.toml values
VALIDATOR1_APP_TOML=$HOME/.oraid/validator1/config/app.toml

# change config.toml values
VALIDATOR1_CONFIG=$HOME/.oraid/validator1/config/config.toml

# Pruning - comment this configuration if you want to run upgrade script
pruning="custom"
pruning_keep_recent="5"
pruning_keep_every="10"
pruning_interval="10000"

sed -i -e "s%^pruning *=.*%pruning = \"$pruning\"%; " $VALIDATOR1_APP_TOML
sed -i -e "s%^pruning-keep-recent *=.*%pruning-keep-recent = \"$pruning_keep_recent\"%; " $VALIDATOR1_APP_TOML
sed -i -e "s%^pruning-keep-every *=.*%pruning-keep-every = \"$pruning_keep_every\"%; " $VALIDATOR1_APP_TOML
sed -i -e "s%^pruning-interval *=.*%pruning-interval = \"$pruning_interval\"%; " $VALIDATOR1_APP_TOML

# state sync  - comment this configuration if you want to run upgrade script
snapshot_interval="10"
snapshot_keep_recent="2"

sed -i -e "s%^snapshot-interval *=.*%snapshot-interval = \"$snapshot_interval\"%; " $VALIDATOR1_APP_TOML
sed -i -e "s%^snapshot-keep-recent *=.*%snapshot-keep-recent = \"$snapshot_keep_recent\"%; " $VALIDATOR1_APP_TOML

# validator1
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR1_CONFIG

# start all three validators
BACKGROUND=${BACKGROUND:-""}
# if the $BACKGROUND env var is empty, then we run foreground mode
if [ -z "$BACKGROUND" ]
then
    oraid start --home=$HOME/.oraid/validator1 --minimum-gas-prices=0orai
else
    screen -S validator1 -d -m oraid start --home=$HOME/.oraid/validator1 --minimum-gas-prices=0orai
    echo "Validator 1 are up and running!"
fi

