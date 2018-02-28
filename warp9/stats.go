package warp9

type StatsOps interface {
	statsRegister()
	statsUnregister()
}
