package entities

import (
	"math"
	"testing"

	"antlab/source/geometry_utils"
)

func TestRotationSignalPriority(t *testing.T) {
	for _, scenario := range []struct {
		name                    string
		job                     AntJob
		food, pheromone         bool
		wantPheromone, wantHome bool
		direction               float64
	}{
		{"search_food", Search, true, true, false, false, 1},
		{"tired_food", Tired, true, true, false, false, 1},
		{"search_pheromone", Search, false, true, true, false, -1},
		{"tired_pheromone", Tired, false, true, true, false, -1},
		{"tired_home", Tired, false, false, true, true, 1},
		{"carrying_home", CarryFood, true, true, false, true, 1},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			ant := Ant{Position: geometry_utils.Point{X: 50, Y: 50}, Speed: 35, Job: scenario.job}
			var pheromoneCalled, homeCalled bool
			foodSmellAt := func(x, y float64) float64 {
				if scenario.food {
					return 0.2 + y*0.001 // Positive Y gradient.
				}
				return 0
			}
			foodPheromoneAt := func(x, y float64) (float64, float64) {
				pheromoneCalled = true
				if scenario.pheromone {
					return 50, 40 // Turn toward negative Y.
				}
				return -1, -1
			}
			homeSmellAt := func(x, y float64) (float64, float64) {
				homeCalled = true
				return 50, 60 // Turn toward positive Y.
			}

			const dt = 1.0 / 60
			got := ant.getRotation(0, dt, 100, 100, foodSmellAt, homeSmellAt, foodPheromoneAt)
			if math.IsNaN(got) || got*scenario.direction <= 0 || math.Abs(got) > math.Pi {
				t.Errorf("rotation = %g; want a finite turn with sign %g", got, scenario.direction)
			}
			if pheromoneCalled != scenario.wantPheromone || homeCalled != scenario.wantHome {
				t.Errorf("queries: pheromone=%t, home=%t; want %t, %t",
					pheromoneCalled, homeCalled, scenario.wantPheromone, scenario.wantHome)
			}
		})
	}
}
