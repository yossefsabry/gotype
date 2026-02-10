package app

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/yossefsabry/gotype/internal/storage"
)

type App struct {
	screen   tcell.Screen
	model    *Model
	renderer *Renderer
	store    *Persister
	data     storage.Data
	prefs    storage.Preferences
	finished bool
}

func Run() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := screen.Init(); err != nil {
		return err
	}
	screen.EnableMouse()
	defer screen.Fini()

	model := NewModel()
	path, data := loadPersistedData()
	if applyPreferences(model, data.Preferences) {
		model.Reset()
	}
	if data.BestScores == nil {
		data.BestScores = map[string]storage.BestScore{}
	}
	data.Preferences = preferencesFromModel(model)
	width, height := screen.Size()
	model.Layout.Recalculate(width, height, model.Options.Mode)
	var persister *Persister
	if path != "" {
		persister = NewPersister(path)
	}

	app := &App{
		screen:   screen,
		model:    model,
		renderer: NewRenderer(screen),
		store:    persister,
		data:     data,
		prefs:    data.Preferences,
	}
	if persister != nil {
		defer persister.Close()
	}
	return app.loop()
}

func (a *App) loop() error {
	events := make(chan tcell.Event, 32)
	quit := make(chan struct{})

	go func() {
		for {
			event := a.screen.PollEvent()
			if event == nil {
				close(events)
				return
			}
			select {
			case events <- event:
			case <-quit:
				return
			}
		}
	}()

	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	needsRender := true
	for {
		if needsRender {
			a.renderer.Render(a.model)
			needsRender = false
		}
		select {
		case event, ok := <-events:
			if !ok {
				return nil
			}
			switch ev := event.(type) {
			case *tcell.EventResize:
				a.screen.Sync()
				width, height := a.screen.Size()
				a.model.Layout.Recalculate(width, height, a.model.Options.Mode)
				if a.model.Timer.Finished {
					a.model.InitReviewStart()
				}
				needsRender = true
			case *tcell.EventKey:
				now := time.Now()
				changed, shouldQuit := a.model.HandleKey(ev, now)
				if shouldQuit {
					close(quit)
					return nil
				}
				if changed {
					needsRender = true
				}
				a.syncPersistence(now)
			case *tcell.EventMouse:
				if ev.Buttons()&tcell.Button1 != 0 {
					now := time.Now()
					x, y := ev.Position()
					if a.model.HandleClick(x, y, now) {
						needsRender = true
					}
					a.syncPersistence(now)
				}
			}
		case <-ticker.C:
			now := time.Now()
			if a.model.Update(now) {
				needsRender = true
			}
			a.syncPersistence(now)
		}
	}
}

func (a *App) syncPersistence(now time.Time) {
	if a.store != nil {
		current := preferencesFromModel(a.model)
		if current != a.prefs {
			a.prefs = current
			a.data.Preferences = current
			a.store.Save(a.data)
		}
	}
	if a.model.Timer.Finished && !a.finished {
		key := scoreKey(a.model.Options)
		previous, ok := a.data.BestScores[key]
		a.model.FinalizeResults(previous, ok)
		a.model.InitReviewStart()
		if updateBestScore(&a.data, a.model.Options, a.model.Stats, now) && a.store != nil {
			a.store.Save(a.data)
		}
	}
	a.finished = a.model.Timer.Finished
}
