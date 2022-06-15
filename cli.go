package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rpgsave",
	Short: "RpgSAVE is a cli tool for editing game save file(.rpgsave) by RPG Maker MV",
}

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open fileN.rpgsave",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.ErrOrStderr().Write([]byte("error: invalid filename, .rpgsave is required\n"))
			return
		}

		pathname := args[0]
		if path.Ext(pathname) != ".rpgsave" {
			cmd.ErrOrStderr().Write([]byte("error: invalid filename, .rpgsave is required\n"))
			return
		}

		if _, err := os.Stat(pathname); os.IsNotExist(err) {
			cmd.ErrOrStderr().Write([]byte("error: file not found\n"))
			return
		}

		DATA.Save.FilePath = pathname
		if err := DATA.WriteConfig(); err != nil {
			cmd.ErrOrStderr().Write([]byte("error: failed to update config\n"))
			return
		}

		msg := fmt.Sprintf("ok, save file is set to: <%s>\n", pathname)
		cmd.OutOrStdout().Write([]byte(msg))
	},
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Auto backup fileN.rpgsave",
	RunE: func(cmd *cobra.Command, args []string) error {
		src := DATA.Save.FilePath
		dst := fmt.Sprintf("%s.backup-%s", src, time.Now().Format("0102-1504"))

		srcFile, err := os.Open(src)
		if err != nil {
			return err
		}

		dstFile, err := os.Create(dst)
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}

		return nil
	},
}

var printCmd = &cobra.Command{
	Use:     "print",
	Aliases: []string{"p", "info"},
	Short:   "Print info current actor",
	Run: func(cmd *cobra.Command, args []string) {
		DATA.save.Print()
	},
}

func init() {
	var expCmd = &cobra.Command{
		Use:   "exp",
		Short: "Change current actor's exp",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("invalid arguments")
			}

			num, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := fmt.Sprintf("Exp: %d => %d\n", DATA.save.Exp(), num)
			cmd.OutOrStdout().Write([]byte(msg))

			DATA.save.SetExp(int64(num))
			return nil
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			DATA.Unload()
		},
	}
	expCmd.AddCommand(&cobra.Command{
		Use: "add",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("invalid arguments")
			}

			num, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := fmt.Sprintf("Exp: %d => %d (%d)\n", DATA.save.Exp(), DATA.save.Exp()+num, num)
			cmd.OutOrStdout().Write([]byte(msg))

			DATA.save.AddExp(int64(num))
			return nil
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			DATA.Unload()
		},
	})

	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(printCmd)
	rootCmd.AddCommand(expCmd)
	rootCmd.AddCommand(backupCmd)
}
