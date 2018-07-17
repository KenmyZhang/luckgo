package route

import (
	"luckgo/api"
)

func (r *Router) InitStats() {
	r.Stats = r.Root.Group("/stats")
	r.Stats.GET("/master/db/connections/total", api.TotalMasterDbConnections)
	r.Stats.GET("/read/db/connections/total", api.TotalReadDbConnections)
	r.Stats.GET("/search/db/connections/total", api.TotalSearchDbConnections)
}
