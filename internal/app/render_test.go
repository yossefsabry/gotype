package app

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func BenchmarkRender(b *testing.B) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		b.Fatalf("init screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(120, 40)

	model := NewModel()
	model.Layout.Recalculate(120, 40, model.Options.Mode, model.focusActive())
	renderer := NewRenderer(screen)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer.Render(model)
	}
}
