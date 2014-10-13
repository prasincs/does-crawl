# Does Crawl 

_Because DO Crawl didn't sound quite right but YMMV_

## How to Build and run

The build script included should handle the process.

`./build && ./bin/does-crawl` should work well.

It should start a web server in port 9999, but an optional `--port` parameter can be supplied that allows you to change the port number.


## Credits

I took pieces of the code from the following locations, mostly to avoid the cold-start/writer's block problem

* Martini Framework https://github.com/PuerkitoBio/martini-api-example
* Crawler with Multiplex https://gist.github.com/dyoo/6064879 

And various library maintainers that I have pulled dependencies from.

## Caveats

This is my first from-scratch-ish golang project built under heavy time constraints (and context-switching, writing Go and Clojure at the same time is not recommended). I'll probably hate the code in a few days/weeks/months once I understand more things about organizing go code.