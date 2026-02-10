package router

import (
	"net/http"
)

func MainRouter() *http.ServeMux {

	sRouter := studentsRouter()
	tRouter := teachersRouter()
	eRouter := execsRouter()

	sRouter.Handle("/", eRouter)
	tRouter.Handle("/", sRouter)

	return tRouter
}
