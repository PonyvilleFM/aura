gophy
=====

[![Build Status](https://travis-ci.org/paddycarey/gophy.svg?branch=master)](https://travis-ci.org/paddycarey/gophy)
[![GoDoc](https://godoc.org/github.com/paddycarey/gophy?status.svg)](https://godoc.org/github.com/paddycarey/gophy)

gophy is a Golang wrapper for the [Giphy.com](http://www.giphy.com) API. gophy
is a WIP but aims to eventually have 100% coverage of the API. Gophy requires
Go >= 1.1.

See [GoDoc](https://godoc.org/github.com/paddycarey/gophy) for full library documentation.


### TODO

- Random endpoints
- Full documentation


### Example

Using gophy should be simple, just import it into your project, create a client and call the appropriate method. A contrived example is shown below:

```go
package main

import (
	"fmt"
	"github.com/paddycarey/gophy"
)

def main() {
	co := &gophy.ClientOptions{}
	client := gophy.NewClient(co)

	gifs, err := gophy.TrendingGifs("", 20)
	if err != nil {
		panic(err)
	}

	for _, gif := range gifs {
		fmt.Printf("%s: %s", gif.Id, gif.URL)
	}
}
```


### Copyright & License

- Copyright © 2015 Patrick Carey (https://github.com/paddycarey)
- Licensed under the **MIT** license.
