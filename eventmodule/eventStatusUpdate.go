package eventmodule

import (
	"fmt"
	"net/http"
)

func EventStatusUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Event Status Update...")
}