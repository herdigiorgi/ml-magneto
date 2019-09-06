# ML - MAGNETO

You can run this exercise locally by just typing:

```
make
```

You need Docker Compose installed in your Linux system.

## Level 1

The main function that calculates if the DNA sequence corresponds to a human or a mutant can be found in `src/ml/mutant/mutant.go` with the `IsMutant(input []string) (bool, error)` signature.

The original `input []string` is transformed into an NXM array of interconnected `Nodes`. Each node is connected with the `Node.nodes` variable. And to prevent re-testing combinations the `Node.dicarted` variable is used. The code has the concept of directions, and is a way of describing if the algorithm took the path to the right, left, up or down relative to the `x` element:

```
0  1  2
7  X  3
6  5  4
```

I'm sure there is some cool graph theory algorithm that could greatly simplify what I have done. If you know how to improve it, tell me.

## Level 2

`src/ml/main.go` has the standard go webserver. If you run `make` you should be able to interact with the API:

- `/mutant`
- `/stats` (it updates every 5seconds)

Each endpoint sets `Cache-Control`. It may be useful if the server is behind a reverse proxy, or maybe a CDN. In this way.

## Level 3

`src/stats/stats.go` stores the results of the DNA checks in an SQLite database. It uses the `getCachedStats` function to cache results. I created that "eventually consistent" cache function because I thought that was what I was expected to do, so that constant `COUNT`s queries are avoided.

## 100 to 1M QPS

It's a little bit hard to assume this program could support 1M qps, but because the `Cache-Control` headers are set, a CDN like AWS CloudFront could be configured to pass `POST` requests while honoring the `Cache-Control`.