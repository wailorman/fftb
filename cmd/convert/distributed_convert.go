package convert

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/local"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/registry"
	"github.com/wailorman/fftb/pkg/distributed/ukvs/localfile"
	"github.com/wailorman/fftb/pkg/distributed/worker"
	"github.com/wailorman/fftb/pkg/files"
)

// DistributedCliConfig _
func DistributedCliConfig() *cli.Command {
	flags := convertParamsFlags()

	flags = append(flags, &cli.BoolFlag{
		Name: "worker",
	})

	return &cli.Command{
		Name:    "distributed-convert",
		Aliases: []string{"dconv"},
		Flags:   flags,
		Subcommands: []*cli.Command{
			&cli.Command{
				Name:  "add",
				Flags: convertParamsFlags(),
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}

					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					err = app.StartContracter(c)

					if err != nil {
						return errors.Wrap(err, "Starting contracter")
					}

					<-app.Wait()

					return nil
				},
			},
			&cli.Command{
				Name: "work",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "worker",
					},
				},
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}

					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					err = app.StartWorker()

					if err != nil {
						return errors.Wrap(err, "Starting worker")
					}

					<-app.Wait()

					return nil
				},
			},
		},
	}
}

// DistributedConvertApp _
type DistributedConvertApp struct {
	storage        models.IStorageController
	registry       models.IRegistry
	dealer         models.IDealer
	contracter     models.IContracter
	workerInstance *worker.Instance
	ctx            context.Context
	cancel         func()
	logger         *logrus.Entry
	closed         chan struct{}
}

// Init _
func (a *DistributedConvertApp) Init() error {
	var err error

	a.closed = make(chan struct{})

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.Formatter = &prefixed.TextFormatter{
		FullTimestamp: true,
		SpacePadding:  10,
	}
	a.logger = logger.WithField("prefix", "fftb.distributed")

	a.ctx, a.cancel = context.WithCancel(context.WithValue(context.Background(), ctxlog.LoggerContextKey, a.logger))

	a.storage, err = initStorage(a.ctx)

	if err != nil {
		return errors.Wrap(err, "Initializing storage")
	}

	a.registry, err = initRegistry(a.ctx)

	if err != nil {
		return errors.Wrap(err, "Initializing registry")
	}

	a.dealer, err = local.NewDealer(a.ctx, a.storage, a.registry)

	if err != nil {
		return errors.Wrap(err, "Initializing delaer")
	}

	return nil
}

// StartContracter _
func (a *DistributedConvertApp) StartContracter(c *cli.Context) error {
	defer a.cancel()

	var err error

	segmentsPath := files.NewPath(".fftb/segments")

	err = segmentsPath.Create()

	if err != nil {
		return errors.Wrap(err, "Creating segments path")
	}

	inFile := files.NewFile(c.Args().Get(0))

	a.contracter, err = local.NewContracter(&local.ContracterParameters{
		TempPath: segmentsPath,
		Dealer:   a.dealer,
	})

	if err != nil {
		return errors.Wrap(err, "Initializing contracter")
	}

	_, err = a.contracter.PrepareOrder(&models.ConvertContracterRequest{
		InFile: inFile,
		Params: convertParamsFromFlags(c),
	})

	a.logger.Debug("Order published")

	if err != nil {
		return errors.Wrap(err, "Publishing order")
	}

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

	a.workerInstance, err = worker.NewWorker(a.ctx, workerPath, a.dealer)

	if err != nil {
		return errors.Wrap(err, "Initializing worker instance")
	}

	// workerDone := a.workerInstance.Start()
	a.workerInstance.Start()

	return nil
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
					<-a.workerInstance.Closed()
					a.logger.Debug("Worker terminated")
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

func (a *DistributedConvertApp) terminateApp() {
	a.logger.Info("Terminating")
	a.cancel()

	if a.workerInstance != nil {
		<-a.workerInstance.Closed()
		a.logger.Debug("Worker terminated")
	}

	<-a.registry.Closed()
	close(a.closed)
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
	store, err := localfile.NewClient(ctx, ".fftb/store.json")

	if err != nil {
		return nil, errors.Wrap(err, "Initializing localfile ukvs store")
	}

	registry, err := registry.NewRegistry(ctx, store)

	if err != nil {
		return nil, errors.Wrap(err, "Initializing registry")
	}

	return registry, nil
}
