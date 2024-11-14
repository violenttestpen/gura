//go:build all || remote || ssh

package cmd

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var sshVars = struct {
	Host string
	Port uint

	Username             string
	Password             string
	IdentityFile         string
	IdentityFilePassword string
}{}

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH REPL command",
	Long:  `This subcommand supports an interactive session for the SSH protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		commands := cmd.Flags().Args()
		authMethod := make([]ssh.AuthMethod, 0, 2)

		if _, err := os.Stat(sshVars.IdentityFile); err == nil {
			keyBytes, err := os.ReadFile(sshVars.IdentityFile)
			if err != nil {
				panic(fmt.Sprintf("Failed to read private key: %s", err))
			}

			var key ssh.Signer
			if sshVars.IdentityFilePassword != "" {
				key, err = ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte(sshVars.IdentityFilePassword))
				if err != nil {
					panic(fmt.Sprintf("Failed to parse private key: %s", err))
				}
			} else {
				key, err = ssh.ParsePrivateKey(keyBytes)
				if err != nil {
					panic(fmt.Sprintf("Failed to parse private key: %s", err))
				}
			}
			authMethod = append(authMethod, ssh.PublicKeys(key))
		}

		if sshVars.Password != "" {
			authMethod = append(authMethod, ssh.Password(sshVars.Password))
		}

		sshClientConfig := &ssh.ClientConfig{
			User:            sshVars.Username,
			Auth:            authMethod,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         5 * time.Second,
		}

		address := net.JoinHostPort(sshVars.Host, strconv.Itoa(int(sshVars.Port)))
		client, err := ssh.Dial("tcp", address, sshClientConfig)
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to %s: %s", sshVars.Host, err))
		}

		if verbose {
			fmt.Printf("Connected to %s\n", sshVars.Host)
		}
		defer client.Close()

		for _, command := range commands {
			session, err := client.NewSession()
			if err != nil {
				panic(fmt.Sprintf("failed to create session: %s", err))
			}
			defer session.Close()

			stdout, err := session.StdoutPipe()
			if err != nil {
				panic(fmt.Sprintf("failed to get stdout: %s", err))
			}
			if err := session.Start(command); err != nil {
				panic(fmt.Sprintf("failed to start command: %s", err))
			}

			output, err := io.ReadAll(stdout)
			if err != nil {
				panic(fmt.Sprintf("failed to read stdout: %s", err))
			}
			if err := session.Wait(); err != nil {
				panic(fmt.Sprintf("command failed: %s", err))
			}
			fmt.Print(string(output))
		}
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)

	sshCmd.PersistentFlags().StringVarP(&sshVars.Host, "host", "H", "localhost", "Address of target")
	sshCmd.PersistentFlags().UintVarP(&sshVars.Port, "port", "P", 22, "Port of target")

	sshCmd.Flags().StringVarP(&sshVars.Username, "username", "u", "root", "Username to connect to database")
	sshCmd.Flags().StringVarP(&sshVars.Password, "password", "p", "", "Password to connect to database")
	sshCmd.Flags().StringVarP(&sshVars.IdentityFile, "identity_file", "i", "", "Identity file to connect to database")
	sshCmd.Flags().StringVarP(&sshVars.IdentityFilePassword, "identity_file_password", "I", "", "Password to unlock identity file")
}
