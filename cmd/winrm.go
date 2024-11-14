//go:build all || remote || winrm

package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/masterzen/winrm"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var winrmVars = struct {
	Hostname string
	Username string
	Password string
	Ntlm     bool
	Cmd      string
	Port     int
	Encoded  bool
	Https    bool
	Insecure bool
	Cacert   string
	Gencert  bool
	Certsize string
	Timeout  string
}{}

// winrmCmd represents the winrm command
var winrmCmd = &cobra.Command{
	Use:   "winrm",
	Short: "WinRM REPL command",
	Long:  `This subcommand supports an interactive session for the WinRm protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		commands := cmd.Flags().Args()

		pass := winrmVars.Password
		if winrmVars.Encoded {
			data, err := base64.StdEncoding.DecodeString(pass)
			if err != nil {
				panic(err)
			}
			pass = strings.TrimRight(string(data), "\r\n")
		}

		var (
			certBytes      []byte
			err            error
			connectTimeout time.Duration
			exitCode       int
		)

		// certBytes := nil
		// if winrmVars.Cacert != "" {
		// 	certBytes, err := os.ReadFile(cacert)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }

		for _, command := range commands {
			if command == "" {
				check(errors.New("ERROR: Please enter the command to execute on the command line"))
			}

			connectTimeout, err = time.ParseDuration(winrmVars.Timeout)
			if err != nil {
				panic(err)
			}

			endpoint := winrm.NewEndpoint(winrmVars.Hostname, winrmVars.Port, winrmVars.Https, winrmVars.Insecure, nil, certBytes, nil, connectTimeout)

			params := winrm.DefaultParameters
			if winrmVars.Ntlm {
				params.TransportDecorator = func() winrm.Transporter { return &winrm.ClientNTLM{} }
			}

			client, err := winrm.NewClientWithParameters(endpoint, winrmVars.Username, winrmVars.Password, params)
			if err != nil {
				panic(err)
			}

			if isatty.IsTerminal(os.Stdin.Fd()) {
				exitCode, err = client.Run(command, os.Stdout, os.Stderr)
			} else {
				exitCode, err = client.RunWithInput(command, os.Stdout, os.Stderr, os.Stdin)
			}
			if err != nil {
				panic(err)
			}
		}

		os.Exit(exitCode)
	},
}

func pickSizeCert(size string) int {
	switch size {
	case "512":
		return 512
	case "1024":
		return 1024
	case "2048":
		return 2048
	case "4096":
		return 4096
	default:
		return 2048
	}
}

// generic check error func
func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(winrmCmd)

	winrmCmd.PersistentFlags().StringVarP(&winrmVars.Hostname, "host", "H", "localhost", "winrm host")
	winrmCmd.PersistentFlags().IntVarP(&winrmVars.Port, "port", "P", 5985, "winrm port")

	winrmCmd.Flags().StringVarP(&winrmVars.Username, "username", "u", "root", "winrm admin username")
	winrmCmd.Flags().StringVarP(&winrmVars.Password, "password", "p", "", "winrm admin password")
	winrmCmd.Flags().BoolVarP(&winrmVars.Ntlm, "ntlm", "N", false, "use use ntlm auth")
	winrmCmd.Flags().BoolVarP(&winrmVars.Encoded, "encoded", "e", false, "use base64 encoded password")
	winrmCmd.Flags().BoolVarP(&winrmVars.Https, "https", "S", false, "use https")
	winrmCmd.Flags().BoolVarP(&winrmVars.Insecure, "insecure", "i", false, "skip SSL validation")
	winrmCmd.Flags().StringVarP(&winrmVars.Cacert, "cacert", "c", "", "CA certificate to use")
	winrmCmd.Flags().BoolVarP(&winrmVars.Gencert, "gencert", "g", false, "Generate x509 client certificate to use with secure connections")
	winrmCmd.Flags().StringVarP(&winrmVars.Certsize, "certsize", "s", "", "Priv RSA key between 512, 1024, 2048, 4096. Default :2048")
	winrmCmd.Flags().StringVarP(&winrmVars.Timeout, "timeout", "t", "0s", "connection timeout")
}
