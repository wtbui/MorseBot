# MorseBot
Morsebot 2.0<br />
Current Features<br />
1. Multi user govee IOT light control and sync<br />

Setup for Dev<br />
Set env variable $MORSEBOT to discord key<br />
Set env variable $GOVEEDB to path of GOVEEDB<br />

To add a new module: <br />
1. Add a folder in pkg + import utils. <br />
2. Add command to command map in events/events.go. <br />
3. Follow structure for messagecreate commands (see echo for example)<br />
Event handler in events.go, add new event if want to use something other than messagecreate<br />

Build Commands <br />
"make" to build project <br />
"make clean" to delete
