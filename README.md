# MorseBot
Morsebot 2.0<br />

Setup for Dev<br />
Set env variable $MORSEBOT to discord key<br />
Set env variable $GOVEEDB to path of GOVEEDB<br />

To add a new module: <br />
Add a folder in pkg + import utils. Add command to command map. Follow structure for messagecreate commands (see echo for example)<br />
Event handler in events.go, add new event if want to use something other than messagecreate<br />

