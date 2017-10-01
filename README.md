# miles 

[![CircleCI](https://circleci.com/gh/tmatias/miles.svg?style=shield&)](https://circleci.com/gh/tmatias/miles) [![travis-ci](https://travis-ci.org/tmatias/miles.svg?&branch=master)](https://travis-ci.org/tmatias/miles) [![Documentation](https://godoc.org/github.com/tmatias/miles?status.svg)](http://godoc.org/github.com/tmatias/miles) [![Go Report Card](https://goreportcard.com/badge/github.com/tmatias/miles)](https://goreportcard.com/report/github.com/tmatias/miles)

miles may be the start of a library to ask users of CLI apps to choose from a list of options.

__It's just me toying with Go and learning how to setup a project, though, so, don't expect much :).__

## usage

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tmatias/miles"
)

func main() {
	opt := miles.Options{
		Prompt:  "What's your preferred color?",
		Allowed: []string{"b", "g"},
	}

	ans, err := opt.Choose()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("You choose %s\n", ans)
}
```

Gives you:

``` console
> What's your preferred color? [b/g]: g
> You choose g
```
