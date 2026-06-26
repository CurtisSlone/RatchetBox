name: App
role: behavior
intent: the program entry point - wire the dispatcher and server together and serve
behavior:
  - in func main: create a Dispatcher with NewDispatcher(4, 1024) and call Start()
  - create a Server with NewServer(dispatcher)
  - log that it is listening, then http.ListenAndServe(":8080", server.Routes())
  - if ListenAndServe returns an error, log.Fatal it
constraints: standard library only (net/http, log); package main; this file (main.go) is the ONLY file
  with func main; uses the existing Dispatcher and Server API verbatim
