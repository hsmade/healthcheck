package healthcheck

import (
	"fmt"
)

func Example() {
	h := New()

	fooHealth := h.Add("foo")

	// go log.Fatal(h.Serve(":8000")) // to start the web server

	fmt.Println(h.Status())
	fmt.Println(fooHealth.Status)

	fooHealth.Update(true)
	fmt.Println(h.Status())
	fmt.Println(fooHealth.Status)

	// Output:
	// false
	// false
	// true
	// true
}
