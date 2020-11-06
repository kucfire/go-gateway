package config

type LoadBalance interface {
	Add(...string) error
	Get(key string) (string, error)
	Next() string
}

type LoadBalanceConf interface {
}
