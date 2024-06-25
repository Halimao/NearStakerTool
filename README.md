# NearStakeTool

A tool that helps users to view current stake statistics of the near chain

## Install

### Build from source

```shell
git clone git@github.com:Halimao/NearStakerTool.git
cd NearStakerTool
go build
```

## Quickstart

Running following command will export every delegator stake statistics of different near validator to stake.db

```shell
./NearStakerTool
```

You can use `sqlite3` to view the data by `sqlite3 stake.db`