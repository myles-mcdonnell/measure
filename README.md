## Measure

[![Go Report Card](https://goreportcard.com/badge/github.com/myles-mcdonnell/measure)](https://goreportcard.com/report/github.com/myles-mcdonnell/measure)
 
cli utility to measure average response time of a get request.

Args

* url 
* concurrency: maximum concurrent requests, default 5
* averageWindow: the average period in seconds, so only the requests within that window from the current time will be included

eg.

```bash
$ measure --url https://tour.golang.org --concurrency 20 --averageWindow 5

Total Reqs: 638 - Window Reqs: 627 - Window Average (ms): 136^
```