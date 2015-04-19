package main

import (
	"log";
	"io/ioutil";
	"encoding/json";
	sf "bitbucket.org/krepa098/gosfml2";
)

type gamemap struct {
	X int `json:"x"`;
	Y int `json:"y"`;
	Width int `json:"width"`;
	Height int `json:"height"`;
	Sprite string `json:"sprite"`;
	Color color `json:"color"`;
	Objects []object `json:"objects"`;
	Shape *sf.RectangleShape;
}

type color struct {
	R uint8 `json:"r"`;
	G uint8 `json:"G"`;
	B uint8 `json:"B"`;
	A uint8 `json:"A"`;
}

type object struct {
	X int `json:"x"`;
	Y int `json:"y"`;
	Health float32 `json:"health"`;
	Sprite string `json:"sprite"`;
	Collide bool `json:"collide"`;
	Shape *sf.RectangleShape;
}

func newMap(filename string, textures *textures) *gamemap {
	file, err := ioutil.ReadFile(filename);
	if err != nil {
		log.Fatal(err);
	}
	gamemap := new(gamemap);
	err = json.Unmarshal(file, &gamemap);
	if err != nil {
		log.Fatal(err);
	}
	gamemap.Shape, err = sf.NewRectangleShape();
	if err != nil {
		log.Fatal(err);
	}

	gamemap.Shape.SetFillColor(sf.Color{gamemap.Color.R, gamemap.Color.G, gamemap.Color.B, gamemap.Color.A});
	gamemap.Shape.SetSize(sf.Vector2f{float32(gamemap.Width), float32(gamemap.Height)});
	gamemap.Shape.Move(sf.Vector2f{float32(gamemap.X), float32(gamemap.Y)});

	for texture := range textures.Texture {
		if gamemap.Sprite == textures.Texture[texture].Name {
			gamemap.Shape.SetTexture(textures.Texture[texture].Data, false);
		}
	}

	for object, _ := range gamemap.Objects {
		if gamemap.Objects[object].Sprite != "" {
			for texture, _ := range textures.Texture {
				if textures.Texture[texture].Name != "" && textures.Texture[texture].Name == gamemap.Objects[object].Sprite {
					gamemap.Objects[object].Shape, err = sf.NewRectangleShape();
					if err != nil {
						log.Fatal(err);
					}
					gamemap.Objects[object].Shape.SetTexture(textures.Texture[texture].Data, false);
					textureSize := textures.Texture[texture].Data.GetSize();
					gamemap.Objects[object].Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
					gamemap.Objects[object].Shape.Move(sf.Vector2f{float32(gamemap.X + gamemap.Objects[object].X), float32(gamemap.Y + gamemap.Objects[object].Y)});
					break;
				}
			}
		}
	}
	x := 0.0;
	for i := 10.0; i < 1260; i+=x {
		// top
		ob := new(object);
		ob.Shape, err = sf.NewRectangleShape();
		ob.Health = 100;
		if err != nil {
			log.Fatal(err);
		}
		for texture, _ := range textures.Texture {
			if textures.Texture[texture].Name != "" && textures.Texture[texture].Name == "tree.png" {
				ob.Shape.SetTexture(textures.Texture[texture].Data, false);
				textureSize := textures.Texture[texture].Data.GetSize();
				ob.Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
				break;
			}
		}
		size := ob.Shape.GetSize();
		x = float64(size.X + 5);
		ob.Shape.Move(sf.Vector2f{float32(float64(gamemap.X) + i), float32(gamemap.Y)});
		gamemap.Objects = append(gamemap.Objects, (*ob));
		// bottom
		ob = new(object);
		ob.Shape, err = sf.NewRectangleShape();
		ob.Health = 100;
		if err != nil {
			log.Fatal(err);
		}
		for texture, _ := range textures.Texture {
			if textures.Texture[texture].Name != "" && textures.Texture[texture].Name == "tree.png" {
				ob.Shape.SetTexture(textures.Texture[texture].Data, false);
				textureSize := textures.Texture[texture].Data.GetSize();
				ob.Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
				break;
			}
		}
		size = ob.Shape.GetSize();
		ob.Shape.Move(sf.Vector2f{float32(float64(gamemap.X) + i), float32(float32(gamemap.Height)-size.Y)});
		gamemap.Objects = append(gamemap.Objects, (*ob));
	}
	y := 0.0;
	for i := 50.0; i < 670; i+=y {
		// left
		ob := new(object);
		ob.Shape, err = sf.NewRectangleShape();
		ob.Health = 100;
		if err != nil {
			log.Fatal(err);
		}
		for texture, _ := range textures.Texture {
			if textures.Texture[texture].Name != "" && textures.Texture[texture].Name == "tree.png" {
				ob.Shape.SetTexture(textures.Texture[texture].Data, false);
				textureSize := textures.Texture[texture].Data.GetSize();
				ob.Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
				break;
			}
		}
		size := ob.Shape.GetSize();
		y = float64(size.Y + 5);
		ob.Shape.Move(sf.Vector2f{float32(gamemap.X), float32(float64(gamemap.Y) + i)});
		gamemap.Objects = append(gamemap.Objects, (*ob));
		// right
		ob = new(object);
		ob.Shape, err = sf.NewRectangleShape();
		ob.Health = 100;
		if err != nil {
			log.Fatal(err);
		}
		for texture, _ := range textures.Texture {
			if textures.Texture[texture].Name != "" && textures.Texture[texture].Name == "tree.png" {
				ob.Shape.SetTexture(textures.Texture[texture].Data, false);
				textureSize := textures.Texture[texture].Data.GetSize();
				ob.Shape.SetSize(sf.Vector2f{float32(textureSize.X), float32(textureSize.Y)});
				break;
			}
		}
		size = ob.Shape.GetSize();
		ob.Shape.Move(sf.Vector2f{float32(float32(gamemap.Width) - size.X), float32(float64(gamemap.Y) + i)});
		gamemap.Objects = append(gamemap.Objects, (*ob));
	}
	return gamemap;
}

func (this *gamemap) Update() {
	for object, _ := range this.Objects {
		this.Objects[object].Shape.Move(sf.Vector2f{1.0, 1.0});
	}
}

func (this *gamemap) DrawGround(renderWindow *sf.RenderWindow, renderStates sf.RenderStates) {
	renderWindow.Draw(this.Shape, renderStates);
}

func (this *gamemap) DrawObjects(renderWindow *sf.RenderWindow, renderStates sf.RenderStates) {
	for object, _ := range this.Objects {
		if this.Objects[object].Health > 0 {
			renderWindow.Draw(this.Objects[object].Shape, renderStates);
		}
	}
}


