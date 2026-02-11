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

// first initialization of the application
func Run() error {
	// creating new window for application
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	// initialize the screen
	if err := screen.Init(); err != nil {
		return err
	}
	screen.EnableMouse()
	// ensure the screen is finalized when the application exits 
	// (freeing resources, restoring terminal state, etc.)
	defer screen.Fini()

	// from here get the data and initialize the model, renderer, and app
	model := NewModel()
	// load default preferences and best scores, and apply them to the model
	path, data := loadPersistedData()
	if applyPreferences(model, data.Preferences) {
		model.Reset()
	}
	if data.BestScores == nil {
		data.BestScores = map[string]storage.BestScore{}
	}
	data.Preferences = preferencesFromModel(model)
	// auto calculate resize the layout based on the current screen size 
	// and model options
	width, height := screen.Size()
	model.Layout.Recalculate(width, height, model.Options.Mode, model.focusActive())

	// persister is responsible for saving the preferences and best scores to disk,
	// it runs in a separate goroutine and listens for changes to the data and 
	// saves it asynchronously to avoid blocking the main UI thread
	var persister *Persister
	if path != "" {
		persister = NewPersister(path)
	}

	// there is were too start the main loop
	app := &App{
		screen:   screen,
		model:    model,
		renderer: NewRenderer(screen),
		store:    persister,
		data:     data,
		prefs:    data.Preferences,
	}
	// ensure the persister is closed when the application exits to clean up resources
	if persister != nil {
		defer persister.Close()
	}
	return app.loop()
}

// loop main event that is run each 80 milliseconds to check for user input
// and update the UI accordingly
func (a *App) loop() error {
	// for handle any interaction with the window
	events := make(chan tcell.Event, 32)
	// for handle any singnals for close the application
	quit := make(chan struct{})

	// start goroutine to listen for events
	go func() {
		for {
			event := a.screen.PollEvent()
			if event == nil {
				close(events)
				return
			}
			select {
				// store the event in the events channel to be processed by the main loop
			case events <- event:
				// if event is quit signal, close the events channel and exit the goroutine
			case <-quit:
				return
			}
		}
	}()

	// update UI each 80ms
	ticker := time.NewTicker(80 * time.Millisecond)
	// ensure stop when exit the function to clean up
	defer ticker.Stop()

	// main loop to handle events and update the UI
	needsRender := true
	for {
		if needsRender {
			a.renderer.Render(a.model)
			needsRender = false
		}
		select {
			// reading the events from the events channel and handle them
		case event, ok := <-events:
			// no events and the channel is closed, exit the loop
			if !ok {
				return nil
			}
			switch ev := event.(type) {
			case *tcell.EventResize:
				a.screen.Sync()
				width, height := a.screen.Size()
				a.model.Layout.Recalculate(width, height, a.model.Options.Mode,
				a.model.focusActive())
				if a.model.Timer.Finished {
					a.model.InitReviewStart()
				}
				needsRender = true
			case *tcell.EventKey:
				now := time.Now()
				// handle the key event and update the model accordingly,
				// if the model
				changed, shouldQuit := a.model.HandleKey(ev, now)
				if shouldQuit {
					close(quit)
					return nil
				}
				if changed {
					needsRender = true
				}
				// handle new data and save it to disk if needed
				a.syncPersistence(now)
			case *tcell.EventMouse:
				if ev.Buttons()&tcell.Button1 != 0 {
					now := time.Now()
					x, y := ev.Position()
					// check for the clicks on the interactive regions and 
					// apply the corresponding actions,
					if a.model.HandleClick(x, y, now) {
						needsRender = true
					}
					a.syncPersistence(now)
				}
			}

		// update the timer for application 
		case <-ticker.C:
			now := time.Now()
			if a.model.Update(now) {
				needsRender = true
			}
			a.syncPersistence(now)
		}
	}
}

// saving new preferences and best scores to disk when they change,
// and calculating the final results
func (a *App) syncPersistence(now time.Time) {
	// if there is a store, check if the preferences have changed and save them
	// make sure to only save when there is a change to avoid 
	// unnecessary writes to disk
	if a.store != nil {
		current := preferencesFromModel(a.model)
		if current != a.prefs {
			a.prefs = current
			a.data.Preferences = current
			a.store.Save(a.data)
		}
	}
	// when the timer finishes for the first time,
	// calculate the final results and update the best score if needed
	if a.model.Timer.Finished && !a.finished {
		key := scoreKey(a.model.Options)
		previous, ok := a.data.BestScores[key]
		a.model.FinalizeResults(previous, ok)
		a.model.InitReviewStart()
		if updateBestScore(&a.data, a.model.Options, a.model.Stats, now) && a.store != nil {
			a.store.Save(a.data)
		}
	}
	// set the finished flag to true to avoid recalculating the results 
	// and updating the best score
	a.finished = a.model.Timer.Finished
}
