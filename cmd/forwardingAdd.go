// Copyright © 2018 Christian Nolte
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"

	"github.com/drlogout/iredmail-cli/iredmail"
	"github.com/goware/emailx"
	"github.com/spf13/cobra"
)

// forwardingAddCmd represents the add command
var forwardingAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add forwarding (e.g. post@domain.com -> info@example.com)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires user and destination email")
		}

		err := emailx.Validate(args[0])
		if err != nil {
			return fmt.Errorf("Invalid user email format: \"%v\"", args[0])
		}

		args[0] = emailx.Normalize(args[0])

		err = emailx.Validate(args[1])
		if err != nil {
			return fmt.Errorf("Invalid destination email format: \"%v\"", args[1])
		}

		args[1] = emailx.Normalize(args[1])

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n")
		}
		defer server.Close()

		userEmail, destinationEmail := args[0], args[1]

		f, err := server.UserAddForwarding(userEmail, destinationEmail)
		if err != nil {
			fatal("%v\n")
		}

		success("Successfully added forwarding %v -> %v\n", f.Address, f.Forwarding)
	},
}

func init() {
	forwardingCmd.AddCommand(forwardingAddCmd)

	forwardingAddCmd.SetUsageTemplate(usageTemplate("forwarding add [user_email] [destination_email]"))
}