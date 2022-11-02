package convert

import (
	"github.com/urfave/cli/v2"
)

// // IRegistry _
// type IRegistry interface {
// 	models.IRegistry
// 	models.IContracterRegistry
// }

// type IDealer interface {
// 	models.IWorkerDealer
// 	models.IContracterDealer
// }

// DistributedConvertApp _
type DistributedConvertApp struct {
}

// Init _
func (a *DistributedConvertApp) Init() error {
	return nil
}

// InitLocal _
func (a *DistributedConvertApp) InitLocal() error {
	return nil
}

// StartContracter _
func (a *DistributedConvertApp) StartContracter() error {
	return nil
}

// StartAPI _
func (a *DistributedConvertApp) StartAPI() error {
	return nil
}

// AddTask _
func (a *DistributedConvertApp) AddTask(c *cli.Context) error {
	return nil
}

// StartWorker _
func (a *DistributedConvertApp) StartWorker() error {
	return nil
}

// StartRemoteWorker _
func (a *DistributedConvertApp) StartRemoteWorker() error {
	return nil
}

// ListOrders _
func (a *DistributedConvertApp) ListOrders(cliCtx *cli.Context) (string, error) {
	return "", nil
}

// ShowOrder _
func (a *DistributedConvertApp) ShowOrder(orderID string) (string, error) {
	return "", nil
}

// ListSegments _
func (a *DistributedConvertApp) ListSegments(cliCtx *cli.Context, orderID string) (string, error) {
	return "", nil
}

// CancelOrder _
func (a *DistributedConvertApp) CancelOrder(orderID string) error {
	return nil
}

// Wait _
func (a *DistributedConvertApp) Wait() <-chan struct{} {
	return nil
}
