Burrow
======

A alternate interface to the [Delve](https://github.com/derekparker/delve) Go debugger


Burrow has 2 modes for entering commands, command mode and key mode.

###Command Mode:
* `:run $interval` - Steps through your program one line at a time, pausing for a specified interval in between.  Requires breakpoint

* `:go $steps` - Runs $steps from the current break point in the program.  Requires breakpoint


* `:break $file:$number` - Sets a breakpoint at $file line $number.

* `:unbreak file:number` - Removes breakpoint at $file line $number if it exists.

* `:breaks` - Shows all of the current breakpoints in the progam.


* `:threads [-v]` - Shows the status of all traced threads.  Verbose flag to show file context

* `:routines [-v]` - Shows the position of all gorountines.  Verbose flag to show file context

* `:print $var` - Show value of $var.  Requires breakpoint

* `:scope` - Show all variables in the current scope.  Requires breakpoint


###Key Mode:
* i - Step and try to step into the next scope

* o - Step and try to step out of the current scope

* k - Step to the next line in the current scope

* r - run

* p - pause

