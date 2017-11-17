/* go-one-password-ui.go

   Wraps a simple GTK user interface around the go-one-password library

*/

package main

import (
	"flag"
	"fmt"
	"github.com/dpapathanasiou/go-one-password/onepassword"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"strconv"
	"strings"
)

const (
	INPUT_DEFAULT = ""
	SPEC_DEFAULT  = "16"
)

func clearResult(win *gtk.TextView) {
	var start, end gtk.TextIter
	buffer := win.GetBuffer()
	buffer.GetBounds(&start, &end)
	buffer.Delete(&start, &end)
}

func setResult(win *gtk.TextView, tag *gtk.TextTag, msg string) {
	var start, end gtk.TextIter
	buffer := win.GetBuffer()
	buffer.GetStartIter(&start)
	buffer.Insert(&start, msg)
	buffer.GetBounds(&start, &end)
	buffer.ApplyTag(tag, &start, &end)
}

func main() {
	// provide an option to pre-fill the UI field inputs based on command line switches
	var hostCL, userCL, specCL, pwdLenCL string
	flag.StringVar(&hostCL, "host", INPUT_DEFAULT, "the website you want to login to (e.g. \"amazon.com\")")
	flag.StringVar(&userCL, "user", INPUT_DEFAULT, "the username or email address you use to login")
	flag.StringVar(&specCL, "spec", INPUT_DEFAULT, "if the website requires one or more \"special\" characters in the password (e.g., \"#%*\" etc.) specify one or more of them here")
	flag.StringVar(&pwdLenCL, "plen", SPEC_DEFAULT, fmt.Sprintf("set the resulting password length (the default is %s)", SPEC_DEFAULT))
	flag.Parse()

	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("go-one-password")
	window.SetIconName("dialog-password")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	})

	vbox := gtk.NewVBox(false, 1)
	vpaned := gtk.NewVPaned()
	vbox.Add(vpaned)

	// credential input frame
	credframe := gtk.NewFrame("Credentials")
	credbox := gtk.NewVBox(false, 2)
	credframe.Add(credbox)

	// results frame
	resframe := gtk.NewFrame("")
	resbox := gtk.NewVBox(false, 2)
	resframe.Add(resbox)

	vpaned.Pack1(credframe, false, false)
	vpaned.Pack2(resframe, false, false)

	// credentials input
	hostbox := gtk.NewHBox(true, 1)
	hostlabel := gtk.NewLabel("Site Name")
	hostlabel.SetJustify(gtk.JUSTIFY_RIGHT)
	hostname := gtk.NewEntry()
	hostname.SetText(hostCL)
	hostbox.Add(hostlabel)
	hostbox.Add(hostname)
	credbox.PackStart(hostbox, false, false, 2)

	userbox := gtk.NewHBox(true, 1)
	userlabel := gtk.NewLabel("Username")
	userlabel.SetJustify(gtk.JUSTIFY_RIGHT)
	username := gtk.NewEntry()
	username.SetText(userCL)
	userbox.Add(userlabel)
	userbox.Add(username)
	credbox.PackStart(userbox, false, false, 2)

	passbox := gtk.NewHBox(true, 1)
	passlabel := gtk.NewLabel("Passphrase")
	passlabel.SetJustify(gtk.JUSTIFY_RIGHT)
	passname := gtk.NewEntry()
	passname.SetVisibility(false)
	passbox.Add(passlabel)
	passbox.Add(passname)
	credbox.PackStart(passbox, false, false, 2)

	visibox := gtk.NewHBox(true, 0)
	visilabel := gtk.NewLabel("")
	checkbutton := gtk.NewCheckButtonWithLabel("Show Passphrase")
	checkbutton.Connect("toggled", func() {
		if checkbutton.GetActive() {
			passname.SetVisibility(true)
		} else {
			passname.SetVisibility(false)
		}
	})
	visibox.Add(visilabel)
	visibox.Add(checkbutton)
	credbox.PackStart(visibox, false, false, 0)

	lenbox := gtk.NewHBox(true, 2)
	lenlabel := gtk.NewLabel("Length")
	lenlabel.SetJustify(gtk.JUSTIFY_RIGHT)
	lenname := gtk.NewEntry()
	lenname.SetText(pwdLenCL)
	lenbox.Add(lenlabel)
	lenbox.Add(lenname)
	credbox.PackStart(lenbox, false, false, 2)

	specialbox := gtk.NewHBox(true, 2)
	speclabel := gtk.NewLabel("Special Chars")
	speclabel.SetJustify(gtk.JUSTIFY_RIGHT)
	specname := gtk.NewEntry()
	specname.SetText(specCL)
	specialbox.Add(speclabel)
	specialbox.Add(specname)
	credbox.PackStart(specialbox, false, false, 2)

	// results window
	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.SetShadowType(gtk.SHADOW_IN)
	textview := gtk.NewTextView()
	buffer := textview.GetBuffer()
	highlight := buffer.CreateTag("highlighted", map[string]string{
		"background": "#FFFF99", "weight": "bold"})
	swin.Add(textview)
	resbox.Add(swin)

	// action buttons
	buttons := gtk.NewHBox(false, 1)
	generate := gtk.NewButtonWithLabel("Generate")
	generate.Clicked(func() {
		size, err := strconv.Atoi(lenname.GetText())
		if err != nil || size < 6 {
			badlen := gtk.NewMessageDialog(
				generate.GetTopLevelAsWindow(),
				gtk.DIALOG_MODAL,
				gtk.MESSAGE_ERROR,
				gtk.BUTTONS_OK,
				"Please use a positive number greater than or equal to five (5)")
			badlen.Response(func() {
				lenname.SetText(SPEC_DEFAULT)
				badlen.Destroy()
			})
			badlen.Run()
		} else {
			host := hostname.GetText()
			user := username.GetText()
			pass := passname.GetText()
			spec := specname.GetText()

			i := 0
			valid := false
			password := ""
			// keep generating passwords (using an updated iteration number)
			// until we get one that meets the PwdIsValid() criteria
			for !valid {
				password = onepassword.GetCandidatePwd(pass, user, host, 12, i)
				valid = onepassword.PwdIsValid(password, size)
				i += 1
			}

			clearResult(textview)
			setResult(textview, highlight, strings.Join([]string{password[0:(size - len(spec))], spec}, ""))
		}
	})
	buttons.Add(generate)

	clear := gtk.NewButtonWithLabel("Clear")
	clear.Clicked(func() {
		hostname.SetText("")
		username.SetText("")
		passname.SetText("")
		lenname.SetText(SPEC_DEFAULT)
		clearResult(textview)
	})
	buttons.Add(clear)

	credbox.Add(buttons)

	// start it up
	window.Add(vbox)
	window.ShowAll()
	gtk.Main()
}
