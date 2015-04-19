package main


import (
	"log";
	"math";
	"io/ioutil";
	"encoding/json";
	sf "bitbucket.org/krepa098/gosfml2";
)

type zombie struct {
	X int `json:"x"`;
	Y int `json:"y"`;
	Health int `json:"health"`;
	Shape *sf.RectangleShape;
	Walk int;
	Textures *textures;
}

func newZombie(x int, y int, health int, textures *textures) *zombie {
	zombie := new(zombie);
	zombie.X = x;
	zombie.Y = y;
	zombie.Health = health;
	zombie.Textures = textures;
	zombie.Walk = 1;
	var err error;
	zombie.Shape, err = sf.NewRectangleShape();
	if err != nil {
		log.Fatal(err);
	}
	for texture, _ := range textures.Texture {
		if textures.Texture[texture].Name == "zombie-stand-two.png" {
			zombie.Shape.SetTexture(textures.Texture[texture].Data, false);
			textureSize := textures.Texture[texture].Data.GetSize();
			zombie.Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
			zombie.Shape.SetOrigin(sf.Vector2f{float32(textureSize.X/2), float32(textureSize.Y/2)});
		}
	}
	zombie.Shape.Move(sf.Vector2f{float32(zombie.X), float32(zombie.Y)});
	return zombie;
}

func newZombieFromFile(filename string) *zombie {
	file, err := ioutil.ReadFile(filename);
	if err != nil {
		log.Fatal(err);
	}
	zombie := new(zombie);
	err = json.Unmarshal(file, &zombie);
	if err != nil {
		log.Fatal(err);
	}
	return zombie;
}

func (this *zombie) WorldCollision(x float32, y float32, gamemap *gamemap) bool {
	for object, _ := range gamemap.Objects {
		objectRect := gamemap.Objects[object].Shape.GetGlobalBounds();
		tmpZombie, err := sf.NewRectangleShape();
		if err != nil {
			log.Fatal(err);
		}
		tmpZombie.SetSize(this.Shape.GetSize());
		tmpZombie.Move(this.Shape.GetPosition());
		tmpZombie.Move(sf.Vector2f{x,y});
		test, _ := objectRect.Intersects(tmpZombie.GetGlobalBounds());
		if test == true {
			return test;
		}
	}
	return false;
}

func (this *zombie) PlayerCollision(x float32, y float32, player *player) bool {
	playerRect := player.Shape.GetGlobalBounds();
	tmpZombie, err := sf.NewRectangleShape();
	if err != nil {
		log.Fatal(err);
	}
	tmpZombie.SetSize(this.Shape.GetSize());
	tmpZombie.Move(this.Shape.GetPosition());
	tmpZombie.Move(sf.Vector2f{x,y});
	test, _ := playerRect.Intersects(tmpZombie.GetGlobalBounds());
	return test;
}

func (this *zombie) BitePlayer(x float32, y float32, player *player) *player {
	playerRect := player.Shape.GetGlobalBounds();
	tmpZombie, err := sf.NewRectangleShape();
	if err != nil {
		log.Fatal(err);
	}
	tmpZombie.SetSize(this.Shape.GetSize());
	tmpZombie.Move(this.Shape.GetPosition());
	tmpZombie.Move(sf.Vector2f{x,y});
	test, _ := playerRect.Intersects(tmpZombie.GetGlobalBounds());
	if test == true {
		if player.Health > 0.5 {
			player.Health -= 0.5;
		} else {
			player.Health = 0;
			return player;
		}
	}
	return player;
}

func (this *zombie) HitHouse(x float32, y float32, gamemap *gamemap) *gamemap {
	for object, _ := range gamemap.Objects {
		if gamemap.Objects[object].Sprite == "house.png" {
			houseRect := gamemap.Objects[object].Shape.GetGlobalBounds();
			tmpZombie, err := sf.NewRectangleShape();
			if err != nil {
				log.Fatal(err);
			}
			tmpZombie.SetSize(this.Shape.GetSize());
			tmpZombie.Move(this.Shape.GetPosition());
			tmpZombie.Move(sf.Vector2f{x,y});
			test, _ := houseRect.Intersects(tmpZombie.GetGlobalBounds());
			if test == true {
				if gamemap.Objects[object].Health > 0.5 {
					gamemap.Objects[object].Health -= 0.5;
				} else {
					gamemap.Objects[object].Health = 0;
				}
				return gamemap;
			}
		}
	}
	return gamemap;
}

func (this *zombie) Update(playerCoords sf.Vector2f, gamemap *gamemap, player *player) (*gamemap, *player) {
	var speed float64;
	speed = 3.0;
	zombiePos := this.Shape.GetPosition();
	angle := math.Atan2(float64(playerCoords.Y - zombiePos.Y), float64(playerCoords.X - zombiePos.X)) * (180/math.Pi) - 90;
	this.Shape.SetRotation(float32(angle));
	run := float64(zombiePos.X - playerCoords.X);
	rise := float64(zombiePos.Y - playerCoords.Y);
	length := math.Sqrt((rise * rise) + (run * run));
	x := -float32(run / length * speed);
	y := -float32(rise / length * speed);
	if this.WorldCollision(x, y, gamemap) == false && this.PlayerCollision(x, y, player) == false{
		this.Shape.Move(sf.Vector2f{x,y});
	} else if this.PlayerCollision(x, y, player) == true {
		player = this.BitePlayer(x, y, player);
	} else if this.WorldCollision(x, y, gamemap) == true {
		gamemap = this.HitHouse(x, y, gamemap);
	}
	if this.Walk == 1 {
		this.Walk = 2;
		for texture, _ := range this.Textures.Texture {
			if this.Textures.Texture[texture].Name == "zombie-walk-one.png" {
				this.Shape.SetTexture(this.Textures.Texture[texture].Data, false);
				textureSize := this.Textures.Texture[texture].Data.GetSize();
				this.Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
				this.Shape.SetOrigin(sf.Vector2f{float32(textureSize.X/2), float32(textureSize.Y/2)});
			}
		}
	} else if this.Walk == 2 {
		this.Walk = 1;
		for texture, _ := range this.Textures.Texture {
			if this.Textures.Texture[texture].Name == "zombie-walk-two.png" {
				this.Shape.SetTexture(this.Textures.Texture[texture].Data, false);
				textureSize := this.Textures.Texture[texture].Data.GetSize();
				this.Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
				this.Shape.SetOrigin(sf.Vector2f{float32(textureSize.X/2), float32(textureSize.Y/2)});
			}
		}

	}
	return gamemap, player;
}

func (this *zombie) Draw(renderWindow *sf.RenderWindow, renderStates sf.RenderStates) {
	renderWindow.Draw(this.Shape, renderStates);
}

