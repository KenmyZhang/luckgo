package route

import "luckgo/api"

func (r *Router) InitConfig() {
	r.Stats = r.Root.Group("/utils")
	r.Stats.GET("/config", api.GetConfig)
	r.Stats.GET("/version", api.GetVersionDetails)
}
