package groupby

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type T struct {
	A, B string
	C    int
	D    []int
}

func TestField(t *testing.T) {
	Convey("value", t, func() {
		a := []T{{"a", "b", 0, nil}, {"a", "c", 1, nil}, {"b", "b", 1, nil}, {"b", "d", 3, nil}, {"c", "x", 3, nil}}
		grouped := Field(a, "B")
		So(grouped, ShouldNotBeEmpty)
		So(len(grouped), ShouldEqual, 4)
		So(grouped["b"], ShouldResemble, []interface{}{
			a[0], a[2],
		})
		So(grouped["c"], ShouldResemble, []interface{}{
			a[1],
		})
		So(grouped["d"], ShouldResemble, []interface{}{
			a[3],
		})
		So(grouped["x"], ShouldResemble, []interface{}{
			a[4],
		})
	})

	Convey("pointer", t, func() {
		a := []*T{{"a", "b", 0, nil}, {"a", "c", 1, nil}, {"b", "b", 1, nil}, {"b", "d", 3, nil}, {"c", "x", 3, nil}}
		grouped := Field(a, "B")
		So(grouped, ShouldNotBeEmpty)
		So(len(grouped), ShouldEqual, 4)
		So(grouped["b"], ShouldResemble, []interface{}{
			a[0], a[2],
		})
		So(grouped["c"], ShouldResemble, []interface{}{
			a[1],
		})
		So(grouped["d"], ShouldResemble, []interface{}{
			a[3],
		})
		So(grouped["x"], ShouldResemble, []interface{}{
			a[4],
		})
	})
}

func TestFunc(t *testing.T) {
	Convey("value", t, func() {
		a := []T{{"a", "b", 0, nil}, {"a", "c", 1, nil}, {"b", "b", 1, nil}, {"b", "d", 3, nil}, {"c", "x", 3, nil}}
		groupFn := func(i interface{}) interface{} {
			return i.(T).C * 2
		}

		grouped := Func(a, groupFn)

		So(grouped, ShouldNotBeEmpty)
		So(len(grouped), ShouldEqual, 3)
		So(grouped[0], ShouldResemble, []interface{}{
			a[0],
		})
		So(grouped[2], ShouldResemble, []interface{}{
			a[1],
			a[2],
		})
		So(grouped[6], ShouldResemble, []interface{}{
			a[3],
			a[4],
		})
	})

	Convey("pointer", t, func() {
		a := []*T{{"a", "b", 0, nil}, {"a", "c", 1, nil}, {"b", "b", 1, nil}, {"b", "d", 3, nil}, {"c", "x", 3, nil}}

		groupFn := func(i interface{}) interface{} {
			// This seems odd. Intuitively, one would expect it to be *T, but it's T.
			// TODO: investigate, but this may not be an issue in practice.
			return i.(T).C * 2
		}

		grouped := Func(a, groupFn)

		So(grouped, ShouldNotBeEmpty)
		So(len(grouped), ShouldEqual, 3)
		So(grouped[0], ShouldResemble, []interface{}{
			a[0],
		})
		So(grouped[2], ShouldResemble, []interface{}{
			a[1],
			a[2],
		})
		So(grouped[6], ShouldResemble, []interface{}{
			a[3],
			a[4],
		})
	})
}

func TestChan(t *testing.T) {
	Convey("value", t, func() {
		a := []T{
			{"a", "b", 0, []int{0, 4}},
			{"a", "c", 1, []int{0, 1}},
			{"b", "b", 1, []int{0, 1, 2}},
			{"b", "d", 3, []int{1, 2, 3}},
			{"c", "x", 3, []int{1, 2, 3, 4}},
			{"z", "z", 0, nil},
		}

		groupFn := func(i interface{}) chan interface{} {
			c := make(chan interface{})
			go func() {
				for _, v := range i.(T).D {
					c <- v
				}
				close(c)
			}()
			return c
		}

		grouped := Chan(a, groupFn)

		So(grouped, ShouldNotBeEmpty)
		So(len(grouped), ShouldEqual, 5)
		So(grouped[0], ShouldResemble, []interface{}{
			a[0], a[1], a[2],
		})
		So(grouped[1], ShouldResemble, []interface{}{
			a[1], a[2], a[3], a[4],
		})
		So(grouped[2], ShouldResemble, []interface{}{
			a[2], a[3], a[4],
		})
		So(grouped[3], ShouldResemble, []interface{}{
			a[3], a[4],
		})
		So(grouped[4], ShouldResemble, []interface{}{
			a[0], a[4],
		})
	})

	Convey("pointer", t, func() {
		a := []*T{
			{"a", "b", 0, []int{0, 4}},
			{"a", "c", 1, []int{0, 1}},
			{"b", "b", 1, []int{0, 1, 2}},
			{"b", "d", 3, []int{1, 2, 3}},
			{"c", "x", 3, []int{1, 2, 3, 4}},
			{"z", "z", 0, nil},
		}

		groupFn := func(i interface{}) chan interface{} {
			c := make(chan interface{})
			go func() {
				for _, v := range i.(T).D {
					c <- v
				}
				close(c)
			}()
			return c
		}

		grouped := Chan(a, groupFn)

		So(grouped, ShouldNotBeEmpty)
		So(len(grouped), ShouldEqual, 5)
		So(grouped[0], ShouldResemble, []interface{}{
			a[0], a[1], a[2],
		})
		So(grouped[1], ShouldResemble, []interface{}{
			a[1], a[2], a[3], a[4],
		})
		So(grouped[2], ShouldResemble, []interface{}{
			a[2], a[3], a[4],
		})
		So(grouped[3], ShouldResemble, []interface{}{
			a[3], a[4],
		})
		So(grouped[4], ShouldResemble, []interface{}{
			a[0], a[4],
		})
	})
}

func ExampleField() {
	a := []T{{"a", "b", 0, nil}, {"a", "c", 1, nil}, {"b", "b", 1, nil}, {"b", "d", 3, nil}, {"c", "x", 3, nil}}

	grouped := Field(a, "B")
	fmt.Printf("%+v", grouped["b"])
	// Output: [{A:a B:b C:0 D:[]} {A:b B:b C:1 D:[]}]
}

func ExampleFunc() {
	a := []T{{"a", "b", 0, nil}, {"a", "c", 1, nil}, {"b", "b", 1, nil}, {"b", "d", 3, nil}, {"c", "x", 3, nil}}
	groupFn := func(i interface{}) interface{} {
		return i.(T).C * 2
	}

	grouped := Func(a, groupFn)
	fmt.Printf("%+v", grouped[6])
	// Output: [{A:b B:d C:3 D:[]} {A:c B:x C:3 D:[]}]
}

func ExampleChan() {
	a := []T{
		{"a", "b", 0, []int{0, 4}},
		{"a", "c", 1, []int{0, 1}},
		{"b", "b", 1, []int{0, 1, 2}},
		{"b", "d", 3, []int{1, 2, 3}},
		{"c", "x", 3, []int{1, 2, 3, 4}},
		{"z", "z", 0, nil},
	}

	groupFn := func(i interface{}) chan interface{} {
		c := make(chan interface{})
		go func() {
			for _, v := range i.(T).D {
				c <- v
			}
			close(c)
		}()
		return c
	}

	grouped := Chan(a, groupFn)
	fmt.Printf("%+v", grouped[3])
	// Output: [{A:b B:d C:3 D:[1 2 3]} {A:c B:x C:3 D:[1 2 3 4]}]
}
