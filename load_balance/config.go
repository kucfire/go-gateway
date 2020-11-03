package load_balance

type LoadBalanceConf interface {
	Attach(o Observer)
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
}

type Observer interface {
	Update()
}
