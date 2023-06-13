package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Equationzhao/g/index"
	"github.com/Equationzhao/pathbeautify"
	"github.com/urfave/cli/v2"
)

var indexFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:               "disable-index",
		Aliases:            []string{"di", "no-update"},
		Usage:              "disable updating index",
		Category:           "Index",
		DisableDefaultText: true,
	},
	&cli.BoolFlag{
		Name:               "rebuild-index",
		Aliases:            []string{"ri", "remove-all"},
		Usage:              "rebuild index",
		DisableDefaultText: true,
		Category:           "Index",
		Action: func(context *cli.Context, b bool) error {
			if b {
				err := index.RebuildIndex()
				if err != nil {
					return err
				}
			}
			return Err4Exit{}
		},
	},
	&cli.BoolFlag{
		Name:               "fuzzy",
		Aliases:            []string{"fz", "f"},
		Usage:              "fuzzy search",
		DisableDefaultText: true,
		Category:           "Index",
		EnvVars:            []string{"G_FZF"},
	},
	&cli.StringSliceFlag{
		Name:     "remove-index",
		Aliases:  []string{"rm"},
		Usage:    "remove paths from index",
		Category: "Index",
		Action: func(context *cli.Context, i []string) error {
			var errSum error = nil

			beautification := true
			if context.Bool("np") { // --no-path-transform
				beautification = false
			}

			for _, s := range i {
				if beautification {
					s = pathbeautify.Transform(s)
				}

				// get absolute path
				r, err := filepath.Abs(s)
				if err != nil {
					errSum = errors.Join(errSum, fmt.Errorf("remove-path: %w", err))
					continue
				}

				err = index.Delete(r)
				if err != nil {
					errSum = errors.Join(errSum, fmt.Errorf("remove-path: %w", err))
				}
			}
			if errSum != nil {
				return errSum
			} else {
				return Err4Exit{}
			}
		},
	},
	&cli.BoolFlag{
		Name:               "list-index",
		Aliases:            []string{"li"},
		Usage:              "list index",
		DisableDefaultText: true,
		Category:           "Index",
		Action: func(context *cli.Context, b bool) error {
			if b {
				keys, _, err := index.All()
				if err != nil {
					return err
				}
				for i := 0; i < len(keys); i++ {
					fmt.Println(keys[i])
				}
			}
			return Err4Exit{}
		},
	},
	&cli.BoolFlag{
		Name:     "remove-current-path",
		Aliases:  []string{"rcp", "rc", "rmc"},
		Usage:    "remove current path from index",
		Category: "Index",
		Action: func(context *cli.Context, b bool) error {
			if b {
				r, err := os.Getwd()
				if err != nil {
					return err
				}
				err = index.Delete(r)
				if err != nil {
					return err
				}
			}
			return Err4Exit{}
		},
	},
}
