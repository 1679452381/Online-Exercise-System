package main

import "online_exercise_system/router"

func main() {
	r := router.Router()
	r.Run(":8080")
}
