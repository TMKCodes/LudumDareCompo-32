package main


import (
	"log";
	"math";
	"io/ioutil";
	"encoding/json";
	sf "bitbucket.org/krepa098/gosfml2";
)

type player struct {
	X int `json:"x"`;
	Y int `json:"y"`;
	Walk int;
	Health float32 `json:"health"`;
	Shape *sf.RectangleShape;
	Textures *textures;
}

func newPlayer(x int, y int, health float32, textures *textures) *player {
	player := new(player);
	player.X = x;
	player.Y = y;
	player.Health = health;
	player.Textures = textures;
	player.Walk = 1;
	var err error;
	player.Shape, err = sf.NewRectangleShape();
	if err != nil {
		log.Fatal(err);
	}
	for texture, _ := range textures.Texture {
		if textures.Texture[texture].Name == "stand-one.png" {
			player.Shape.SetTexture(textures.Texture[texture].Data, false);
			textureSize := textures.Texture[texture].Data.GetSize();
			player.Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
			player.Shape.SetOrigin(sf.Vector2f{float32(textureSize.X/2), float32(textureSize.Y/2)});
		}
	}
	player.Shape.Move(sf.Vector2f{float32(player.X), float32(player.Y)});
	return player;
}

func newPlayerFromFile(filename string) *player {
	file, err := ioutil.ReadFile(filename);
	if err != nil {
		log.Fatal(err);
	}
	player := new(player);
	err = json.Unmarshal(file, &player);
	if err != nil {
		log.Fatal(err);
	}
	return player;
}

func (this *player) WorldCollision(x float32, y float32, gamemap *gamemap) bool {
	for object, _ := range gamemap.Objects {
		objectRect := gamemap.Objects[object].Shape.GetGlobalBounds();
		tmpPlayer, err := sf.NewRectangleShape();
		if err != nil {
			log.Fatal(err);
		}
		tmpPlayer.SetSize(this.Shape.GetSize());
		tmpPlayer.Move(this.Shape.GetPosition());
		tmpPlayer.Move(sf.Vector2f{x,y});
		test, _ := objectRect.Intersects(tmpPlayer.GetGlobalBounds());
		if test == true {
			if gamemap.Objects[object].Sprite == "house.png" {
				if this.Health < 100 {
					this.Health += 1;
				}
			}
			return test;
		}
	}
	return false;
}

func (this *player) ZombieCollision(x float32, y float32, zombies []*zombie) bool {
	for zombie, _ := range zombies {
		zombieRect := zombies[zombie].Shape.GetGlobalBounds();
		tmpPlayer, err := sf.NewRectangleShape();
		if err != nil {
			log.Fatal(err);
		}
		tmpPlayer.SetSize(this.Shape.GetSize());
		tmpPlayer.Move(this.Shape.GetPosition());
		tmpPlayer.Move(sf.Vector2f{x,y});
		test, _ := zombieRect.Intersects(tmpPlayer.GetGlobalBounds());
		if test == true {
			return test;
		}
	}
	return false;
}

func (this *player) MeleeDamageZombies(x float32, y float32, zombies []*zombie) (int, []*zombie) {
	for zombie, _ := range zombies {
		zombieRect := zombies[zombie].Shape.GetGlobalBounds();
		tmpPlayer, err := sf.NewRectangleShape();
		if err != nil {
			log.Fatal(err);
		}
		tmpPlayer.SetSize(this.Shape.GetSize());
		tmpPlayer.Move(this.Shape.GetPosition());
		tmpPlayer.Move(sf.Vector2f{x,y});
		test, _ := zombieRect.Intersects(tmpPlayer.GetGlobalBounds());
		if test == true {
			if zombies[zombie].Health > 50 {
				zombies[zombie].Health -= 50;
			} else {
				zombies[zombie].Health = 0;
				zombies = append(zombies[:zombie], zombies[zombie+1:]...);
				return 1, zombies;
			}
		}
	}
	return 0, zombies;
}

func (this *player) Update(mouseCoords sf.Vector2f, gamemap *gamemap, zombies []*zombie) (int, []*zombie) {
	kills := 0;
	if this.Health > 0 {
		var speed float64;
		speed = 5.0;
		playerPos := this.Shape.GetPosition();
		angle := math.Atan2(float64(mouseCoords.Y - playerPos.Y), float64(mouseCoords.X - playerPos.X)) * (180/math.Pi) - 90;
		this.Shape.SetRotation(float32(angle));
		if sf.KeyboardIsKeyPressed(sf.KeyW) == true || sf.IsMouseButtonPressed(sf.MouseRight) {
			run := float64(playerPos.X - mouseCoords.X);
			rise := float64(playerPos.Y - mouseCoords.Y);
			length := math.Sqrt((rise * rise) + (run * run));
			x := -float32(run / length * speed);
			y := -float32(rise / length * speed);
			if this.WorldCollision(x, y, gamemap) == false && this.ZombieCollision(x, y, zombies) == false {
				this.Shape.Move(sf.Vector2f{x,y});
			}
			if this.Walk == 1 {
				this.Walk = 2;
				for texture, _ := range this.Textures.Texture {
					if this.Textures.Texture[texture].Name == "walk-one.png" {
						this.Shape.SetTexture(this.Textures.Texture[texture].Data, false);
						textureSize := this.Textures.Texture[texture].Data.GetSize();
						this.Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
						this.Shape.SetOrigin(sf.Vector2f{float32(textureSize.X/2), float32(textureSize.Y/2)});
						break;
					}
				}
			} else if this.Walk == 2 {
				this.Walk = 1;
				for texture, _ := range this.Textures.Texture {
					if this.Textures.Texture[texture].Name == "walk-two.png" {
						this.Shape.SetTexture(this.Textures.Texture[texture].Data, false);
						textureSize := this.Textures.Texture[texture].Data.GetSize();
						this.Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
						this.Shape.SetOrigin(sf.Vector2f{float32(textureSize.X/2), float32(textureSize.Y/2)});
						break;
					}
				}
			}
		} else if sf.KeyboardIsKeyPressed(sf.KeyS) == true {
			run := float64(playerPos.X - mouseCoords.X);
			rise := float64(playerPos.Y - mouseCoords.Y);
			length := math.Sqrt((rise * rise) + (run * run));
			x := float32(run / length * speed);
			y := float32(rise / length * speed);
			if this.WorldCollision(x, y, gamemap) == false && this.ZombieCollision(x, y, zombies) == false {
				this.Shape.Move(sf.Vector2f{x,y});
			}
			if this.Walk == 1 {
				this.Walk = 2;
				for texture, _ := range this.Textures.Texture {
					if this.Textures.Texture[texture].Name == "walk-one.png" {
						this.Shape.SetTexture(this.Textures.Texture[texture].Data, false);
						textureSize := this.Textures.Texture[texture].Data.GetSize();
						this.Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
						this.Shape.SetOrigin(sf.Vector2f{float32(textureSize.X/2), float32(textureSize.Y/2)});
						break;
					}
				}
			} else if this.Walk == 2 {
				this.Walk = 1;
				for texture, _ := range this.Textures.Texture {
					if this.Textures.Texture[texture].Name == "walk-two.png" {
						this.Shape.SetTexture(this.Textures.Texture[texture].Data, false);
						textureSize := this.Textures.Texture[texture].Data.GetSize();
						this.Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
						this.Shape.SetOrigin(sf.Vector2f{float32(textureSize.X/2), float32(textureSize.Y/2)});
						break;
					}
				}
			}
		} else {
			for texture, _ := range this.Textures.Texture {
				if this.Textures.Texture[texture].Name == "stand-one.png" {
					this.Shape.SetTexture(this.Textures.Texture[texture].Data, false);
					textureSize := this.Textures.Texture[texture].Data.GetSize();
					this.Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
					this.Shape.SetOrigin(sf.Vector2f{float32(textureSize.X/2), float32(textureSize.Y/2)});
					break;
				}
			}
		}
		if sf.IsMouseButtonPressed(sf.MouseLeft) == true {
			run := float64(playerPos.X - mouseCoords.X);
			rise := float64(playerPos.Y - mouseCoords.Y);
			length := math.Sqrt((rise * rise) + (run * run));
			x := -float32(run / length * speed);
			y := -float32(rise / length * speed);
			kills, zombies = this.MeleeDamageZombies(x,y, zombies);
			for texture, _ := range this.Textures.Texture {
				if this.Textures.Texture[texture].Name == "stand-two.png" {
					this.Shape.SetTexture(this.Textures.Texture[texture].Data, false);
					textureSize := this.Textures.Texture[texture].Data.GetSize();
					this.Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
					this.Shape.SetOrigin(sf.Vector2f{float32(textureSize.X/2), float32(textureSize.Y/2)});
					break;
				}
			}
		}
	}
	return kills, zombies;
}

func (this *player) Draw(renderWindow *sf.RenderWindow, renderStates sf.RenderStates) {
	if this.Health > 0 {
		renderWindow.Draw(this.Shape, renderStates);
	}
}

