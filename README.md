# 诚泰融租链

## 概览

TODO

## 区块链节点

### 编译

操作系统：任意的 Linux 发行版。

依赖软件包：make、git、Go 1.14 或更高版本、UPX（可选）。

```bash
make build
```

编译后，在 build 目录下，会出现以下可执行文件。

- `chengtay-chain` 区块链节点程序。

建议将该文件移动至目标机器的 `/usr/local/bin` 目录下。

#### 中国大陆用户

在编译前，需要设置 Go 代理。推荐使用由七牛云运行的 [goproxy.cn](https://goproxy.cn/)。

```bash
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

#### UPX 压缩

可以使用 UPX 程序将生成的可执行文件进行压缩，通常可将体积压缩至原先的一半。

```bash
make build-upx
```

#### 在离线环境下编译

保险起见，可将依赖的源代码全部保存至本地的 vendor 目录。在极端环境下，代码依然可以成功编译。

```bash
go mod vendor
```

离线编译时，须修改 `Makefile` 文件，将

`BUILD_FLAGS = -mod=readonly -ldflags "$(LD_FLAGS)"`

改为

`BUILD_FLAGS = -mod=vendor -ldflags "$(LD_FLAGS)"`。

#### 交叉编译

通过环境变量 `$GOOS` 和 `$GOARCH` 指定目标平台。例如，编译到 arm64 平台的 Linux 的命令如下。

```bash
GOOS=linux GOARCH=arm64 make build
```

#### 区块链节点要求

`chengtay-chain` 是静态编译的程序，对区块链节点的目标操作系统没有要求，可以是任意的 Linux 发行版，无任何依赖组件。可以用 file 命令验证，结果大致如下。

```
chengtay-chain: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, Go BuildID=mz73U2-osJg4YY28a1Cg/7VjC64v5DrESnXX6ANJ3/T0TXzzUR_0PCGhyZFqwy/Yci4hmN0ziL57Takz5CI, stripped
```
```
chengtay-chain: ELF 64-bit LSB executable, ARM aarch64, version 1 (SYSV), statically linked, Go BuildID=_Q5vQ0Fw0aD3ZCHajbXC/1IgI5zUugKk0upLQjY1p/pcI7jl5IusOIVP-rgfez/GRLfwEDtadWsg0QYYYBo, stripped
```

### 创建并分发密钥和创世区块

#### 创建

在部署区块链节点之前，需要事先确定好区块链验证者节点的个数，然后在一台机器上创建各节点的密钥和创世区块数据。

在以下示例中，假设验证者节点个数为 5。注：如需单机调试，请将验证者节点个数设为 1。假设区块链主程序位于 `/usr/local/bin/chengtay-chain`。

确保 `$HOME/.chengtaychain` 目录**不存在**。

注意，`~` 仅**在 Shell 下**代指当前用户的 `$HOME` 目录，编程时请勿使用 `~` 代指 `$HOME`。如果不小心生成了名为 `~` 的文件夹，删除时请一定**不要**使用 `rm -rf ~` 命令，而是使用`rm -rf ./~` 或 `rm -rf \~` 命令。

执行以下命令，生成 5 个节点的密钥和对应的创世区块。
```bash
chengtay-chain gen_genesis 5
```

`$HOME/.chengtaychain` 目录下产生了以下文件。

- `config/genesis.json` 创世区块文件。此文件应分发给所有的区块链节点。
- `config/config.toml` 节点配置文件。此文件既可分发给其他节点，也可以让节点自动生成。
- `config/priv_validator_key.json.$NUM.json` 编号为`$NUM`的验证者节点的公钥和私钥。此文件应保密，仅由该节点保存。
- `data/priv_validator_state.json.$NUM.json` 编号为`$NUM`的验证者节点的附加信息。此文件应保密，仅由该节点保存。

#### 分发

- 将区块链节点程序 `chengtay-chain` 和创世区块文件 `config/genesis.json` 分发给所有的区块链节点。
- 将 `config/priv_validator_key.json.$NUM.json` 和 `data/priv_validator_state.json.$NUM.json` 分发给编号为 `$NUM`的验证者节点。

### 部署节点

1. 将区块链节点程序 `chengtay-chain` 移动到 `/usr/local/bin/chengtay-chain`，并赋予可执行权限。
2. 将创世区块文件 `genesis.json` 移动到 `$HOME/.chengtaychain/config/genesis.json`。
3. 对于编号为 `$NUM`的验证者节点，将文件 `priv_validator_key.json.$NUM.json` 移动到 `$HOME/.chengtaychain/config/priv_validator_key.json`，将文件 `priv_validator_state.json.$NUM.json` 移动到 `$HOME/.chengtaychain/data/priv_validator_state.json`。
4. 初始化节点。
```bash
chengtay-chain init
```
注意，如果在跳过第 3 步的情况下初始化节点，则节点将自动生成一对新的密钥。由于该密钥不在创世区块中，节点将作为非验证者节点加入区块链网络。

### 运行节点

```bash
chengtay-chain node
```

### 节点的连接、发现

TODO

### 示例：Docker 方案

见 [这里](https://github.com/ChengtayChain/ChengtayChain/tree/master/DOCKER)。该方案可一键运行多个节点组成的区块链网络。

## 数据服务器

数据服务器运行在权重最大的验证节点上，使用该节点的公私钥对。

对于区块链而言，数据服务器是一个客户端。

简单的客户端示例见 [这里](https://github.com/ChengtayChain/ChengtayChain/tree/master/chengtay/cmd/example-client)。该示例不断创建随机的交易，并向区块链节点发送这些交易。

TODO

## Web 界面

TODO

## 数据结构

TODO
