package dconfig

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/twitchtv/twirp"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/yukitsune/lokirus"

	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
	"github.com/wailorman/fftb/pkg/distributed/worker"
	"github.com/wailorman/fftb/pkg/stdouthook"
)

type Instance struct {
	viper *viper.Viper
}

func New() (*Instance, error) {
	i := &Instance{}
	i.viper = viper.New()

	i.viper.SetConfigType("yaml")
	i.viper.SetConfigName("fftbd_config")
	i.viper.AddConfigPath(".")

	i.viper.SetDefault("worker_name", "default")

	i.viper.SetDefault("ffmpeg_path", "ffmpeg")
	i.viper.SetDefault("ffprobe_path", "ffprobe")
	i.viper.SetDefault("rclone_path", "rclone")
	i.viper.SetDefault("tmp_path", "tmp/")
	i.viper.SetDefault("local_remotes_map", map[string]string{})

	i.viper.SetDefault("dealer.url", "http://localhost:3000")
	i.viper.SetDefault("dealer.secret", nil)

	i.viper.SetDefault("threads_count", 1)

	i.viper.SetDefault("log_level", "info")
	i.viper.SetDefault("log_disable_colors", false)

	i.viper.SetDefault("loki.url", nil)
	i.viper.SetDefault("loki.user", nil)
	i.viper.SetDefault("loki.password", nil)

	if err := i.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			i.viper.SafeWriteConfigAs("./fftbd_config.yml")
		} else {
			return nil, errors.Wrap(err, "Reading config")
		}
	}

	return i, nil
}

func (i *Instance) ThreadsCount() int {
	return i.viper.GetInt("threads_count")
}

type ThreadConfigParams struct {
	Ctx    context.Context
	Dealer pb.Dealer
	Logger *logrus.Entry
	Wg     *chwg.ChannelledWaitGroup

	ThreadNumber  int
	Authorization string
}

func (i *Instance) ThreadConfig(params *ThreadConfigParams) worker.WorkerParams {
	tmpPath := i.viper.GetString("threads_config." + strconv.Itoa(params.ThreadNumber) + ".tmp_path")

	if tmpPath == "" {
		tmpPath = i.viper.GetString("tmp_path")
	}

	return worker.WorkerParams{
		Ctx:           params.Ctx,
		Dealer:        params.Dealer,
		Logger:        params.Logger,
		Wg:            params.Wg,
		Authorization: params.Authorization,

		TmpPath:          tmpPath,
		RcloneConfigPath: i.viper.GetString("rclone_config_path"),
		RclonePath:       i.viper.GetString("rclone_path"),
		FFmpegPath:       i.viper.GetString("ffmpeg_path"),
		FFprobePath:      i.viper.GetString("ffprobe_path"),
		LocalRemotesMap:  i.viper.GetStringMapString("local_remotes_map"),
	}
}

func (i *Instance) BuildDealer(logger *logrus.Entry) pb.Dealer {
	return pb.NewDealerProtobufClient(
		i.viper.GetString("dealer.url"),
		&http.Client{},
		twirp.WithClientInterceptors(dlog.TwirpLogInterceptor(logger)),
	)
}

func (i *Instance) BuildLogger() (*logrus.Entry, error) {
	logger := logrus.New()
	parsedLevel, err := logrus.ParseLevel(i.viper.GetString("log_level"))

	if err != nil {
		return nil, errors.Wrap(err, "Parsing log level")
	}

	logger.SetLevel(logrus.TraceLevel)
	logger.SetOutput(io.Discard)

	stdoutFormatter := &prefixed.TextFormatter{
		ForceFormatting: true,
		DisableColors:   i.viper.GetBool("log_disable_colors"),
		ForceColors:     !i.viper.GetBool("log_disable_colors"),
	}

	logger.AddHook(stdouthook.New(stdouthook.HookParams{
		Level:     parsedLevel,
		Formatter: stdoutFormatter,
	}))

	if i.viper.GetString("loki.url") != "" {
		logger.AddHook(lokiHook(&lokiParams{
			workerName: i.viper.GetString("worker_name"),
			lokiUrl:    i.viper.GetString("loki.url"),
			lokiUser:   i.viper.GetString("loki.user"),
			lokiPass:   i.viper.GetString("loki.password"),
		}))
	}

	return ctxlog.WithPrefix(logger, "fftb"), nil
}

func (i *Instance) BuildAccessToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"worker_name": i.viper.GetString("worker_name"),
	})

	return token.SignedString([]byte(i.viper.GetString("dealer.secret")))
}

type lokiParams struct {
	workerName string
	lokiUrl    string
	lokiUser   string
	lokiPass   string
}

func lokiHook(params *lokiParams) *lokirus.LokiHook {
	opts := lokirus.NewLokiHookOptions().
		// Grafana doesn't have a "panic" level, but it does have a "critical" level
		// https://grafana.com/docs/grafana/latest/explore/logs-integration/
		WithLevelMap(lokirus.LevelMap{logrus.PanicLevel: "critical"}).
		WithFormatter(&prefixed.TextFormatter{ForceFormatting: true, ForceColors: true}).
		WithStaticLabels(lokirus.Labels{
			"app":         "fftb",
			"worker_name": params.workerName,
		})

	if params.lokiUser != "" {
		opts = opts.WithBasicAuth(params.lokiUser, params.lokiPass)
	}

	hook := lokirus.NewLokiHookWithOpts(
		params.lokiUrl,
		opts)

	return hook
}
