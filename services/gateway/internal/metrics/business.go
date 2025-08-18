package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	UsersVerifiedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "gateway",
		Subsystem: "business",
		Name:      "users_verified_total",
		Help:      "Total number of users successfully verified (registration completed).",
	})

	PaymentsCompletedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "gateway",
		Subsystem: "business",
		Name:      "payments_completed_total",
		Help:      "Total number of successful payments completed.",
	})

	UsersDeletedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "gateway",
		Subsystem: "business",
		Name:      "users_deleted_total",
		Help:      "Total number of users who deleted their account.",
	})

	AdminUserBlockedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "gateway",
		Subsystem: "business",
		Name:      "admin_user_blocked_total",
		Help:      "Total number of admin 'block user' actions.",
	})

	AdminUserUnblockedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "gateway",
		Subsystem: "business",
		Name:      "admin_user_unblocked_total",
		Help:      "Total number of admin 'unblock user' actions.",
	})
)
