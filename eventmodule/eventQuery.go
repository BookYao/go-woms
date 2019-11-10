package  eventmodule

import (
	"fmt"
	"net/http"
)

func QuerySingleEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Single Event Query...")
}
