package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	orderLabel = "order"
)

var (
	issuedOrdersTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pvzservice_issued_orders_total",
		Help: "total number of orders issued",
	}, []string{
		orderLabel,
	})
)

func AddIssuedOrdersTotal(cnt int, order string) {
	issuedOrdersTotal.With(prometheus.Labels{
		orderLabel: order,
	}).Add(float64(cnt))
}
