# did-cell-indexer

did cell indexing service, which you can use to query the did cell assets in the chain, as well as to structure
transactions.

* [Prerequisites](#prerequisites)
* [Install &amp; Run](#install--run)
    * [Source Compile](#source-compile)
    * [Docker](#docker)
* [APIs](#apis)

## Prerequisites

* Ubuntu 18.04 or newer
* MYSQL >= 8.0
* Redis >= 5.0 (for cache)
* GO version >= 1.17.10
* [ckb-node](https://github.com/nervosnetwork/ckb) (Must be synced to latest height and add `Indexer` module to
  ckb.toml)
* If the version of the dependency package is too low, please install `gcc-multilib` (apt install gcc-multilib)
* Machine configuration: 4c8g500G

## Install & Run

### Source Compile

```bash
# get the code
git clone https://github.com/dotbitHQ/did-cell-indexer.git

# compile
cd did-cell-indexer
make svr

# edit config/config.yaml for svr
vim config/config.yaml

# run
make svr
./did_indexer_svr --config=config/config.yaml
```

### Docker

* docker >= 20.10
* docker-compose >= 2.2.2

```bash
# install docker compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.2.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose

# get the code
git clone https://github.com/dotbitHQ/did-cell-indexer.git
cd did-cell-indexer

# edit config/config.yaml for svr
vim config/config.yaml

# run
docker-compose up -d
```

_if you already have mysql,redis installed, just run_

```bash
docker run -dp 9132:9132 -v $PWD/config:/app/config -v $PWD/logs:/app/logs --name did-indexer-svr admindid/did-indexer-svr:latest
```

## APIs
More APIs see [API.md](https://github.com/dotbitHQ/did-cell-indexer/blob/main/API.md)