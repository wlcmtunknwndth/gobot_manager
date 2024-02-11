package telegram

// import "errors"

const msgHelp = `I can save and keep your pages. Also i'm able to offer u them randomly.

To save the page, send me a link of it

To get a random saved page, send me a command /rnd

Beware! should i delete it?
`

const msgHello = "Hello there! \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command"
	msgNoSavedPages   = "You havent saved a page yet"
	msgSaved          = "Succesfully saved!"
	msgAlreadyExists  = "You have already saved this one"
	msgNoSavedLinks   = "Sorry, there are no links to this time"
	msgNotEnoughArg   = "Not enough arguments"
	msgNoPerm         = "Not enough permissions"
)
