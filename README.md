go-one-password
===============

About
-----

This project is inspired by <a href="https://github.com/maxtaco/oneshallpass">oneshallpass</a> but is written in <a href="http://golang.org/">Go</a> instead of javascript, and runs as a self-contained, statically compiled binary, either on the command line or as a gui, instead of a web browser.

The technical implementation is similar, i.e., HMAC-SHA512() hashing a combination of the host, username, generator and indicator numbers, but using <a href="http://www.tarsnap.com/scrypt.html">scrypt</a> instead of <a href="http://en.wikipedia.org/wiki/PBKDF2">PBKDF2</a> for generating the shared private key (dk) value from the <a href="https://en.wikipedia.org/wiki/Passphrase">passphrase</a>.

Building and Installing
-----------------------

This program now comes in two versions, a command line interface (cli), and a graphical user interface (gui).

In order to build the gui version, you will need this library (optional):

```sh
$ go get github.com/mattn/go-gtk/gtk
```

Note that [go-gtk](http://mattn.github.io/go-gtk/) requires that the [GTK-Development packages](https://github.com/mattn/go-gtk#install) for your system are already installed.

Next, import these two repositories (required):

```sh
$ go get github.com/howeyc/gopass
$ go get github.com/dpapathanasiou/go-one-password
```

Use the [Makefile](Makefile) to build either or both versions:

```sh
$ make all # build both the cli and gui versions
$ make cli # build just the cli version
$ make gui # build just the gui version
```

### Command Line Interface Version

The resulting binary is <tt>go-one-password-cl</tt>.

Update your $PATH to include the folder where go-one-password-cl was built, and add a shorter alias, if you prefer (e.g., "g1p", assuming that doesn't conflict with anything on your system).

### Graphical User Interface Version

The resulting binary is <tt>go-one-password-ui</tt>.

You can add a launcher from your desktop menu to run it that way, if you prefer.

Usage
-----

The core idea is that by remembering just one <a href="https://en.wikipedia.org/wiki/Passphrase#Passphrase_selection">quality passphrase</a> (known _only_ by you), you can generate unique and secure passwords for multiple website logins.

There are <a href="http://www.queen.clara.net/pgp/pass.html">many</a> different <a href="https://en.wikipedia.org/wiki/Passphrase#Example_methods">ways</a> of selecting a quality passphrase, but if you cannot come up with one on your own, there are <a href="http://passphra.se/">several free sites</a> which can <a href="https://oneshallpass.com/pp.html">pick one for you</a>.

Once you settle on a passphrase, just *make sure you commit it to memory*; it's not stored anywhere by this code, and if lost or forgotten, is unrecoverable.

### Command Line Interface Version

If you forget how to use go-one-password-cl type it (or whatever alias you've used for it) in a shell prompt followed by "-help":

```
$ ./go-one-password-cl -help
Usage of g1p:
  -host="": (required) the website you want to login to (e.g. "amazon.com")
  -plen=16: (optional) set the resulting password length (the default is 16)
  -spec="": (optional) if the website requires one or more "special" characters in the password (e.g., "#%*" etc.) specify one or more of them here
  -user="": (required) the username or email address you use to login
```

#### Usage Examples

Here's how to use it in practice (the passphrase is asked in an interactive prompt, instead of from a command line argument, because we don't want to save the passphrase in your shell history by accident).

Note that while the passphrase is hidden on Mac OSX, Windows and Linux systems, it may appear as viewable text on other operating systems. To keep the passphrase text hidden on such systems, use the gui version instead.

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

o95gZHxeh7D9LY#%
```

### Graphical User Interface Version

The gui version supports all the same features of the cli version, with the additional benefit that it hides the passphrase by default:

![](http://i.imgur.com/FAPKYtm.png "Graphical User Interface Version")

Passphrases can be made visible if desired, and "special" characters work too:

![](http://i.imgur.com/cBqIWHi.png "Graphical User Interface Version")

Launching the gui version with any of these command-line switches automatically pre-populates the corresponding fields in the input:

```
$ ./go-one-password-ui -help
Usage of ./go-one-password-ui:
  -host string
    	the website you want to login to (e.g. "amazon.com")
  -plen string
    	set the resulting password length (the default is 16) (default "16")
  -spec string
    	if the website requires one or more "special" characters in the password (e.g., "#%*" etc.) specify one or more of them here
  -user string
    	the username or email address you use to login
```

Using these command-line switches to launch the gui is optional.
