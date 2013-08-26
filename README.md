Early, experimental version of a "new generation" Graphite server in Golang.

I want to build something akin to Graphite's rich function API leveraging Go's efficient concurrency constructs.
Goals are: speed, ease of deployment. elegant code.

However:
 * only the json output, not the png renderer. (because [client side rendering](https://github.com/vimeo/timeserieswidget/) is best)
 * No web UI (because there are plenty of graphite dashboards out there)
 * No carbon/whisper/ceres alike at this point. (I later want to hook this up to a reliable timeseries store, maybe whisper, ceres, kairosdb, ...).
 * No events system (graphite events sucks, [anthracite](https://github.com/Dieterbe/anthracite/) is better)

So what this does is spawn a webserver that gives you a /render/ endpoint where you can do queries like
`/render/sum(stats.web1.bytes_received,scale(stats.web2.bytes_received,5))from=123&until=456`

It's not entirely working yet, but close.

Note that the program converts all user input into a real, functioning Go program, compiles and runs it, and runs the stdout.
It can do this because the graphite api notation can easily be converted to real program code.  Great power, great responsability.

to try it out, run this from the code checkout:
```
rm -f executor-*.go ; go install github.com/Dieterbe/graphite-ng && graphite-ng
```
