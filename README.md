# Graphite-ng

Experimental version of a new generation Graphite API server in Golang, leveraging Go's efficient concurrency constructs.
Goals are: speed, ease of deployment. elegant code.  
Furthermore, this rewrite allows to redesign and fundamentally improve some specific annoyances.

# Limitations

 * Not all functions are supported yet
 * only the json output, not the png renderer. (because [client side rendering](https://github.com/vimeo/timeserieswidget/) is best)
 * No web UI (because there are plenty of graphite dashboards out there)
 * No reinventing carbon/whisper/ceres at this point. (I later want to hook this up to a reliable timeseries store, maybe whisper, ceres, kairosdb, ...).
   (there's a `carbon-es` dir which is a carbon-cache that stores metrics in elasticsearch, but I'm still experimenting with it)
 * No events system (graphite events sucks, [anthracite](https://github.com/Dieterbe/anthracite/) is better)
 * No wildcards yet

# How it works

`graphite-ng` is a webserver that gives you a `/render/` http endpoint where you can do queries like
`/render/sum(stats.web1,scale(stats.web2,5.2))from=123&until=456`

`graphite-ng` converts all user input into a real, functioning Go program, compiles and runs it, and returns the output.
It can do this because the graphite api notation can easily be converted to real program code.  Great power, great responsability.
The worker functions use channels to stream data around and avoid blocking.

# Installation & running

run this from the code checkout:
```
rm -f executor-*.go ; go install github.com/Dieterbe/graphite-ng && graphite-ng
```

then open something like this in your browser:

```
http://localhost:8080/render/?target=stats.web2&target=derivative(stats.web2)
http://localhost:8080/render/?target=sum(stats.web1,scale(stats.web2,5))&from=60&until=300
```

# Which metrics and functions can I use?

Look at data.go and the functions directory.
Since graphite-ng is not hooked up yet to a decent timeseries store, for now we'll have to do with the
examples in data.go.

# Function plugins 

all functions come in plugin files. want to add a new function? just drop a .go file in the functions folder and restart.  You can easily add your own functions
that get data from external sources, manipulate data, or represent data in a different way; and then call those functions from your target string.

# other interesting things & diff with real graphite:

* every function can request a different timerange from the functions it depends on.   E.g.:
  * `derivative` needs the datapoint from before the requested timerange
  * `movingAverage(foo, X)` needs x previous datapoints.
  Regular graphite doesn't support this so you end up with gaps in the beginning of the graph.
* clever automatic rollups based on tags (TODO)
* The keys in Graphite's json output are sometimes not exactly the requested target string (usually manifests itself as floats being rounded), it's not so easily fixed in Graphite
  due to the pathExpression system,  which means client renderes have to implement ugly hacks to work around this. 
  With graphite-ng we just use the exact same string.
* avoid any results being dependent on any particular potentially unknown variable, aim for per second instead of per current-interval, etc. specifically:
  * derivative is a true derivative (ie `(y2-y1)/(x2-x1)`) unlike graphite's derivative where you depend on a factor that depends on whatever the resolution is at each point in time.
