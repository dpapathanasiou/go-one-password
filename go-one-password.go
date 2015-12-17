/* go-one-password.go

   This is a Go library capable of generating a unique password for
   different sites and usernames based on a single, private (i.e.,
   known only to the person running this program) passphrase

*/

package onepassword

import (
	"bufio"
	"bytes"
	"code.google.com/p/go.crypto/scrypt"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
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

// GetCandidatePwd computes a password string based on the passphrase, hostname, username, a generation number, and an iteration number
func GetCandidatePwd(pp, usr, host string, g, i int) string {
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

// PwdIsValid checks the string created by getCandidatePwd and returns a boolean depending on whether or not it meets the criteria:
// the first pwdLen characters are all alphanumeric, and there is at least one (but no more than 5) uppercase, lowercase and numeric characters
func PwdIsValid(pwd string, pwdLen int) bool {
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
