Early, experimental version of a "new generation" Graphite server in Golang.

I want to build something akin to Graphite's rich function API leveraging Go's efficient concurrency constructs.
Goals are: speed, ease of deployment. elegant code.

However:
 * only the json output, not the png renderer. (because [client side rendering](https://github.com/vimeo/timeserieswidget/) is best)
 * No web UI (because there are plenty of graphite dashboards out there)
 * No reinventing carbon/whisper/ceres at this point. (I later want to hook this up to a reliable timeseries store, maybe whisper, ceres, kairosdb, ...).
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
then open something like this in your browser:
```
http://localhost:8080/render/?target=sum(stats.web1.bytes_received,scale(stats.web2.bytes_received,5))&from=60&until=300
```
look at data.go and functions.go for which metrics and functions you can use so far.


interesting things & diff with real graphite:
* consistently treat datapoint as the value covering the timespan leading up to it, this matters esp. for derivative, integral, etc
* make functions that need extra info outside of the from-until range (i.e. derivative needs the from-60 datapoint; movingAverage needs x previous datapoints, etc)
  able to get that info in an elegant way. unlike graphite where sometimes the beginning of your graph is empty because a movingAverage only has enough data halfway the graph.
* clever automatic rollups based on tags (TODO)
* The keys in Graphite's json output are sometimes not exactly the requested target string (usually manifests itself as floats being rounded), it's not so easily fixed in Graphite
  due to the pathExpression system,  which means client renderes have to implement ugly hacks to work around this.  With graphite-ng we just use the exact same string.
* it should be easy to add your own functions, by loading them all as plugins (TODO)
