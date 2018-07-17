package route

import (
	"luckgo/api"
)

func (r *Router) HelloWorld() {
	r.Helloworld = r.Root.Group("/hello")
	r.Helloworld.POST("/world", api.Helloworld)

}
