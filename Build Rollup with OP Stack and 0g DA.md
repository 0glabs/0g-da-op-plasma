# Introduction

[OP Stack](https://stack.optimism.io/) is the set of [software components](https://github.com/ethereum-optimism/optimism) that run the [Optimism](https://www.optimism.io/) rollup and can be deployed independently to power third-party rollups.

By default, OP Stack sequencers write batches to Ethereum in the form of calldata or 4844 blobs to commit to the transactions included in the canonical L2 chain. In Alt-DA mode, OP Stack sequencers and full nodes are configured talk to a third-party HTTP server for writing and reading tx batches to and from DA. Optimism's Alt-DA spec contains a more in-depth breakdown of how this works.

To implement this server spec, 0G DA provides the da-server, which can be run alongside OP Stack sequencers and full nodes to securely communicate with the 0G DA Client.

# Prerequisites

Ensure you have installed the following software.

| Software | Version    |
| -------- | ---------- |
| Git      | 2.34.1     |
| Go       | 1.23.3     |
| Node     | ^20        |
| just     | 1.34.0     |
| Make     | 4.3        |
| jq       | 1.6        |
| direnv   | 2.25.2     |
| foundry  | 0.2.0      |

Use the following releases while following the guide:
* op-node/v1.9.1
* op-proposer/v1.9.1
* op-batcher/v1.9.1
* op-geth v1.101408.0

# Setup 0g DA-server

## Run 0g-da server (manually)

1. Navigate to

```bash
git clone https://github.com/0glabs/0g-da-op-plasma.git
cd 0g-da-op-plasma
```

2. Build 0g DA Server:

```bash
make da-server
```

3. Run your DA sidecar:

```bash
./bin/da-server --addr 127.0.0.1 --port 3100 --zg.server <rpc to an da client>
```

## Run using docker

1. Build the Docker image

```bash
docker build -t 0g-da-op-plasma .
```

2. Run the Docker container

```bash
docker run -p 3100:3100 0g-da-op-plasma:latest da-server --addr 0.0.0.0 --port 3100 --zg.server <rpc to an da client>
```

# Build the Optimism Monorepo

1. Clone and navigate to the Optimism Monorepo:

```bash
git clone https://github.com/ethereum-optimism/optimism.git
cd optimism
git fetch --tag --all
git checkout v1.9.1
git submodule update --init --recursive
```

2. Check your dependencies
```bash
# Install Git
sudo apt install -y git curl make jq
 
# Install Go
wget https://go.dev/dl/go1.23.3 linux-amd64.tar.gz
tar xvzf go1.23.3 linux-amd64.tar.gz
sudo cp go/bin/go /usr/bin/go
sudo mv go /usr/lib
echo export GOROOT=/usr/lib/go >> ~/.bashrc
 
# Install Node.js
curl -sL https://deb.nodesource.com/setup_14.x | sudo -E bash -
sudo apt-get install -y nodejs
 
# Install Pnpm
sudo npm install -g pnpm
 
# Install Make
sudo apt install -y make
 
# Install jq
sudo apt install -y jq
 
# Install direnv
sudo apt install -y direnv

# install just
sudo apt install cargo
cargo install just

# install semgrep
sudo apt install snapd
sudo snap install semgrep
sudo apt install python3-pip
python3 -m pip install --upgrade pip
python3 -m pip install semgrep==1.90.0
sudo apt-get install python3-virtualenv
python3 -m virtualenv myenv
source myenv/bin/activate
python3 -m pip install --upgrade semgrep

# install foundry
cargo install --git https://github.com/foundry-rs/foundry --profile release --locked forge cast chisel anvil
```bash

Run the following script and double check that you have all of the required versions installed. If you don't have the correct versions installed, you may run into unexpected errors.
```bash
./packages/contracts-bedrock/scripts/getting-started/versions.sh
```

3. Compile the necessary packages:

```bash
make op-node op-batcher op-proposer
make build
```

# Build the Optimism Geth Source

1. Clone and navigate to op-geth:

```bash
cd
git clone https://github.com/ethereum-optimism/op-geth.git
cd op-geth
git fetch --tag --all
git checkout v1.101408.0
```

2. Compile op-geth:

```bash
make geth
```

# Get Access to a Sepolia Node

For deploying to Sepolia, access an L1 node using a provider like [Alchemy](https://www.alchemy.com/) (easier) or run your own Sepolia node (harder).

# Fill out environment variables

You'll need to fill out a few environment variables before you can start deploying your chain.


1. Enter the Optimism Monorepo
```bash
cd ~/optimism
```

2. Duplicate the sample environment variable file
```bash
cp .envrc.example .envrc
```

3. Fill out the environment variable file
Open up the environment variable file and fill out the following variables:

```
| Variable Name | Description                                                                                                                                                                                                  |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `L1_RPC_URL`  | URL for your L1 node (a Sepolia node in this case).                                                                                                                                                          |
| `L1_RPC_KIND` | Kind of L1 RPC you're connecting to, used to inform optimal transactions receipts fetching. Valid options: `alchemy`, `quicknode`, `infura`, `parity`, `nethermind`, `debug_geth`, `erigon`, `basic`, `any`. |
```

# Generate addresses

You'll need four addresses and their private keys when setting up the chain:

*   The `Admin` address has the ability to upgrade contracts.
*   The `Batcher` address publishes Sequencer transaction data to L1.
*   The `Proposer` address publishes L2 transaction results (state roots) to L1.
*   The `Sequencer` address signs blocks on the p2p network.

You can use `cast wallet` in the `contracts-bedrock` package for key generation:

1. In the Optimism repo, navigate to the [contracts-bedrock package](https://github.com/ethereum-optimism/optimism/tree/129032f15b76b0d2a940443a39433de931a97a44/packages/contracts-bedrock):

```bash
cd ~/optimism/packages/contracts-bedrock
```

2. Generate accounts:

```bash
./scripts/getting-started/wallets.sh
```

You should see an output similar to:

```text
Copy the following into your .envrc file:
  
# Admin address
export GS_ADMIN_ADDRESS=0x9625B9aF7C42b4Ab7f2C437dbc4ee749d52E19FC
export GS_ADMIN_PRIVATE_KEY=0xbb93a75f64c57c6f464fd259ea37c2d4694110df57b2e293db8226a502b30a34

# Batcher address
export GS_BATCHER_ADDRESS=0xa1AEF4C07AB21E39c37F05466b872094edcf9cB1
export GS_BATCHER_PRIVATE_KEY=0xe4d9cd91a3e53853b7ea0dad275efdb5173666720b1100866fb2d89757ca9c5a
  
# Proposer address
export GS_PROPOSER_ADDRESS=0x40E805e252D0Ee3D587b68736544dEfB419F351b
export GS_PROPOSER_PRIVATE_KEY=0x2d1f265683ebe37d960c67df03a378f79a7859038c6d634a61e40776d561f8a2
  
# Sequencer address
export GS_SEQUENCER_ADDRESS=0xC06566E8Ec6cF81B4B26376880dB620d83d50Dfb
export GS_SEQUENCER_PRIVATE_KEY=0x2a0290473f3838dbd083a5e17783e3cc33c905539c0121f9c76614dda8a38dca
```

Record and securely store these key details. You'll need to fund `Admin`, `Proposer`, and `Batcher` with Sepolia ETH (0.5 ETH for `Admin`, 0.2 ETH for `Proposer`, 0.1 ETH for `Batcher`).

> **NOTE FOR PRODUCTION:**  
> Use secure hardware for key management in production environments. cast wallet is not designed for production deployments.

3. Save the addresses

Copy the output from the previous step and paste it into your .envrc file as directed.

# Load environment variables

Now that you've filled out the environment variable file, you need to load those variables into your terminal.

1. Enter the Optimism Monorepo

```bash
cd ~/optimism
```

2. Load the variables with direnv

Next you'll need to allow direnv to read this file and load the variables into your terminal using the following command.

```bash
eval "$(direnv hook bash)"
direnv allow
```

3. Confirm that the variables were loaded

After running direnv allow you should see output that looks something like the following (the exact output will vary depending on the variables you've set, don't worry if it doesn't look exactly like this):

```bash
direnv: loading ~/optimism/.envrc                                                            
direnv: export +DEPLOYMENT_CONTEXT +ETHERSCAN_API_KEY +GS_ADMIN_ADDRESS +GS_ADMIN_PRIVATE_KEY +GS_BATCHER_ADDRESS +GS_BATCHER_PRIVATE_KEY +GS_PROPOSER_ADDRESS +GS_PROPOSER_PRIVATE_KEY +GS_SEQUENCER_ADDRESS +GS_SEQUENCER_PRIVATE_KEY +IMPL_SALT +L1_RPC_KIND +L1_RPC_URL +PRIVATE_KEY +TENDERLY_PROJECT +TENDERLY_USERNAME
```

# Core Contract Deployment

Deploy essential L1 contracts for the chain’s functionality:

1. Update `/optimism/packages/contracts-bedrock/deploy-config` and update file `getting_started.json.`

```bash
cd packages/contracts-bedrock
./scripts/getting-started/config.sh
```

2. Add the following at the bottom of the config generated:

```
  "useAltDA": true,
  "daCommitmentType": "GenericCommitment",
  "daChallengeWindow": 160,
  "daResolveWindow": 160,
  "daBondSize": 1000000,
  "daResolverRefundPercentage": 0
```

Example config ( for reference purpose ):

```
{

 "l1StartingBlockTag": "0x0b2b81474a22fc1122bbb3a465985c5cb40dfd8ef18bfe4fc0a9fa47e775d692",

  "l1ChainID": 11155111,
  "l2ChainID": 42069,
  "l2BlockTime": 2,
  "l1BlockTime": 12,

  "maxSequencerDrift": 600,
  "sequencerWindowSize": 3600,
  "channelTimeout": 300,

  "p2pSequencerAddress": "0xd34514056DBE102dF5c24DAf3e78701b502F34c7",
  "batchInboxAddress": "0xff00000000000000000000000000000000042069",
  "batchSenderAddress": "0x4eD53FeB5b06c60368490683e2b317d7C20eC41E",

  "l2OutputOracleSubmissionInterval": 120,
  "l2OutputOracleStartingBlockNumber": 0,
  "l2OutputOracleStartingTimestamp": 1729571268,

  "l2OutputOracleProposer": "0x7d32557e4e79F494836037aD622AaB60125D8323",
  "l2OutputOracleChallenger": "0x6A06226C406f7298E1F7CF3F1aa72f3Bc5eF4247",

  "finalizationPeriodSeconds": 12,

  "proxyAdminOwner": "0x6A06226C406f7298E1F7CF3F1aa72f3Bc5eF4247",
  "baseFeeVaultRecipient": "0x6A06226C406f7298E1F7CF3F1aa72f3Bc5eF4247",
  "l1FeeVaultRecipient": "0x6A06226C406f7298E1F7CF3F1aa72f3Bc5eF4247",
  "sequencerFeeVaultRecipient": "0x6A06226C406f7298E1F7CF3F1aa72f3Bc5eF4247",
  "finalSystemOwner": "0x6A06226C406f7298E1F7CF3F1aa72f3Bc5eF4247",
  "superchainConfigGuardian": "0x6A06226C406f7298E1F7CF3F1aa72f3Bc5eF4247",

  "baseFeeVaultMinimumWithdrawalAmount": "0x8ac7230489e80000",
  "l1FeeVaultMinimumWithdrawalAmount": "0x8ac7230489e80000",
  "sequencerFeeVaultMinimumWithdrawalAmount": "0x8ac7230489e80000",
  "baseFeeVaultWithdrawalNetwork": 0,
  "l1FeeVaultWithdrawalNetwork": 0,
  "sequencerFeeVaultWithdrawalNetwork": 0,

  "gasPriceOracleOverhead": 0,
  "gasPriceOracleScalar": 1000000,

  "enableGovernance": true,
  "governanceTokenSymbol": "OP",
  "governanceTokenName": "Optimism",
  "governanceTokenOwner": "0x6A06226C406f7298E1F7CF3F1aa72f3Bc5eF4247",

  "l2GenesisBlockGasLimit": "0x1c9c380",
  "l2GenesisBlockBaseFeePerGas": "0x3b9aca00",

  "eip1559Denominator": 50,
  "eip1559DenominatorCanyon": 250,
  "eip1559Elasticity": 6,
  "l2GenesisFjordTimeOffset": "0x0",
  "l2GenesisRegolithTimeOffset": "0x0",
  "l2GenesisEcotoneTimeOffset": "0x0",
  "l2GenesisDeltaTimeOffset": "0x0",
  "l2GenesisCanyonTimeOffset": "0x0",
  "systemConfigStartBlock": 0,

  "requiredProtocolVersion": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "recommendedProtocolVersion": "0x0000000000000000000000000000000000000000000000000000000000000000",

  "faultGameAbsolutePrestate": "0x03c7ae758795765c6664a5d39bf63841c71ff191e9189522bad8ebff5d4eca98",
  "faultGameMaxDepth": 44,
  "faultGameClockExtension": 0,
  "faultGameMaxClockDuration": 1200,
  "faultGameGenesisBlock": 0,
  "faultGameGenesisOutputRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "faultGameSplitDepth": 14,
  "faultGameWithdrawalDelay": 600,

  "preimageOracleMinProposalSize": 1800000,
  "preimageOracleChallengePeriod": 300,

  "useAltDA": true,
  "daCommitmentType": "GenericCommitment",
  "daChallengeWindow": 160,
  "daResolveWindow": 160,
  "daBondSize": 1000000,
  "daResolverRefundPercentage": 0
}

```

2. Navigate to `/optimism/packages/contracts-bedrock/` and the deploy contracts (this can take up to 15 minutes):

```bash
DEPLOYMENT_OUTFILE=deployments/artifact.json \
DEPLOY_CONFIG_PATH=deploy-config/getting-started.json \
forge script scripts/deploy/Deploy.s.sol:Deploy  --broadcast --private-key $GS_ADMIN_PRIVATE_KEY \
--rpc-url $L1_RPC_URL --slow
```

3. L2 Allocs

```bash
CONTRACT_ADDRESSES_PATH=deployments/artifact.json DEPLOY_CONFIG_PATH=deploy-config/getting-started.json STATE_DUMP_PATH=deploy-config/statedump.json forge script scripts/L2Genesis.s.sol:L2Genesis --sig 'runWithStateDump()' --chain 42069
```

# Setting Up L2 Configuration

After configuring the L1 layer, focus shifts to establishing the L2 infrastructure. This involves generating three key files:

- genesis.json for the genesis block
- rollup.json for rollup configurations
- jwt.txt for secure communication between op-node and op-geth

1. Navigate to the op-node directory:

```bash
cd ~/optimism/op-node
```

2. Run the following command, ensuring you replace <RPC> with your specific L1 RPC URL. This generates the genesis.json and rollup.json files:

```bash
go run cmd/main.go genesis l2 \
--deploy-config ../packages/contracts-bedrock/deploy-config/getting-started.json \
--l1-deployments ../packages/contracts-bedrock/deployments/artifact.json \
--outfile.l2 genesis.json \
--outfile.rollup rollup.json \
--l1-rpc $L1_RPC_URL \
--l2-allocs ../packages/contracts-bedrock/deploy-config/statedump.json
```

You'll find the newly created genesis.json and rollup.json in the op-node package.

3. Add the following at the end of rollup.json:

```
 "alt_da": {
    "da_challenge_contract_address": "0x0000000000000000000000000000000000000000",
    "da_commitment_type": "GenericCommitment",
    "da_challenge_window": 160,
    "da_resolve_window": 160
  }
```

4. Generate a `jwt.txt` file, which is crucial for the secure interaction between nodes:

```bash
openssl rand -hex 32 > jwt.txt
```

5. To get op-geth ready, move the genesis.json and jwt.txt files into its directory:

```bash
cp genesis.json ~/op-geth
cp jwt.txt ~/op-geth
```

# Initialize and Configure Geth

Prepare `op-geth` for running the chain:

1. Navigate to op-geth:

```bash
cd ~/op-geth
```

2. Create a data directory:

```bash
mkdir datadir
```

3. Initialize with the genesis file:

```bash
build/bin/geth init --datadir=datadir genesis.json
```

### Running op-geth

To initiate `op-geth`, navigate to its directory and execute the following commands:

```bash
cd ~/op-geth
./build/bin/geth \
  --datadir ./datadir \
  --http \
  --http.corsdomain="*" \
  --http.vhosts="*" \
  --http.addr=0.0.0.0 \
  --http.port=9545 \
  --http.api=web3,debug,eth,txpool,net,engine \
  --ws \
  --ws.addr=0.0.0.0 \
  --ws.port=9546 \
  --ws.origins="*" \
  --ws.api=debug,eth,txpool,net,engine \
  --syncmode=full \
  --nodiscover \
  --maxpeers=0 \
  --networkid=42069 \
  --authrpc.vhosts="*" \
  --authrpc.addr=0.0.0.0 \
  --authrpc.port=9551 \
  --authrpc.jwtsecret=./jwt.txt \
  --rollup.disabletxpoolgossip=true \
  --state.scheme=path
```

`op-geth` is now active, but block creation will begin once `op-node` is operational.


### Running op-node

To launch op-node, which acts as a consensus client, run:

```bash
cd ~/optimism/op-node
./bin/op-node \
  --l2=http://localhost:9551 \
  --l2.jwt-secret=./jwt.txt \
  --sequencer.enabled \
  --sequencer.l1-confs=5 \
  --verifier.l1-confs=4 \
  --rollup.config=./rollup.json \
  --rpc.addr=0.0.0.0 \
  --rpc.port=8547 \
  --p2p.disable \
  --rpc.enable-admin \
  --p2p.sequencer.key=$GS_SEQUENCER_PRIVATE_KEY \
  --l1=$L1_RPC_URL \
  --l1.rpckind=$L1_RPC_KIND \
  --altda.enabled=true \
  --altda.da-server=<DA_SERVER_HTTP_URL> \
  --altda.da-service=true \
  --l1.beacon.ignore=true
```

Block creation will commence once op-node starts processing L1 information and interfaces with op-geth.

### P2P Synchronization

To optimize synchronization and avoid network resource waste:

- Disable p2p sync (--p2p.disable) by default.
- Use specific command line parameters for synchronization among multiple nodes.


# Running op-batcher

`op-batcher` is crucial in publishing transactions from the Sequencer to L1. Ensure it has at least 1 Sepolia ETH for operational continuity.

```bash
cd ~/optimism/op-batcher
./bin/op-batcher \
  --l2-eth-rpc=http://localhost:9545 \
  --rollup-rpc=http://localhost:8547 \
  --poll-interval=1s \
  --sub-safety-margin=6 \
  --num-confirmations=1 \
  --safe-abort-nonce-too-low-count=3 \
  --resubmission-timeout=30s \
  --rpc.addr=0.0.0.0 \
  --rpc.port=8548 \
  --rpc.enable-admin \
  --max-channel-duration=1 \
  --l1-eth-rpc=$L1_RPC_URL \
  --private-key=$GS_BATCHER_PRIVATE_KEY \
  --altda.enabled=true \
  --altda.da-service=true \
  --altda.da-server=<DA_SERVER_HTTP_URL>
```

### Controlling Batcher Costs

The --max-channel-duration=n setting tells the batcher to write all the data to L1 every n L1 blocks. When it is low, transactions are written to L1 frequently and other nodes can synchronize from L1 quickly. When it is high, transactions are written to L1 less frequently and the batcher spends less ETH. If you want to reduce costs, either set this value to 0 to disable it or increase it to a higher value.

# Running op-proposer

Finally, start `op-proposer` to propose new state roots:

```bash
cd ~/optimism/op-proposer
./bin/op-proposer \
  --poll-interval=12s \
  --rpc.port=9560 \
  --rollup-rpc=http://localhost:8547 \
  --l2oo-address=$L2OO_ADDR \
  --private-key=$PROPOSER_KEY \
  --l1-eth-rpc=$L1_RPC_URL
```

# Acquire Sepolia ETH for Layer 2

To obtain ETH on your Rollup:

1. Go to `contracts-bedrock`:

```bash
cd ~/optimism/packages/contracts-bedrock
```

2. Find the L1 standard bridge contract address:

```bash
cat deployments/artifact.json | jq -r .L1StandardBridgeProxy
```

3. Send Sepolia ETH to the bridge contract address.

# Conduct Test Transactions

You now have a fully operational 0gDA-Powered Optimism-based EVM Rollup. Experiment with it as you would with any other test blockchain.

**Congratulations on setting up your chain!**