package proxy

import (
	chengtay "github.com/ChengtayChain/ChengtayChain/chengtay/abci"
	"sync"

	"github.com/pkg/errors"

	abcicli "github.com/ChengtayChain/ChengtayChain/abci/client"
	"github.com/ChengtayChain/ChengtayChain/abci/types"
)

// NewABCIClient returns newly connected client
type ClientCreator interface {
	NewABCIClient() (abcicli.Client, error)
}

//----------------------------------------------------
// local proxy uses a mutex on an in-proc app

type localClientCreator struct {
	mtx *sync.Mutex
	app types.Application
}

func NewLocalClientCreator(app types.Application) ClientCreator {
	return &localClientCreator{
		mtx: new(sync.Mutex),
		app: app,
	}
}

func (l *localClientCreator) NewABCIClient() (abcicli.Client, error) {
	return abcicli.NewLocalClient(l.mtx, l.app), nil
}

//---------------------------------------------------------------
// remote proxy opens new connections to an external app process

type remoteClientCreator struct {
	addr        string
	transport   string
	mustConnect bool
}

func NewRemoteClientCreator(addr, transport string, mustConnect bool) ClientCreator {
	return &remoteClientCreator{
		addr:        addr,
		transport:   transport,
		mustConnect: mustConnect,
	}
}

func (r *remoteClientCreator) NewABCIClient() (abcicli.Client, error) {
	remoteApp, err := abcicli.NewClient(r.addr, r.transport, r.mustConnect)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to proxy")
	}
	return remoteApp, nil
}

//-----------------------------------------------------------------
// default

func DefaultClientCreator(addr, transport, dbDir string) ClientCreator {
	// Set ChengtayChain to default.
	return NewLocalClientCreator(chengtay.NewApplication(dbDir))

	//switch addr {
	//case "kvstore":
	//	return NewLocalClientCreator(kvstore.NewApplication())
	//case "persistent_kvstore":
	//	return NewLocalClientCreator(kvstore.NewPersistentKVStoreApplication(dbDir))
	//case "noop":
	//	return NewLocalClientCreator(types.NewBaseApplication())
	//default:
	//	// 非 go 的版本就是通过 socket 或者 grpc 进行通信。也就是最后的默认分支创建的 abci 实例
	//	mustConnect := false // loop retrying
	//	return NewRemoteClientCreator(addr, transport, mustConnect)
	//}
}
