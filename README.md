go-one-password
===============

About
-----

This project is inspired by <a href="https://github.com/maxtaco/oneshallpass">oneshallpass</a> but is written in <a href="http://golang.org/">Go</a> instead of javascript, and runs on the command line as a self-contained, statically compiled binary, instead of a web browser.

The technical implementation is similar, i.e., HMAC-SHA512() hashing a combination of the host, username, generator and indicator numbers, but using <a href="http://www.tarsnap.com/scrypt.html">scrypt</a> instead of <a href="http://en.wikipedia.org/wiki/PBKDF2">PBKDF2</a> for generating the shared private key (dk) value from the <a href="https://en.wikipedia.org/wiki/Passphrase">passphrase</a>.

Building and Installing
-----------------------

This code requires Go's <a href="https://code.google.com/p/go/source/browse/scrypt/scrypt.go?repo=crypto">scrypto package</a> before building:

```
$ sudo go get code.google.com/p/go.crypto/scrypt
$ go build go-one-password.go
```

Update your $PATH to include the folder where go-one-password was built, and add a shorter alias, if you prefer (e.g., "g1p", assuming that doesn't conflict with anything on your system).

Usage
-----

The core idea is that by remembering just one <a href="https://en.wikipedia.org/wiki/Passphrase#Passphrase_selection">quality passphrase</a> (known _only_ by you), you can generate unique and secure passwords for multiple website logins.

There are <a href="http://www.queen.clara.net/pgp/pass.html">many</a> different <a href="https://en.wikipedia.org/wiki/Passphrase#Example_methods">ways</a> of selecting a quality passphrase, but if you cannot come up with one on your own, there are <a href="http://passphra.se/">several free sites</a> which can <a href="https://oneshallpass.com/pp.html">pick one for you</a>.

Once you settle on a passphrase, just *make sure you commit it to memory*; it's not stored anywhere by this code, and if lost or forgotten, is unrecoverable.

If you forget how to use go-one-password type it (or whatever alias you've used for it) in a shell prompt followed by "-help":

```
$ g1p -help
Usage of g1p:
  -host="": (required) the website you want to login to (e.g. "amazon.com")
  -spec="": (optional) if the website requires one or more "special" characters in the password (e.g., "#%*" etc.) specify one or more of them here
  -user="": (required) the username or email address you use to login
```

Usage Examples
--------------

Here's how to use it in practice (the passphrase is asked in an interactive prompt, instead of from a command line argument, because we don't want to save the passphrase in your shell history by accident):

```
$ g1p -host example.org -user me@example.com
What's your passphrase? (or ctrl-c to quit) close introduced when lunch
Your password for example.org logging in as user me@example.com is:

o95gZHxeh7D9LYnp

```

This is another example, for when the site requires one or more "special" characters:

```
$ g1p -host example.org -user me@example.com -spec="#%"
What's your passphrase? (or ctrl-c to quit) close introduced when lunch
Your password for example.org logging in as user me@example.com is:

o95gZHxe#%h7D9LYnp

```
