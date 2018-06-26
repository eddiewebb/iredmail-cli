package integrationTest

import (
	"fmt"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("user list", func() {
	BeforeEach(func() {
		err := resetDB()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can list users", func() {
		if skipUserList && !isCI {
			Skip("can list users")
		}

		users := []string{
			userName1,
			userName2,
			userName3,
			userName4,
			userName5,
		}

		for _, user := range users {
			cli := exec.Command(cliPath, "user", "add", user, userPW)
			err := cli.Run()
			Expect(err).NotTo(HaveOccurred())
		}

		cli := exec.Command(cliPath, "user", "list")
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := loadGolden("can_list_users")

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})

	It("can list users and filter result", func() {
		if skipUserList && !isCI {
			Skip("can list users and filter result")
		}

		users := []string{
			userName1,
			userName2,
			userName3,
			userName4,
			userName5,
		}

		for _, user := range users {
			cli := exec.Command(cliPath, "user", "add", user, userPW)
			err := cli.Run()
			Expect(err).NotTo(HaveOccurred())
		}

		cli := exec.Command(cliPath, "user", "list", "-f", "domain.com")
		output, err := cli.CombinedOutput()
		if err != nil {
			Fail(string(output))
		}

		actual := string(output)
		expected := loadGolden("can_list_users_and_filter_result")

		if !reflect.DeepEqual(actual, expected) {
			Fail(fmt.Sprintf("actual = %s, expected = %s", actual, expected))
		}
	})
})
