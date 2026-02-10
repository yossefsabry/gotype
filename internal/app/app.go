package app

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type App struct {
	screen   tcell.Screen
	model    *Model
	renderer *Renderer
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
	width, height := screen.Size()
	model.Layout.Recalculate(width, height)

	app := &App{
		screen:   screen,
		model:    model,
		renderer: NewRenderer(screen),
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
				a.model.Layout.Recalculate(width, height)
				needsRender = true
			case *tcell.EventKey:
				changed, shouldQuit := a.model.HandleKey(ev, time.Now())
				if shouldQuit {
					close(quit)
					return nil
				}
				if changed {
					needsRender = true
				}
			case *tcell.EventMouse:
				if ev.Buttons()&tcell.Button1 != 0 {
					x, y := ev.Position()
					if a.model.HandleClick(x, y, time.Now()) {
						needsRender = true
					}
				}
			}
		case <-ticker.C:
			if a.model.Update(time.Now()) {
				needsRender = true
			}
		}
	}
}
