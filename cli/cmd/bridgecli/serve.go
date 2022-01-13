package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	b "github.com/tak1827/evm-bridge/cli/bridge"
	"github.com/tak1827/evm-bridge/cli/client"
	"github.com/tak1827/evm-bridge/cli/log"
	"github.com/tak1827/evm-bridge/cli/pb"
	"github.com/tak1827/transaction-confirmer/confirm"
)

const (
	QueueSize              = 65536
	MIN_LOG_FETCH_INTERVAL = 3000 // 3s
)

var (
	InEndpoint  string
	OutEndpoint string
	HexBank     string

	PrivKey string

	LogFetchInterval int
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve briging functions",
	Long:  `Start fetching cotract events, and mint equivalent asset to the other chain`,
	Run: func(cmd *cobra.Command, args []string) {
		getServeConfig()
		start()
	},
}

func init() {
	serveCmd.Flags().StringVarP(&InEndpoint, "in-endpoint", "i", "http://localhost:8545", "in chain endpoint")
	serveCmd.Flags().StringVarP(&OutEndpoint, "out-endpoint", "o", "http://localhost:8545", "out chain endpoint")
	serveCmd.Flags().StringVar(&HexBank, "bank", "", "the bank contract address")
	rootCmd.AddCommand(serveCmd)
}

func getServeConfig() {
	getConfig()

	if PrivKey = viper.GetString("pri_key"); PrivKey == "" {
		logger.Fatal().Msg("please set `BRIDGECLI_PRI_KEY` as the env variable")
	}

	getConfigString("in-endpoint", &InEndpoint)
	getConfigString("out-endpoint", &OutEndpoint)
	getConfigString("bank", &HexBank)
	getConfigInt("log-fetch-interval", &LogFetchInterval)
	if LogFetchInterval < MIN_LOG_FETCH_INTERVAL {
		logger.Fatal().Msgf("`log-fetch-interval` is %d milisec, please set grater than %d milisic", LogFetchInterval, MIN_LOG_FETCH_INTERVAL)
	}
}

func confirmerOps() (ops []confirm.Opt) {
	workers := viper.GetInt("confirmer.workers")
	if workers != 0 {
		ops = append(ops, confirm.WithWorkers(workers))
		logger.Info().Msgf("confirmer.workers: %d", workers)
	}
	interval := viper.GetInt("confirmer.interval")
	if interval != 0 {
		ops = append(ops, confirm.WithWorkerInterval(int64(interval)))
		logger.Info().Msgf("confirmer.interval: %d", interval)
	}
	blocks := viper.GetInt("confirmer.confirmation-blocks")
	if blocks != 0 {
		ops = append(ops, confirm.WithConfirmationBlock(uint64(blocks)))
		logger.Info().Msgf("confirmer.confirmation-blocks: %d", blocks)
	}
	return
}

func start() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		rotator     = b.NewRotator(2)
	)

	c, err := client.NewClient(ctx, OutEndpoint, HexBank)
	handleErr(err)

	rc, err := client.NewReadClient(ctx, InEndpoint, HexBank)
	handleErr(err)

	confirmer := confirm.NewConfirmer(&c, QueueSize, confirmerOps()...)

	bridge, err := b.NewBridge(ctx, &c, &rc, &confirmer, PrivKey, homeDir)
	handleErr(err)

	err = bridge.Start(ctx)
	handleErr(err)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	timer := time.NewTicker(time.Duration(int64(LogFetchInterval)/2) * time.Millisecond)
	defer timer.Stop()

	for {
		select {
		case <-sigCh:
			log.Logger.Info().Msg("shutting down...")
			bridge.Close(cancel, 3, true)
			return
		case <-timer.C:
			switch rotator.Rotate() {
			case b.SlotERC20:
				if !(len(bridge.EventMapERC20) == 0) {
					continue
				}

				err = bridge.ConfirmedBlockERC20.Put(bridge.DB, pb.BlockERC20)
				handleErr(err)

				bridge.ConfirmedBlockERC20.Number, err = bridge.FetchERC20(ctx)
				handleErr(err)

			case b.SlotNFT:
				if !(len(bridge.EventMapNFT) == 0) {
					continue
				}

				bridge.ConfirmedBlockNFT.Put(bridge.DB, pb.BlockNFT)
				handleErr(err)

				bridge.ConfirmedBlockNFT.Number, err = bridge.FetchNFT(ctx)
				handleErr(err)
			}
		default:
		}
	}
}
