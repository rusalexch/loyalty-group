package user

func New(conf Config) *userModule {
	module := &userModule{
		pool: conf.Pool,
	}
	module.initRepo()

	return module
}
