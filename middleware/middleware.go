package middleware

import (
	"github.com/ahmedkhaeld/booking/data"
	"github.com/ahmedkhaeld/jazz"
)

type Middleware struct {
	*jazz.Jazz
	Models data.Models
}
