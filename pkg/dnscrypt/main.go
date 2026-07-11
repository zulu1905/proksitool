package dnscrypt

import (
    "context"
    "flag"
    "fmt"
    "os"
    "runtime"
	"github.com/jedisct1/dlog"
)

const AppVersion = "umbrella-1.0"
const DefaultConfigFileName = "dnscrypt-proxy.toml"

type App struct {
    flags *ConfigFlags
    proxy *Proxy
    quit  chan os.Signal
}

var dnsCancel context.CancelFunc

func Run() {
    ctx, cancel := context.WithCancel(context.Background())
    dnsCancel = cancel

    tzErr := TimezoneSetup()
    dlog.Init("dnscrypt-proxy", dlog.SeverityNotice, "DAEMON")
    if tzErr != nil {
        dlog.Warnf("Timezone setup failed: [%v]", tzErr)
    }
    runtime.MemProfileRate = 0

    version := flag.Bool("version", false, "print current proxy version")

    flags := ConfigFlags{}
    flags.Resolve = flag.String("resolve", "", "resolve a DNS name (string can be <name> or <name>,<resolver address>)")
    flags.List = flag.Bool("list", false, "print the list of available resolvers for the enabled filters")
    flags.ListAll = flag.Bool("list-all", false, "print the complete list of available resolvers, ignoring filters")
    flags.IncludeRelays = flag.Bool("include-relays", false, "include the list of available relays in the output of -list and -list-all")
    flags.JSONOutput = flag.Bool("json", false, "output list as JSON")
    flags.Check = flag.Bool("check", false, "check the configuration file and exit")
    flags.ConfigFile = flag.String("config", DefaultConfigFileName, "Path to the configuration file")
    flags.Child = flag.Bool("child", false, "Invokes program as a child process")
    flags.NetprobeTimeoutOverride = flag.Int("netprobe-timeout", 60, "Override the netprobe timeout")
    flags.ShowCerts = flag.Bool("show-certs", false, "print DoH certificate chain hashes")

    flag.Parse()

    if *version {
        fmt.Println(AppVersion)
        os.Exit(0)
    }

    if fullexecpath, err := os.Executable(); err == nil {
        WarnIfMaybeWritableByOtherUsers(fullexecpath)
    }

    app := &App{
        flags: &flags,
    }

    app.proxy = NewProxy()
    _ = ServiceManagerStartNotify()

    // Başlat
 //   go app.AppMain()

    // Çıkış bekleme
    select {
    case <-app.quit:
        dlog.Notice("Quit signal received...")
    case <-ctx.Done():
        dlog.Notice("Context cancelled, shutting down...")
    }
}

func Stop() {
    if dnsCancel != nil {
        dnsCancel()
        dnsCancel = nil
    }
}
