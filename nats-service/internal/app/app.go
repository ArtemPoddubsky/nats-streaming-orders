package app

import (
	"errors"
	"github.com/gorilla/mux"
	"html/template"
	"main/internal/config"
	"main/internal/inmemory"
	"main/internal/log"
	"main/internal/utils"
	"net/http"
	"time"
)

// App holds data needed for HTTP server.
type App struct {
	cfg         *config.Config
	memoryCache *inmemory.Cache
	router      *mux.Router
}

// NewApp returns new instance of App.
func NewApp(cfg *config.Config, memoryCache *inmemory.Cache) App {
	return App{
		cfg:         cfg,
		memoryCache: memoryCache,
		router:      mux.NewRouter(),
	}
}

// Run configures and runs HTTP server with GracefullShutdown.
func (a App) Run() {
	a.router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("materials/website/css"))))
	a.router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "materials/website/favicon.ico")
	})
	a.router.HandleFunc("/{id}", a.RenderPage).Methods(http.MethodGet)

	server := &http.Server{
		Addr:         a.cfg.Port,
		Handler:      a.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Logger.Infoln("Running server")
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Logger.Fatalln("ListenAndServe:", err)
		}
	}()

	utils.GracefullShutdown(server)
}

// RenderPage generates HTML page from cache by requested ID.
func (a App) RenderPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	a.memoryCache.Mutex.Lock()

	if _, ok := a.memoryCache.Storage[mux.Vars(r)["id"]]; !ok {
		a.memoryCache.Mutex.Unlock()
		utils.PageNotFound(w)
		return
	}

	t := template.Must(template.New("index.html").ParseFiles("materials/website/index.html"))

	if err := t.Execute(w, a.memoryCache.Storage[mux.Vars(r)["id"]]); err != nil {
		a.memoryCache.Mutex.Unlock()
		log.Logger.Errorln("template.Execute: ", err)
		utils.PageInternalError(w)
		return
	}
	a.memoryCache.Mutex.Unlock()
}
