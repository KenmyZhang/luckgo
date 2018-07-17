package route

import (
	"luckgo/tools"
)

func (r *Router) InitPrometheus() {
	r.Prometheus = r.Root.Group("")
	r.Prometheus.GET(tools.DefaultMetricPath, tools.LatestMetrics)
}
