package bootstrap

import "database/sql"

type Application struct {
	Env *Env
	DB  *sql.DB
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.DB = NewPostgresDatabase(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	ClosePostgresDBConnection(app.DB)
}
