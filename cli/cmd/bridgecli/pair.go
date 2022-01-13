package main

import (
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/tak1827/evm-bridge/cli/pb"
	"github.com/tak1827/go-store/store"
)

var (
	IsWrapped bool
)

var pairCmd = &cobra.Command{
	Use:                        "pair",
	Short:                      "set and get contract pair address",
	DisableFlagParsing:         true,
	SuggestionsMinimumDistance: 2,
}

var pairSetCmd = &cobra.Command{
	Use:   "set [in-addr] [out-addr]",
	Short: "register contract address pair",
	Long:  `Register the in and out chain contract address pair`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		getConfig()

		inAddr, err := cast.ToStringE(args[0])
		handleErr(err)
		outAddr, err := cast.ToStringE(args[1])
		handleErr(err)

		if !common.IsHexAddress(inAddr) {
			handleErr(errors.New(fmt.Sprintf("invalid address formt in-addr: %s", inAddr)))
		}
		if !common.IsHexAddress(outAddr) {
			handleErr(errors.New(fmt.Sprintf("invalid address formt out-addr: %s", outAddr)))
		}

		pair := pb.Pair{
			Inaddr:  common.HexToAddress(inAddr).Hex(),
			Outaddr: common.HexToAddress(outAddr).Hex(),
			Intype:  pb.Pair_ORIGINAL,
		}

		if IsWrapped {
			pair.Intype = pb.Pair_WRAPPED
		}

		db, err := store.NewLevelDB(homeDir)
		handleErr(err)

		err = pair.Put(db)
		handleErr(err)

		fmt.Println("succeeded!")
	},
}

var pairGetCmd = &cobra.Command{
	Use:   "get [in-addr]",
	Short: "get contract address pair",
	Long:  `show the Registerd the in and out chain contract address pair by "in-addr"`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getConfig()

		inAddr, err := cast.ToStringE(args[0])
		handleErr(err)

		if !common.IsHexAddress(inAddr) {
			handleErr(errors.New(fmt.Sprintf("invalid address formt in-addr: %s", inAddr)))
		}

		db, err := store.NewLevelDB(homeDir)
		handleErr(err)

		pair, err := pb.GetPair(db, common.HexToAddress(inAddr).Hex())
		handleErr(err)

		spew.Dump(pair)
	},
}

func init() {
	pairSetCmd.Flags().BoolVar(&IsWrapped, "in-type-wrapped", false, "the type of in chain contract (`ORIGINAL` or `WRAPPED`) is `WRAPPED`")
	pairCmd.AddCommand(pairSetCmd)
	pairCmd.AddCommand(pairGetCmd)
	rootCmd.AddCommand(pairCmd)
}
