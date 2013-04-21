/* go-one-password.go

This code creates a small, self-contained binary executable which runs on the command-line,
capable of generating a unique password for different sites and usernames based on a single,
private (i.e., known only to the person running this program) passphrase.

This code requires the scrypt package, so run:

$ go get code.google.com/p/go.crypto/scrypt

before trying to build.

*/

package main

import (
    "bufio"
    "bytes"
    "code.google.com/p/go.crypto/scrypt"
    "crypto/hmac"
    "crypto/sha512"
    "encoding/base64"
    "flag"
    "fmt"
    "os"
    "regexp"
    "strings"
)

// encodeBase64 accepts a byte array and returns base64 representation of it as pointer to a bytes.Buffer
func encodeBase64(data []byte) *bytes.Buffer {
    bb := &bytes.Buffer{}
    encoder := base64.NewEncoder(base64.StdEncoding, bb)
    encoder.Write([]byte(data))
    encoder.Close()
    return bb
}

// getPassphrase requests the passphrase string as a standard input prompt instead of the command line args
// (because otherwise the passphrase would be visible in this user's shell history)
func getPassphrase() string {
    fmt.Print("What's your passphrase? (or ctrl-c to quit) ")
    reader := bufio.NewReader(os.Stdin)
    passphrase, err := reader.ReadString('\n')
    if err != nil {
        fmt.Fprintf(os.Stderr, "Whoa, error! %s\n", err.Error())
        return ""
    }
    return strings.TrimSpace(passphrase)
}

// getCandidatePwd computes a password string based on the passphrase, hostname, username, a generation number, and an iteration number
func getCandidatePwd(pp, usr, host string, g, i int) string {
    // use scrypt to generate a derived key (dk) based on the passphrase and username
    dk, err := scrypt.Key([]byte(pp), []byte(usr), 16384, 8, 1, 32)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Whoa, error! %s\n", err.Error())
        return ""
    }
    // use dk as the shared private key in the HMAC-SHA512() generator
    // the value to be hashed consists of a concatenation of the username, hostname, generation and iteration numbers
    h := hmac.New(sha512.New, dk)
    h.Write([]byte(strings.Join([]string{"G1P,v.1.0", usr, host, fmt.Sprintf("%d", g), fmt.Sprintf("%d", i)}, "")))
    // return the result as a base64 encoded string
    return fmt.Sprintf("%s", encodeBase64(h.Sum(nil)))
}

var (
    IS_ALPHANUMERIC = regexp.MustCompile("[A-Za-z0-9]+")
    HAS_UPPERCASE   = regexp.MustCompile("[A-Z]{1,5}")
    HAS_LOWERCASE   = regexp.MustCompile("[a-z]{1,5}")
    HAS_NUMERICS    = regexp.MustCompile("[0-9]{1,5}")
)

// pwdIsValid checks the string created by getCandidatePwd and returns a boolean depending on whether or not it meets the criteria:
// the first pwdLen characters are all alphanumeric, and there is at least one (but no more than 5) uppercase, lowercase and numeric characters
func pwdIsValid(pwd string, pwdLen int) bool {
    result := false
    // make sure the first pwdLen characters are all alphanumeric
    i := IS_ALPHANUMERIC.FindSubmatchIndex([]byte(pwd))
    if i != nil && i[0] == 0 && i[1] >= pwdLen {
        // now make sure there is at least one (but no more than 5) uppercase, lowercase and numeric characters
        pwdPrefix := []byte(pwd[0:pwdLen])
        if HAS_UPPERCASE.Match(pwdPrefix) && HAS_LOWERCASE.Match(pwdPrefix) && HAS_NUMERICS.Match(pwdPrefix) {
            result = true
        }
    }
    return result
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
            result = getCandidatePwd(passphrase, username, hostname, 12, i)
            valid = pwdIsValid(result, passwordLength)
            i += 1
        }
        // success: display the result (the first passwordLength characters, with special chars at the end, if any)
        fmt.Print(fmt.Sprintf("Your password for %s logging in as user %s is:\n\n%s\n\n", hostname, username, strings.Join([]string{result[0:(passwordLength - len(specials))], specials}, "")))
    }
}
