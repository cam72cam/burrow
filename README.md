Burrow
======

A alternate interface to the [Delve](https://github.com/derekparker/delve) Go debugger


Burrow has 2 modes for entering commands, command mode and key mode.

###Command Mode:
:run [timestep] (must be at breakpoint)

:go 10 (same level, step 10 times)
:break file:number
:unbreak file:number
:breaks

:threads
:print $var
:scope


###Key Mode:
i - next in
o - next over
k - keep going same level

r - run
p - pause

