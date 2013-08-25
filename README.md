Early, experimental version of a "new generation" Graphite server in Golang.

I want to build something akin to Graphite's rich function API leveraging Go's efficient concurrency constructs.
Goals are: speed, ease of deployment. elegant code.

However:
 * only the json output, not the png renderer. (because [client side rendering](https://github.com/vimeo/timeserieswidget/) is best)
 * No web UI (because there are plenty of graphite dashboards out there)
 * No carbon/whisper/ceres alike at this point. (I later want to hook this up to a reliable timeseries store, maybe whisper, ceres, kairosdb, ...).
 * No events system (graphite events sucks, [anthracite](https://github.com/Dieterbe/anthracite/) is better)

Ideas of converting commands to metrics processing:
* translate user input to go code directly, compile and execute: small code and very powerfull. but hard to do validation
* use a Go scanner/lexer: harder and verboser, but safer.

