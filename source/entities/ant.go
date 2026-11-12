package entities

import (
	"antlab/source/geometry_utils"
	"math"
	"math/rand/v2"
)

type AntJob int

const (
	Search AntJob = iota
	CarryFood
	Tired
)

type Ant struct {
	Position              geometry_utils.Point
	Speed                 float64
	Rotation              float64
	Job                   AntJob
	Capacity              int
	Stamina               float64
	FoodPheromoneCapacity float64
	homeWander            float64
	homeWanderTime        float64
}

func (ant *Ant) getRotation(rotation, dt, width, height float64, foodSmellAt func(x, y float64) float64,
	homeSmellAt func(x, y float64) (float64, float64),
	foodPheromoneAt func(x, y float64) (float64, float64)) float64 {
	x, y := ant.Position.X, ant.Position.Y

	if ant.Job == Search {
		if foodSmellAt(x, y) < 0.05 {
			targetX, targetY := foodPheromoneAt(x, y)
			if targetX >= 0 && targetY >= 0 && math.Hypot(targetX-x, targetY-y) >= 1e-6 {
				return ant.turnToFoodPheromone(rotation, dt, targetX, targetY)
			}

			const turnsPerSecond = 0.08 * 60
			if rand.Float64() < math.Min(1, turnsPerSecond*dt) {
				rotation += (rand.Float64() - 0.5) * math.Pi / 4
			}
			return rotation
		} else {
			dx, dy := (ant.Speed*dt)*1, (ant.Speed*dt)*1
			fx_, fy_ := (foodSmellAt(x+dx, y)-foodSmellAt(x-dx, y))/dx, (foodSmellAt(x, y+dy)-foodSmellAt(x, y-dy))/dy

			if math.Hypot(fx_, fy_) < 1e-6 {
				return rotation
			}

			target := math.Atan2(fy_, fx_)

			delta := math.Atan2(
				math.Sin(target-rotation),
				math.Cos(target-rotation),
			)

			maxTurn := math.Pi * dt
			delta = math.Max(-maxTurn, math.Min(maxTurn, delta))

			return rotation + delta
		}
	}

	if ant.Job == CarryFood {
		const maxWanderAngle = math.Pi / 12 // Up to 15 degrees either side of the home direction.
		const wanderTurnSpeed = math.Pi / 4
		ant.homeWanderTime -= dt
		if ant.homeWanderTime <= 0 {
			ant.homeWander = rand.Float64()*2 - 1
			ant.homeWanderTime = 0.2 + rand.Float64()*0.2
		}

		targetX, targetY := homeSmellAt(x, y)

		if targetX < 0 || targetY < 0 || math.Hypot(targetX-x, targetY-y) < 1e-6 {
			return rotation + ant.homeWander*wanderTurnSpeed*dt
		}

		target := math.Atan2(targetY-y, targetX-x) + ant.homeWander*maxWanderAngle
		delta := math.Atan2(math.Sin(target-rotation), math.Cos(target-rotation))
		maxTurn := math.Pi * dt
		delta = math.Max(-maxTurn, math.Min(maxTurn, delta))

		return rotation + delta
	}

	if ant.Job == Tired {
		if foodSmellAt(x, y) < 0.05 {
			targetX, targetY := foodPheromoneAt(x, y)
			if targetX >= 0 && targetY >= 0 && math.Hypot(targetX-x, targetY-y) >= 1e-6 {
				return ant.turnToFoodPheromone(rotation, dt, targetX, targetY)
			}
			targetX, targetY = homeSmellAt(x, y)

			if targetX < 0 || targetY < 0 || math.Hypot(targetX-x, targetY-y) < 1e-6 {
				return rotation
			}

			target := math.Atan2(targetY-y, targetX-x)
			delta := math.Atan2(math.Sin(target-rotation), math.Cos(target-rotation))
			maxTurn := math.Pi * dt
			delta = math.Max(-maxTurn, math.Min(maxTurn, delta))

			return rotation + delta
		} else {
			dx, dy := (ant.Speed*dt)*1, (ant.Speed*dt)*1
			fx_, fy_ := (foodSmellAt(x+dx, y)-foodSmellAt(x-dx, y))/dx, (foodSmellAt(x, y+dy)-foodSmellAt(x, y-dy))/dy

			if math.Hypot(fx_, fy_) < 1e-6 {
				return rotation
			}

			target := math.Atan2(fy_, fx_)

			delta := math.Atan2(
				math.Sin(target-rotation),
				math.Cos(target-rotation),
			)

			maxTurn := math.Pi * dt
			delta = math.Max(-maxTurn, math.Min(maxTurn, delta))

			return rotation + delta
		}
	}

	return rotation
}

func (ant *Ant) turnToFoodPheromone(rotation, dt, targetX, targetY float64) float64 {
	dx, dy := targetX-ant.Position.X, targetY-ant.Position.Y
	distance := math.Hypot(dx, dy)
	if distance < 1e-6 {
		return rotation
	}

	target := math.Atan2(dy, dx)
	delta := math.Atan2(math.Sin(target-rotation), math.Cos(target-rotation))
	// A fixed turn radius can orbit a nearby cell without ever reaching it.
	maxTurn := max(math.Pi, 2*ant.Speed/distance) * dt
	return rotation + math.Max(-maxTurn, math.Min(maxTurn, delta))
}

func (ant *Ant) Move(width, height, dt float64, foodSmellAt func(x, y float64) float64, getFood func(x, y float64, maxCnt int) int,
	homeSmellAt func(x, y float64) (float64, float64),
	anthillScan func(x, y float64) bool,
	foodPheromoneAt func(x, y float64) (float64, float64)) {
	if dt <= 0 {
		return
	}

	ant.Position.X += ant.Speed * math.Cos(ant.Rotation) * dt
	ant.Position.Y += ant.Speed * math.Sin(ant.Rotation) * dt

	xPadding, yPadding := width*0.01, height*0.01

	if ant.Position.X < xPadding {
		ant.Rotation = math.Pi - ant.Rotation
		ant.Position.X = xPadding
	}
	if ant.Position.Y < yPadding {
		ant.Rotation = -ant.Rotation
		ant.Position.Y = yPadding
	}

	if ant.Position.X > width-xPadding {
		ant.Rotation = math.Pi - ant.Rotation
		ant.Position.X = width - xPadding
	}
	if ant.Position.Y > height-yPadding {
		ant.Rotation = -ant.Rotation
		ant.Position.Y = height - yPadding
	}
}

func (ant *Ant) Update(width, height, dt float64, foodSmellAt func(x, y float64) float64, getFood func(x, y float64, maxCnt int) int,
	homeSmellAt func(x, y float64) (float64, float64),
	anthillScan func(x, y float64) bool,
	foodPheromoneAt func(x, y float64) (float64, float64),
	putFood func(int)) {

	if ant.Stamina < 0.07 && ant.Job == Search {
		ant.Job = Tired
	}

	if f := getFood(ant.Position.X, ant.Position.Y, ant.Capacity); (ant.Job == Search || ant.Job == Tired) &&
		foodSmellAt(ant.Position.X, ant.Position.Y) >= 0.8 && f > 0 {

		ant.Capacity -= f
		ant.Job = CarryFood
		ant.FoodPheromoneCapacity = 1.0
		ant.homeWanderTime = 0
	}

	if (ant.Job == CarryFood || ant.Job == Tired) && anthillScan(ant.Position.X, ant.Position.Y) {
		ant.Stamina = 1
		ant.Job = Search
		putFood(ant.Capacity)
		ant.Capacity = 5
	}

	ant.Rotation = ant.getRotation(ant.Rotation, dt, width, height, foodSmellAt, homeSmellAt, foodPheromoneAt)
}
