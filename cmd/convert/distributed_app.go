package convert

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/adapters"
	"github.com/wailorman/fftb/pkg/distributed/handlers"
	"github.com/wailorman/fftb/pkg/distributed/local"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/registry"
	"github.com/wailorman/fftb/pkg/distributed/remote"
	"github.com/wailorman/fftb/pkg/distributed/ukvs/ubolt"
	"github.com/wailorman/fftb/pkg/distributed/worker"
	"github.com/wailorman/fftb/pkg/files"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
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
	storage              models.IStorageController
	publisher            models.IAuthor
	registry             models.IRegistry
	dealer               models.IDealer
	contracter           *local.ContracterInstance
	contracterInteractor *adapters.ContracterAdapter
	workerInstance       *worker.Instance
	storageClient        models.IStorageClient
	ctx                  context.Context
	cancel               func()
	logger               *logrus.Entry
	closed               chan struct{}
	wg                   *chwg.ChannelledWaitGroup
}

// Init _
func (a *DistributedConvertApp) Init() error {
	var err error

	a.closed = make(chan struct{})
	a.wg = chwg.New()

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.Formatter = &prefixed.TextFormatter{
		FullTimestamp: true,
	}
	a.logger = logger.WithField("prefix", "fftb")

	a.ctx, a.cancel = context.WithCancel(context.WithValue(context.Background(), ctxlog.LoggerContextKey, a.logger))

	a.storage, err = initStorage(a.ctx)

	if err != nil {
		return errors.Wrap(err, "Initializing storage")
	}

	a.storageClient, err = initStorageClient()

	if err != nil {
		return errors.Wrap(err, "Initializing storage client")
	}

	a.registry, err = initRegistry(a.ctx)

	if err != nil {
		return errors.Wrap(err, "Initializing registry")
	}

	a.dealer, err = local.NewDealer(a.ctx, a.storage, a.registry, models.NewSegmentMutation())

	if err != nil {
		return errors.Wrap(err, "Initializing delaer")
	}

	a.publisher = &models.Author{Name: "local"}

	contracterTmpPath := files.NewPath(".fftb/contracter")
	err = contracterTmpPath.Create()

	if err != nil {
		return errors.Wrap(err, "Creating tmp path for contracter")
	}

	a.contracter, err = local.NewContracter(
		a.ctx,
		a.dealer,
		a.registry,
		a.storageClient,
		models.NewOrderMutation(logger), contracterTmpPath)

	if err != nil {
		return errors.Wrap(err, "Initializing contracter")
	}

	a.contracterInteractor = adapters.NewContracterAdapter(a.contracter)

	return nil
}

// StartContracter _
func (a *DistributedConvertApp) StartContracter() error {
	publishWorker := local.NewContracterPublishWorker(a.ctx, a.contracter, models.NewOrderMutation(a.logger))
	concatWorker := local.NewContracterConcatWorker(a.ctx, a.contracter)

	go publishWorker.Start()
	go concatWorker.Start()

	a.dealer.ObserveSegments(a.ctx, a.wg)
	a.contracter.ObserveOrders(a.ctx, a.wg)

	return nil
}

// StartAPI _
func (a *DistributedConvertApp) StartAPI() error {
	h := handlers.NewDealerHandler(a.ctx, a.dealer)

	e := echo.New()

	remote.RegisterHandlers(e, h)

	return e.Start(":8080")
}

// AddTask _
func (a *DistributedConvertApp) AddTask(c *cli.Context) error {
	var err error

	segmentsPath := files.NewPath(".fftb/segments")

	err = segmentsPath.Create()

	if err != nil {
		return errors.Wrap(err, "Creating segments path")
	}

	inFile := files.NewFile(c.Args().Get(0))
	outFile := files.NewFile(c.Args().Get(1))

	// a.contracter, err = local.NewContracter(a.ctx, a.dealer, a.registry, segmentsPath)

	// if err != nil {
	// 	return errors.Wrap(err, "Initializing contracter")
	// }

	_, err = a.contracter.AddOrderToQueue(a.ctx, &models.ConvertContracterRequest{
		InFile:  inFile,
		OutFile: outFile,
		Params:  convertParamsFromFlags(c),
		Author:  a.publisher,
	})

	if err != nil {
		return errors.Wrap(err, "Publishing order")
	}

	a.logger.Debug("Order queued")

	return nil
}

// StartWorker _
func (a *DistributedConvertApp) StartWorker() error {
	var err error

	workerPath := files.NewPath(".fftb/worker")

	err = workerPath.Create()

	if err != nil {
		return errors.Wrap(err, "Initializing worker path")
	}

	a.workerInstance, err = worker.NewWorker(a.ctx, workerPath, a.dealer, a.storageClient)

	if err != nil {
		return errors.Wrap(err, "Initializing worker instance")
	}

	// workerDone := a.workerInstance.Start()
	a.workerInstance.Start()

	return nil
}

// ListOrders _
func (a *DistributedConvertApp) ListOrders(cliCtx *cli.Context) (string, error) {
	filters := make([]models.IOrderSearchCriteria, 0)

	if cliCtx.String("state") != "" {
		filters = append(filters, models.OrderStateFilter(cliCtx.String("state")))
	}

	orders, err := a.contracterInteractor.GetAllOrders(a.ctx, models.MergeOrderFilters(filters...))

	if err != nil {
		return "", err
	}

	headers := []string{"ID", "Input file", "Output file", "Progress", "State"}

	ordersData := make([][]string, 0)

	for _, orderItem := range orders {
		ordersData = append(ordersData,
			[]string{
				orderItem.ID,
				files.NewFile(orderItem.InputFile).Name(),
				files.NewFile(orderItem.OutputFile).Name(),
				fmt.Sprintf("%.2f%%", orderItem.Progress*100),
				orderItem.State,
			},
		)
	}

	return renderTable(headers, ordersData), nil
}

// ShowOrder _
func (a *DistributedConvertApp) ShowOrder(orderID string) (string, error) {
	orderItem, err := a.contracterInteractor.GetOrderByID(a.ctx, orderID)

	if err != nil {
		return "", err
	}

	headers := []string{"Attribute", "Value"}

	data := [][]string{
		{"ID", orderItem.ID},
		{"Input file", orderItem.InputFile},
		{"Output file", orderItem.OutputFile},
		{"State", orderItem.State},
		{"Progress", fmt.Sprintf("%.2f%%", orderItem.Progress*100)},
		{"Segments count", strconv.Itoa(orderItem.SegmentsCount)},
	}

	return renderTable(headers, data), nil
}

// ListSegments _
func (a *DistributedConvertApp) ListSegments(cliCtx *cli.Context, orderID string) (string, error) {
	filters := make([]models.ISegmentSearchCriteria, 0)

	if cliCtx.String("state") != "" {
		filters = append(filters, models.SegmentStateFilter(cliCtx.String("state")))
	}

	segmentItems, err := a.contracterInteractor.GetSegmentsByOrderID(a.ctx, orderID, models.MergeSegmentFilters(filters...))

	if err != nil {
		return "", err
	}

	headers := []string{"ID", "State", "Performer"}

	segmentsData := make([][]string, 0)

	for _, segmentItem := range segmentItems {
		segmentsData = append(segmentsData,
			[]string{
				segmentItem.ID,
				segmentItem.State,
				segmentItem.Performer,
			},
		)
	}

	return renderTable(headers, segmentsData), nil
}

// CancelOrder _
func (a *DistributedConvertApp) CancelOrder(orderID string) error {
	return a.contracterInteractor.CancelOrderByID(a.ctx, orderID)
}

// Wait _
func (a *DistributedConvertApp) Wait() <-chan struct{} {
	cancelSignal := make(chan os.Signal)

	signal.Notify(cancelSignal, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-a.ctx.Done():
				if a.workerInstance != nil {
					select {
					case <-a.workerInstance.Closed():
						a.logger.Debug("Worker terminated")
					case <-time.After(3 * time.Second):
						a.logger.Debug("Worker killed")
					}
				}

				err := a.registry.Persist()

				if err != nil {
					a.logger.WithError(err).Warn("Registry persisting problem")
				}

				<-a.registry.Closed()
				close(a.closed)
				return
			case <-cancelSignal:
				a.logger.Info("Terminating")
				a.cancel()
			}
		}
	}()

	return a.closed
}

func initStorage(ctx context.Context) (models.IStorageController, error) {
	storagePath := files.NewPath(".fftb/storage")

	err := storagePath.Create()

	if err != nil {
		return nil, errors.Wrap(err, "Creating storage path")
	}

	storage := local.NewStorageControl(storagePath)

	return storage, nil
}

func initRegistry(ctx context.Context) (models.IRegistry, error) {
	store, err := ubolt.NewClient(ctx, ".fftb/bolt.db")

	if err != nil {
		return nil, errors.Wrap(err, "Initializing ubolt ukvs store")
	}

	registry, err := registry.NewRegistry(ctx, store)

	if err != nil {
		return nil, errors.Wrap(err, "Initializing registry")
	}

	return registry, nil
}

func initStorageClient() (models.IStorageClient, error) {
	storageClientPath := files.NewPath(".fftb/storage_client")

	err := storageClientPath.Create()

	if err != nil {
		return nil, errors.Wrap(err, "Initializing storage client path")
	}

	storageClient := local.NewStorageClient(storageClientPath.FullPath())

	return storageClient, nil
}
