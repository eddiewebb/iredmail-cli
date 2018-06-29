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

// mailboxAliasAddCmd represents the add-alias command
var mailboxAliasAddCmd = &cobra.Command{
	Use:   "add-alias",
	Short: "Add mailbox alias (e.g. abuse -> post@domain.com)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires alias and mailbox (email address)")
		}

		err := emailx.Validate(args[0])
		if err == nil {
			return fmt.Errorf("Invalid alias format: \"%v\"", args[0])
		}

		err = emailx.Validate(args[1])
		if err != nil {
			return fmt.Errorf("Invalid mailbox format: \"%v\"", args[1])
		}

		args[1] = emailx.Normalize(args[1])

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		server, err := iredmail.New()
		if err != nil {
			fatal("%v\n", err)
		}
		defer server.Close()

		err = server.MailboxAliasAdd(args[0], args[1])
		if err != nil {
			fatal("%v\n", err)
		}

		success("Successfully added mailbox alias %v -> %v\n", args[0], args[1])
	},
}

func init() {
	mailboxCmd.AddCommand(mailboxAliasAddCmd)
	mailboxAliasAddCmd.SetUsageTemplate(usageTemplate("mailbox add-alias [alias] [mailbox_email]"))
}