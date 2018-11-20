# espoke - A prometheus blackbox probe for ES clusters

* Only works using a opinionated consul model as discovery method
* Expose prometheus metrics from a blank search on every ES indexes of every datanodes

```
Start ES discovering mecanism & probe periodically requests against all ES nodes.
Expose all measures using a prometheus compliant HTTP endpoint.

Usage:
  espoke serve [flags]

Flags:
  -h, --help   help for serve

Global Flags:
      --cleaningPeriod string   prometheus metrics cleaning interval (for vanished nodes) (default "600s")
  -c, --config string           config file (default is $HOME/.espoke.yaml)
  -a, --consulApi string        consul target api host:port (default "127.0.0.1:8500")
      --consulPeriod string     nodes discovery update interval (default "120s")
  -l, --loglevel string         log level (default "info")
  -p, --metricsPort int         port where prometheus will expose metrics to (default 2112)
      --probePeriod string      elasticsearch nodes probing interval (default "30s")
```
