package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loadCommand)
	rootCmd.AddCommand(expCommand())
	rootCmd.AddCommand(printCommand())
	rootCmd.AddCommand(goldCommand())
	rootCmd.AddCommand(itemCommand())
	rootCmd.AddCommand(extraCommand())
}

var rootCmd = &cobra.Command{
	Use:   "rpgsave",
	Short: "RpgSAVE is a cli tool for editing game save file(.rpgsave) by RPG Maker MV",
}

var loadCommand = &cobra.Command{
	Use:     "load",
	Aliases: []string{"open"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("invalid argument, filename is required")
		}

		filename := args[0]
		if !isFileExist(filename) {
			return fmt.Errorf("file <%s> not found", filename)
		}

		GAME.FilePath = filename
		if err := GAME.loadSaveFile(filename); err != nil {
			return err
		}

		if err := GAME.WriteConfig(); err != nil {
			return err
		}

		GAME.ItemPath = findItemsJSON(filename)
		fmt.Println("[ load ]", GAME.ItemPath)
		if !isFileExist(GAME.ItemPath) {
			fmt.Printf("warning: %s not found\n", GAME.ItemPath)
			return nil
		}

		if _, err := GAME.loadItems(GAME.ItemPath); err != nil {
			fmt.Printf("warning: failed to open %s not found: %s\n", GAME.ItemPath, err.Error())
			return nil
		}

		if err := GAME.WriteConfig(); err != nil {
			return err
		}

		return nil
	},
}

func extraCommand() *cobra.Command {
	extraNameMap := map[string]int{
		"MaxHP": MaxHP,
		"MHP":   MaxHP,
		"MaxMP": MaxMP,
		"MMP":   MaxMP,
		"ATK":   ATK,
		"DEF":   DEF,
		"MAT":   MAT,
		"MDF":   MDF,
		"AGI":   AGI,
		"LUK":   LUK,
	}

	var setExtra *map[string]int64

	cmd := &cobra.Command{
		Use:   "extra",
		Short: "Update _paramPlus fields in .rpgsave file",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			for k := range *setExtra {
				if _, ok := extraNameMap[k]; !ok {
					return fmt.Errorf("ERROR: invalid key [%s]", k)
				}
			}

			return GAME.Load()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			items := []string{}
			for k, num := range *setExtra {
				name := extraNameMap[k]
				old := GAME.save.Extra(name)
				GAME.save.SetExtra(name, num)
				items = append(items, k, fmt.Sprintf("%d (%d)", num, old))
			}
			cmd.OutOrStdout().Write([]byte(uPrint(items...)))
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return GAME.Unload()
		},
	}
	setExtra = cmd.Flags().StringToInt64P("set", "s", nil, "[MHP|MMP|ATK|DEF|MAT|MDF|AGI|LUK]=NUM")

	return cmd
}

func itemCommand() *cobra.Command {
	var showAll *bool
	cmd := &cobra.Command{
		Use:   "item",
		Short: "See items",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return GAME.Load()
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			log.Println("show-all?", *showAll)

			msgs := []string{}
			items := GAME.Items(*showAll)
			line := "+------+------+-------------"
			msgs = append(msgs, line)
			msgs = append(msgs, fmt.Sprintf("| %4s | %4s | %s", "ID", "NUM", "NAME"))
			msgs = append(msgs, line)
			for _, u := range items {
				msgs = append(msgs, fmt.Sprintf("| %4d | %4d | %s", u.ID, u.count, u.Name))
			}
			msgs = append(msgs, line)
			cmd.OutOrStdout().Write([]byte(strings.Join(msgs, "\n")))

			return nil
		},
	}
	showAll = cmd.Flags().Bool("show-all", false, "item -a")

	cmd.AddCommand(&cobra.Command{
		Use:   "set",
		Short: "set ID NUM",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("invalid argument, ID, NUM required")
			}

			return GAME.Load()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid argument %s, ID required", args[0])
			}

			num, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid argument %s, NUM required", args[1])
			}

			u := GAME.save.Item(id)
			GAME.save.SetItem(id, num)
			u1 := GAME.save.Item(id)

			msg := fmt.Sprintf("------------ ITEM ------------\n#%d | %d (%d) | %s\n------------------------------", id, u1.count, u.count, u1.Name)
			cmd.OutOrStdout().Write([]byte(msg))
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return GAME.Unload()
		},
	})

	return cmd
}

func goldCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gold",
		Short: "Change num of gold",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return GAME.Load()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			msg := uPrint(
				"GOLD", fmt.Sprintf("%d", GAME.save.Gold()),
			)
			cmd.OutOrStdout().Write([]byte(msg))

			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return GAME.Unload()
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use: "set",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("invalid argument, number required")
			}

			return GAME.Load()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			num, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid argument %s, number required", args[0])
			}

			msg := uPrint(
				"Gold", fmt.Sprintf("%d (%d)", num, GAME.save.Gold()),
			)
			cmd.OutOrStdout().Write([]byte(msg))

			GAME.save.SetGold(num)
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return GAME.Unload()
		},
	})

	return cmd
}

func printCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "print",
		Aliases: []string{"info"},
		Short:   "Show info of actor",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return GAME.Load()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			msg := uPrint(
				"NAME", GAME.save.Name(),

				"EXP", fmt.Sprintf("%d", GAME.save.Exp()),
				"GOLD", fmt.Sprintf("%d", GAME.save.Gold()),

				"MaxHP", fmt.Sprintf("%.0f", GAME.save.Extra(MaxHP)),
				"MaxMP", fmt.Sprintf("%.0f", GAME.save.Extra(MaxMP)),

				"ATK", fmt.Sprintf("%.0f", GAME.save.Extra(ATK)),
				"DEF", fmt.Sprintf("%.0f", GAME.save.Extra(DEF)),
				"MAT", fmt.Sprintf("%.0f", GAME.save.Extra(MAT)),
				"MDF", fmt.Sprintf("%.0f", GAME.save.Extra(MDF)),
				"AGI", fmt.Sprintf("%.0f", GAME.save.Extra(AGI)),
				"LUK", fmt.Sprintf("%.0f", GAME.save.Extra(LUK)),
			)
			cmd.OutOrStdout().Write([]byte(msg))

			GAME.save.Items()

			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return GAME.Unload()
		},
	}

	return cmd
}

func expCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exp",
		Short: "Get or set actor's exp",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return GAME.Load()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			msg := uPrint(
				"NAME", GAME.save.Name(),
				"EXP", fmt.Sprintf("%d", GAME.save.Exp()),
			)
			cmd.OutOrStdout().Write([]byte(msg))
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return GAME.Unload()
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use: "add",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("invalid arguments")
			}

			return GAME.Load()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			num, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid argument %s, number required", args[0])
			}

			msg := uPrint(
				"NAME", GAME.save.Name(),
				"EXP", fmt.Sprintf("%d (%d)", GAME.save.Exp()+num, GAME.save.Exp()),
			)
			cmd.OutOrStdout().Write([]byte(msg))

			GAME.save.AddExp(int64(num))
			return nil
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			GAME.Unload()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use: "set",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("invalid argument, number required")
			}

			return GAME.Load()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			num, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid argument %s, number required", args[0])
			}

			msg := uPrint(
				"NAME", GAME.save.Name(),
				"EXP", fmt.Sprintf("%d (%d)", GAME.save.Exp(), num),
			)
			cmd.OutOrStdout().Write([]byte(msg))

			GAME.save.AddExp(int64(num))
			return nil
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			log.Println("post-run")
			// GAME.Unload()
		},
	})

	return cmd
}
