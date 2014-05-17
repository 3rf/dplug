DPlug
======

###Goals
I want to be able set up a dynamic plugin by implementing a simple interface and adding one line to a config file

###Design
DPlug aims to simplify the addition of dynamically "linked" plugins to your Go project.
The first iteration of this will involve moving the plugin code to a seperate main package, registering functions as "hooks", then starting up a server--all the gross stuff is done for you.
This package heavily utilizes the net/rpc package.

####TODO
* Nicer plugin shutdown (instead of killing the process outright)
* Real Unit and Integration Tests
* "RefreshConfig" function to let the client dynamically update plugins

