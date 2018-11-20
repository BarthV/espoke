# espoke - A prometheus blackbox probe for ES clusters

## Infos

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

## Metrics

```
# HELP es_datanode_availability Reflects datanode availabity : 1 is OK, 0 means node unavailable 
# TYPE es_datanode_availability gauge
es_datanode_availability{cluster="test",nodename="elasticsearch-test01.preprod"} 1
es_datanode_availability{cluster="test",nodename="elasticsearch-test02.preprod"} 1
es_datanode_availability{cluster="test",nodename="elasticsearch-test03.preprod"} 0

# HELP es_datanode_count Reports current discovered nodes amount
# TYPE es_datanode_count gauge
es_datanode_count 3

# HELP es_datanode_search_latency Measure latency for every datanode (quantiles - in ns)
# TYPE es_datanode_search_latency summary
es_datanode_search_latency{cluster="test",nodename="elasticsearch-test01.preprod",quantile="0.5"} 3.5991948e+07
es_datanode_search_latency{cluster="test",nodename="elasticsearch-test01.preprod",quantile="0.9"} 3.9996552e+07
es_datanode_search_latency{cluster="test",nodename="elasticsearch-test01.preprod",quantile="0.99"} 6.1367425e+07
es_datanode_search_latency_sum{cluster="test",nodename="elasticsearch-test01.preprod"} 1.1059591307e+10
es_datanode_search_latency_count{cluster="test",nodename="elasticsearch-test01.preprod"} 331

es_datanode_search_latency{cluster="test",nodename="elasticsearch-test02-par.preprod",quantile="0.5"} 3.4431259e+07
es_datanode_search_latency{cluster="test",nodename="elasticsearch-test02-par.preprod",quantile="0.9"} 4.1628496e+07
es_datanode_search_latency{cluster="test",nodename="elasticsearch-test02-par.preprod",quantile="0.99"} 6.1737103e+07
es_datanode_search_latency_sum{cluster="test",nodename="elasticsearch-test02-par.preprod"} 1.0730152625e+10
es_datanode_search_latency_count{cluster="test",nodename="elasticsearch-test02-par.preprod"} 331

# HELP es_probe_consul_discovery_duration Time spent for discovering nodes using Consul API (in ns)
# TYPE es_probe_consul_discovery_duration summary
es_probe_consul_discovery_duration{quantile="0.5"} 1.1310931e+07
es_probe_consul_discovery_duration{quantile="0.9"} 1.2408547e+07
es_probe_consul_discovery_duration{quantile="0.99"} 2.276767e+07
es_probe_consul_discovery_duration_sum 1.676217694e+09
es_probe_consul_discovery_duration_count 111

# HELP es_probe_errors_count Reports Espoke internal errors absolute counter since start
# TYPE es_probe_errors_count counter
es_probe_errors_count 331

# HELP es_probe_metrics_cleaning_duration Time spent for cleaning vanished nodes metrics (in ns)
# TYPE es_probe_metrics_cleaning_duration summary
es_probe_metrics_cleaning_duration{quantile="0.5"} 32429
es_probe_metrics_cleaning_duration{quantile="0.9"} 67287
es_probe_metrics_cleaning_duration{quantile="0.99"} 70555
es_probe_metrics_cleaning_duration_sum 960260
es_probe_metrics_cleaning_duration_count 27
```
