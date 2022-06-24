package tui

var topics = map[pageName]string{
	pgHelp: "Context help shows the relevant topic for the current screen." +
		" Use arrows  [$ptc]↑[-]  and  [$ptc]↓[-]  to scroll the text.  Press  [$ptc]ESC[-]  to close" +
		" the help window and return to the main screen.",
	pgLoginMenu: "Use arrows  [$ptc]↑[-]  and  [$ptc]↓[-]  to navigate, press [$ptc]Enter[-] to choose the item.\n\n" +
		"[$ptc]Login With Browser[-]\n\n   Opens a browser with the slack workspace of your choice" +
		" and waits till Slack Login is done.\n\n" +
		"[$ptc]Login with Token and Cookie[-]\n\n   Allows one to specify the token and cookie or cookie file.",
}
