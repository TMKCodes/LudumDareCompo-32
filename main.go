package main

import (
	"fmt";
	"time";
	"runtime";
)

func main() {
	runtime.LockOSThread();
	fmt.Printf("Zombie hut is a Ludum Dare 32 game.\n\nThe theme for Ludum Dare 32 was An Unconventional Weapon.\n");
	game := NewGame("Zombie hut", 1280, 720, 32, true);
	for game.RenderWindow.IsOpen() {
		game.Update();
		game.Draw();
		time.Sleep(time.Duration(time.Millisecond * 50));
	}
}
