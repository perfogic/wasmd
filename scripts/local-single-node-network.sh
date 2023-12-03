#!/bin/bash
set -ux

# always returns true so set -e doesn't exit if it is not running.
# set install jq for ubuntu & mac. If failed we still can continue running
apt install jq
brew install jq
rm -rf $HOME/.oraid/
killall screen

NODE_HOME=${NODE_HOME:-"$HOME/.oraid/validator1"}

# make four orai directories
mkdir $HOME/.oraid
mkdir $NODE_HOME

# init all three validators
oraid init --chain-id=testing validator1 --home=$NODE_HOME

# create keys for all three validators
(cat .env) | oraid keys add validator1 --recover --keyring-backend test --home=$NODE_HOME

update_genesis () {    
    cat $NODE_HOME/config/genesis.json | jq "$1" > $NODE_HOME/config/tmp_genesis.json && mv $NODE_HOME/config/tmp_genesis.json $NODE_HOME/config/genesis.json
}

# change staking denom to orai
update_genesis '.app_state["staking"]["params"]["bond_denom"]="orai"'

# create validator node 1
oraid genesis add-genesis-account $(oraid keys show validator1 -a --keyring-backend=test --home=$NODE_HOME | grep orai) 1000000000000orai,1000000000000stake --home=$NODE_HOME
oraid genesis gentx validator1 500000000orai --keyring-backend=test --home=$NODE_HOME --chain-id=testing
oraid genesis collect-gentxs --home=$NODE_HOME
oraid genesis validate-genesis --home=$NODE_HOME

# update staking genesis
update_genesis '.app_state["staking"]["params"]["unbonding_time"]="1209600s"'

# update crisis variable to orai
update_genesis '.app_state["crisis"]["constant_fee"]["denom"]="orai"'

# udpate gov genesis
update_genesis '.app_state["gov"]["voting_params"]["voting_period"]="60s"'
update_genesis '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="orai"'

# update mint genesis
update_genesis '.app_state["mint"]["params"]["mint_denom"]="orai"'

# port key (validator1 uses default ports)
# validator1 1317, 9090, 9091, 26658, 26657, 26656, 6060

# change app.toml values
VALIDATOR1_APP_TOML=$NODE_HOME/config/app.toml

# change config.toml values
VALIDATOR1_CONFIG=$NODE_HOME/config/config.toml

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
    oraid start --home=$NODE_HOME --minimum-gas-prices=0orai --rpc.laddr tcp://0.0.0.0:26657 --grpc.address 0.0.0.0:9090
else
    screen -S validator1 -d -m oraid start --home=$NODE_HOME --minimum-gas-prices=0orai --rpc.laddr tcp://0.0.0.0:26657 --grpc.address 0.0.0.0:9090
    echo "Validator 1 are up and running!"
fi

