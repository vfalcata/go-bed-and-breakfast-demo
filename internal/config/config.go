// this config package can be accessed anywhere inside our application and we want to be careful to not importing anything other than what is absolutely mandatory
// we can cause an import cycle if we are not careful with our imports
// to do this because this will be a common use package we will only use the std library and not any package inside or app
package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application config
//
//	This allows for an app wide access to some configuration
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template // template for rendering
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool                // are we in prod or dev mode
	Session       *scs.SessionManager // session data
}
