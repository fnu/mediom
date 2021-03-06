package app

import (
	"fmt"
	_ "fmt"
	"github.com/huacnlee/revel"
	"log"
	"time"
)

var (
	Logger log.Logger
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.AssetsFilter,
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		InstramentFilter,
		revel.ActionInvoker, // Invoke the action.
	}

	Logger = log.Logger{}
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

var InstramentFilter = func(c *revel.Controller, fc []revel.Filter) {
	t1 := time.Now()
	fmt.Println("\nStarted", c.Request.Method, c.Request.URL.String(), "at", time.Now().Format(time.RFC3339), "from", c.Request.RemoteAddr)

	fc[0](c, fc[1:])

	duration := time.Since(t1)
	fmt.Println("\nComplated", c.Response.Status, "in", fmt.Sprintf("%.2fms", float64(duration.Nanoseconds()/1e4)/100.0))
}
