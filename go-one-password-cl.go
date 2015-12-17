/* go-one-password-cl.go

   Creates a small, self-contained binary executable which runs on the
   command-line, using the go-one-password library

*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/dpapathanasiou/go-one-password"
	"os"
	"os/exec"
	"strings"
)

// getPassphrase requests the passphrase string as a standard input prompt instead of the command line args
// (because otherwise the passphrase would be visible in this user's shell history)
func getPassphrase() string {
	fmt.Print("What's your passphrase? (or ctrl-c to quit) ")
	// do not display typed characters to appear on the screen (linux only)
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	reader := bufio.NewReader(os.Stdin)
	passphrase, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Whoa, error! %s\n", err.Error())
		return ""
	}
	// re-enable typed characters to appear on the screen
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	return strings.TrimSpace(passphrase)
}

func main() {

	var (
		passphrase, username, hostname, specials, result string
		passwordLength                                   int
	)
	const (
		defaultInput  = ""
		defaultPwdLen = 16
	)

	flag.StringVar(&hostname, "host", defaultInput, "(required) the website you want to login to (e.g. \"amazon.com\")")
	flag.StringVar(&username, "user", defaultInput, "(required) the username or email address you use to login")
	flag.StringVar(&specials, "spec", defaultInput, "(optional) if the website requires one or more \"special\" characters in the password (e.g., \"#%*\" etc.) specify one or more of them here")
	flag.IntVar(&passwordLength, "plen", defaultPwdLen, fmt.Sprintf("(optional) set the resulting password length (the default is %d)", defaultPwdLen))
	flag.Parse()

	if len(hostname) < 1 && len(username) < 1 {
		fmt.Println("Usage:")
		flag.PrintDefaults()
	} else {
		passphrase = ""
		// prompt the user for a passphrase from standard input
		// (we do this to avoid using a command-line arg for the passphrase,
		// which would log the passphrase in the user's shell history)
		for len(passphrase) < 1 {
			passphrase = getPassphrase()
		}

		i := 0
		valid := false
		// keep generating passwords (using an updated iteration number) until we get one that meets the pwdIsValid() criteria
		for !valid {
			result = onepassword.GetCandidatePwd(passphrase, username, hostname, 12, i)
			valid = onepassword.PwdIsValid(result, passwordLength)
			i += 1
		}
		// success: display the result (the first passwordLength characters, with special chars at the end, if any)
		fmt.Print(fmt.Sprintf("Your password for %s logging in as user %s is:\n\n%s\n\n", hostname, username, strings.Join([]string{result[0:(passwordLength - len(specials))], specials}, "")))
	}
}
