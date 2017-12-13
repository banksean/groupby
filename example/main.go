package main

import (
	"fmt"
	"github.com/banksean/groupby"
)

type typ struct {
	A, B string
	C    int
}

func main() {

	// Use groupby.Field to group values by typ.C:
	a := []typ{{"a", "b", 0}, {"a", "c", 1}, {"b", "b", 1}, {"b", "d", 3}, {"c", "x", 3}}
	fmt.Printf("%+v\n", groupby.Field(a, "C"))

	// Use groupby.Func to group by a value derived from typ.C
	b := groupby.Func(a, func(i interface{}) interface{} {
		return i.(typ).C * 2
	})
	fmt.Printf("%+v\n", b)

	// Use groupby.Chan to group by multiple values derived from typ.C
	c := groupby.Chan(a, func(i interface{}) chan interface{} {
		c := make(chan interface{})
		go func() {
			c <- i.(typ).C
			c <- i.(typ).C + 10
			close(c)
		}()
		return c
	})

	fmt.Printf("%+v\n", c)
}
