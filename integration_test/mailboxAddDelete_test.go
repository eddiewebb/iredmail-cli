package integrationTest

import (
	"database/sql"
	"fmt"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mailbox", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can add a mailbox", func() {
		if skipMailboxAddDelete && !isCI {
			Skip("can add a mailbox")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully added mailbox %s\n", mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		sqlQuery := `SELECT exists
		(SELECT * FROM mailbox
		WHERE username = ?);`

		err = db.QueryRow(sqlQuery, mailboxName1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))

		sqlQuery = `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = ? AND forwarding = ? 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(sqlQuery, mailboxName1, mailboxName1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(true))
	})

	It("can delete a mailbox", func() {
		if skipMailboxAddDelete && !isCI {
			Skip("can delete a mailbox")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "mailbox", "delete", "--force", mailboxName1)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully deleted mailbox %v\n", mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var exists bool

		sqlQuery := `SELECT exists
		(SELECT * FROM mailbox
		WHERE username = ?);`

		err = db.QueryRow(sqlQuery, mailboxName1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))

		sqlQuery = `SELECT exists
		(SELECT * FROM forwardings
		WHERE address = ? AND forwarding = ? 
		AND is_forwarding = 1 AND active = 1 AND is_alias = 0 AND is_maillist = 0);`

		err = db.QueryRow(sqlQuery, mailboxName1, mailboxName1).Scan(&exists)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(false))
	})

	It("can't add an existing mailbox", func() {
		if skipMailboxAddDelete && !isCI {
			Skip("can't add an existing mailbox")
		}

		cli := exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		err := cli.Run()
		Expect(err).NotTo(HaveOccurred())

		cli = exec.Command(cliPath, "mailbox", "add", mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		Expect(err).To(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Mailbox %v already exists\n", mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can add an mailbox with custom quota", func() {
		if skipMailboxAddDelete && !isCI {
			Skip("can add an mailbox with custom quota")
		}

		cli := exec.Command(cliPath, "mailbox", "add", "--quota", strconv.Itoa(customQuota), mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		actual := string(output)
		expected := fmt.Sprintf("Successfully added mailbox %s\n", mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var quota int

		sqlQuery := `SELECT quota FROM mailbox WHERE username = ?;`

		err = db.QueryRow(sqlQuery, mailboxName1).Scan(&quota)
		Expect(err).NotTo(HaveOccurred())
		Expect(quota).To(Equal(customQuota))
	})

	It("can add an mailbox with custom storage path", func() {
		if skipMailboxAddDelete && !isCI {
			Skip("can add an mailbox with custom storage path")
		}

		cli := exec.Command(cliPath, "mailbox", "add", "--storage-path", customStoragePath, mailboxName1, mailboxPW)
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := fmt.Sprintf("Successfully added mailbox %s\n", mailboxName1)

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}

		db, err := sql.Open("mysql", dbConnectionString)
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		var storageBaseDirectory, storageNode string

		sqlQuery := `SELECT storagebasedirectory, storagenode FROM mailbox WHERE username = ?;`

		err = db.QueryRow(sqlQuery, mailboxName1).Scan(&storageBaseDirectory, &storageNode)
		Expect(err).NotTo(HaveOccurred())

		Expect(storageBaseDirectory).To(Equal(filepath.Dir(customStoragePath)))
		Expect(storageNode).To(Equal(filepath.Base(customStoragePath)))
	})
})
