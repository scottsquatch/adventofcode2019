package main

import (
	"github.com/scottsquatch/adventofcode2019/common"
)

// HullPaintingRobot represents a hull painting robot
type HullPaintingRobot struct {
	software common.IntCodeProgram
	pc       *common.Computer
	camera   chan int64
	arm      chan int64
}

// NewHullPaintingRobot Create a new hullpaintingrobot instance
func NewHullPaintingRobot(software common.IntCodeProgram) *HullPaintingRobot {
	camera := make(chan int64, 1)
	arm := make(chan int64)
	return &HullPaintingRobot{software, common.NewComputer(camera, arm), camera, arm}
}

// PaintHandler is executed for every paint action
type PaintHandler func(int64)

// MovementHandler is executed for each movement action
type MovementHandler func(int64)

// CameraHandler generates a number for the camera
type CameraHandler func() int64

// Start the robot
func (robot *HullPaintingRobot) Start(paintAction PaintHandler, movementAction MovementHandler, cameraAction CameraHandler) {

	halt := make(chan bool)
	go func() {
		robot.pc.Run(&robot.software)
		halt <- true
	}()

	isMove := false
	for {

		robot.camera <- cameraAction()
		for i := 0; i < 2; i++ {
			select {
			case <-halt:
				return
			case num := <-robot.arm:
				if isMove {
					movementAction(num)
				} else {
					paintAction(num)
				}
				isMove = !isMove
			}
		}
	}
}
