name: App
role: behavior
intent: production entry point - timeouts and graceful shutdown
behavior:
  - build a Dispatcher with NewDispatcher(4, 1024); Start(); build a Server
  - construct an *http.Server with Addr ":8080", Handler server.Routes(), and ReadHeaderTimeout,
    ReadTimeout, WriteTimeout, IdleTimeout all set (a bare ListenAndServe has none - slowloris risk)
  - use signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM); run
    srv.ListenAndServe() in a goroutine (ignore http.ErrServerClosed); on ctx.Done() call srv.Shutdown
    with a timeout context, then dispatcher.Stop() so queued events drain
constraints: standard library (net/http, context, os/signal, syscall, time, log, errors); package main; func main only here
