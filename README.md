Gophernotify
============

This program consists of two parts.
One is a server side (GAE) program and other one is a command line program.
First you access the server side program on GAE, the program provides a client id (number) as the last segment of URL.
Next you post messages using the client id from the command line program as followings:
	
	gophernotify -m "Hello, Gophers" -c 1

This example uses 1 as client id and posts "Hello, Gophers" to the GAE server.
And you can see the message on a browser without reload.
