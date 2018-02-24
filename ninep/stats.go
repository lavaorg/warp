package ninep

type StatsOps interface {
	statsRegister()
	statsUnregister()
}
