package main

import (
	"fmt";
	"log";
	"time";
	"math/rand";
	sf "bitbucket.org/krepa098/gosfml2";
)

type game struct {
	RenderWindow *sf.RenderWindow;
	GameView *sf.View;
	Map *gamemap;
	Textures *textures;
	Player *player;
	ZombieTimer time.Time;
	ZombieFreeze time.Time;
	ZombieNoMore bool;
	ZombieKills int;
	Zombies []*zombie;
	Status int;
	Welcome *sf.RectangleShape;
	WelcomeTimer time.Time;
	WelcomeTick int;
	Font *sf.Font;
	WelcomeText *sf.Text;
	Dialogue *sf.Text;
	PlayerHealth *sf.Text;
	HutHealth *sf.Text;
	EndingText *sf.Text;
	Music []*sf.Music;
	Sfxs []*sf.Music;
}

func NewGame(title string, width uint, height uint, bpp uint, vsync bool) *game {
	var err error;
	rand.Seed(time.Now().Unix());
	game := new(game);
	game.RenderWindow = sf.NewRenderWindow(sf.VideoMode{width, height, bpp}, title, sf.StyleDefault, sf.DefaultContextSettings());
	game.RenderWindow.SetVSyncEnabled(vsync);

	game.Textures = newTextures("res/textures/");
	game.Map = newMap("res/map.json", game.Textures);
	game.Player = newPlayer(int(width/2), int(height/2), 100, game.Textures);

	game.ZombieTimer = time.Now();
	game.ZombieKills = 0;
	game.ZombieNoMore = false;

	game.Status = 0;

	game.Font, err = sf.NewFontFromFile("res/fonts/UbuntuMono-R.ttf");
	if err != nil {
		log.Fatal(err);
	}

	game.WelcomeTimer = time.Now();
	game.WelcomeTick = 0;
	game.Welcome, err = sf.NewRectangleShape();
	if err != nil {
		log.Fatal(err);
	}
	game.Welcome.SetSize(sf.Vector2f{float32(width), float32(height)});
	for texture, _ := range game.Textures.Texture {
		if game.Textures.Texture[texture].Name == "mouse-hit.png" {
			game.Welcome.SetTexture(game.Textures.Texture[texture].Data, false);
			break;
		}
	}

	game.WelcomeText, err = sf.NewText(game.Font);
	if err != nil {
		log.Fatal(err);
	}
	game.WelcomeText.SetString("Press Hit or Walk to star the game.\nPress Enter to reset the game.");
	game.WelcomeText.SetCharacterSize(24);
	game.WelcomeText.SetColor(sf.Color{255,255,255,255});
	game.WelcomeText.SetPosition(sf.Vector2f{350,450});

	game.EndingText, err = sf.NewText(game.Font);
	if err != nil {
		log.Fatal(err);
	}
	game.EndingText.SetString("");
	game.EndingText.SetCharacterSize(24);
	game.EndingText.SetColor(sf.Color{255,255,255,255});
	game.EndingText.SetPosition(sf.Vector2f{340,450});

	game.Dialogue, err = sf.NewText(game.Font);
	if err != nil {
		log.Fatal(err);
	}
	game.PlayerHealth, err = sf.NewText(game.Font);
	if err != nil {
		log.Fatal(err);
	}
	game.PlayerHealth.SetCharacterSize(24);
	game.PlayerHealth.SetColor(sf.Color{255,255,255,255});
	game.PlayerHealth.SetPosition(sf.Vector2f{50,650});


	game.HutHealth, err = sf.NewText(game.Font);
	if err != nil {
		log.Fatal(err);
	}
	game.HutHealth.SetCharacterSize(24);
	game.HutHealth.SetColor(sf.Color{255,255,255,255});
	game.HutHealth.SetPosition(sf.Vector2f{1000,650});

	game.Music = make([]*sf.Music, 0);
	music, _ := sf.NewMusicFromFile("res/music/Unconventional-Weapon-Bang-Bang.ogg");
	game.Music = append(game.Music, music);
	game.Music[0].Play();
	game.Music[0].SetLoop(true);

	game.Sfxs = make([]*sf.Music, 0);
	sfx, _ := sf.NewMusicFromFile("res/sfx/hit-1.wav");
	game.Sfxs = append(game.Sfxs, sfx);
	sfx, _ = sf.NewMusicFromFile("res/sfx/hit-2.wav");
	game.Sfxs = append(game.Sfxs, sfx);
	sfx, _ = sf.NewMusicFromFile("res/sfx/hit-3.wav");
	game.Sfxs = append(game.Sfxs, sfx);



	return game;
}

func (this *game) Update() {
	for event := this.RenderWindow.PollEvent(); event != nil; event = this.RenderWindow.PollEvent() {
		switch ev := event.(type) {
			case sf.EventKeyReleased:
				if ev.Code == sf.KeyEscape {
					this.RenderWindow.Close();
				}
				if ev.Code == sf.KeyReturn {
					this.Status = 0;
					this.ZombieKills = 0;
					this.ZombieTimer = time.Now();
					this.ZombieNoMore = false;
					this.Zombies = make([]*zombie, 0);
					this.Map = newMap("res/map.json", this.Textures);
					renderWindowSize := this.RenderWindow.GetSize();
					this.Player = newPlayer(int(renderWindowSize.X/2), int(renderWindowSize.Y/2), 100, this.Textures);
				}
			case sf.EventClosed:
				this.RenderWindow.Close();
		}
	}
	if this.Status == 0 {
		if sf.IsMouseButtonPressed(sf.MouseLeft) == true || sf.IsMouseButtonPressed(sf.MouseRight) == true {
			this.Status = 1;
		} else {
			now := time.Now();
			if this.WelcomeTick == 0 {
				if now.Sub(this.WelcomeTimer) > time.Duration(time.Second * 2) {
					for texture, _ := range this.Textures.Texture {
						if this.Textures.Texture[texture].Name == "mouse-walk.png" {
							this.Welcome.SetTexture(this.Textures.Texture[texture].Data, false);
							this.WelcomeTick = 1;
							this.WelcomeTimer = time.Now();
							break;
						}
					}
				}
			} else {
				if now.Sub(this.WelcomeTimer) > time.Duration(time.Second * 2) {
					for texture, _ := range this.Textures.Texture {
						if this.Textures.Texture[texture].Name == "mouse-hit.png" {
							this.Welcome.SetTexture(this.Textures.Texture[texture].Data, false);
							this.WelcomeTick = 0;
							this.WelcomeTimer = time.Now();
							break;
						}
					}
				}
			}
		}
	} else if this.Status == 1 {
		if this.ZombieKills == 0 {
			this.Dialogue.SetString("It seems there is a wave of zombies coming.\nI need to kill them with my bare hands?");
		} else if this.ZombieKills == 1 {
			this.Dialogue.SetString("Do you know how it feels to be the last person alive?");
		} else if this.ZombieKills == 2 {
			this.Dialogue.SetString("Well it is my own fault for creating these zombies.");
		} else if this.ZombieKills == 3 {
			this.Dialogue.SetString("");
		} else if this.ZombieKills >= 10 && this.ZombieKills < 12 {
			this.Dialogue.SetString("These Zombies are very annoying weapons!");
		} else if this.ZombieKills >= 12 && this.ZombieKills < 14 {
			this.Dialogue.SetString("Seriously there is more of them coming?");
		} else if this.ZombieKills >= 14 && this.ZombieKills < 16 {
			this.Dialogue.SetString("On another note. I need to watch my health.");
		} else if this.ZombieKills == 16 {
			this.Dialogue.SetString("");
		} else {
			if this.Player.Health < 30 {
				this.Dialogue.SetString("Better go to the hut now to heal myself!");
			} else {
				this.Dialogue.SetString("");
			}
		}
		this.Dialogue.SetCharacterSize(24);
		this.Dialogue.SetColor(sf.Color{255,255,255,255});
		playerPos := this.Player.Shape.GetPosition();
		this.Dialogue.SetPosition(sf.Vector2f{playerPos.X - 300, playerPos.Y - 100});
		kills, zombies := this.Player.Update(this.RenderWindow.MapPixelToCoords(sf.MouseGetPosition(this.RenderWindow), this.RenderWindow.GetDefaultView()), this.Map, this.Zombies);
		if kills > 0 {
			this.ZombieKills += kills;
			this.Zombies = zombies;
			this.ZombieFreeze = time.Now();
			fmt.Printf("Zombie killed by the player.\n");
			fmt.Printf("Zombies: %v\n", len(this.Zombies));
			fmt.Printf("Zombies killed: %v\n", this.ZombieKills);
			fmt.Printf("Player health: %v\n", this.Player.Health);
			sfx := rand.Intn(2);
			this.Sfxs[sfx].Play();
		}
		if this.ZombieKills < 10 {
			now := time.Now();
			if len(this.Zombies) < 5 {
				if now.Sub(this.ZombieTimer) > time.Duration(time.Second * 5) {
					this.ZombieTimer = time.Now();
					w := rand.Intn(3);
					var zombie *zombie;
					if w == 0 {
						x := rand.Intn(1150/50)+1;
						zombie = newZombie(x*50, 50, 150, this.Textures);
					} else if w == 1 {
						y := rand.Intn(600/50)+1;
						zombie = newZombie(50, y*50, 150, this.Textures);
					} else if w == 2 {
						x := rand.Intn(1150/50)+1;
						zombie = newZombie(x*50, 650, 150, this.Textures);
					} else if w == 3 {
						y := rand.Intn(600/50)+1;
						zombie = newZombie(1215, y*50, 150, this.Textures);
					}
					this.Zombies = append(this.Zombies, zombie);
					fmt.Printf("New Zombie added to game.\n");
					fmt.Printf("Zombies: %v\n", len(this.Zombies));
				}
			}
		} else if this.ZombieKills >= 10 && this.ZombieKills < 20 {
			now := time.Now();
			if len(this.Zombies) < 10 {
				if now.Sub(this.ZombieTimer) > time.Duration(time.Second * 4) {
					this.ZombieTimer = time.Now();
					w := rand.Intn(3);
					var zombie *zombie;
					if w == 0 {
						x := rand.Intn(1150/50)+1;
						zombie = newZombie(x*50, 50, 150, this.Textures);
					} else if w == 1 {
						y := rand.Intn(600/50)+1;
						zombie = newZombie(50, y*50, 150, this.Textures);
					} else if w == 2 {
						x := rand.Intn(1150/50)+1;
						zombie = newZombie(x*50, 650, 150, this.Textures);
					} else if w == 3 {
						y := rand.Intn(600/50)+1;
						zombie = newZombie(1215, y*50, 150, this.Textures);
					}
					this.Zombies = append(this.Zombies, zombie);
					fmt.Printf("New Zombie added to game.\n");
					fmt.Printf("Zombies: %v\n", len(this.Zombies));
				}
			}
		} else if this.ZombieKills >= 20 && this.ZombieKills < 35 {
			now := time.Now();
			if len(this.Zombies) < 15 {
				if now.Sub(this.ZombieTimer) > time.Duration(time.Second * 3) {
					this.ZombieTimer = time.Now();
					w := rand.Intn(3);
					var zombie *zombie;
					if w == 0 {
						x := rand.Intn(1150/50)+1;
						zombie = newZombie(x*50, 50, 150, this.Textures);
					} else if w == 1 {
						y := rand.Intn(600/50)+1;
						zombie = newZombie(50, y*50, 150, this.Textures);
					} else if w == 2 {
						x := rand.Intn(1150/50)+1;
						zombie = newZombie(x*50, 650, 150, this.Textures);
					} else if w == 3 {
						y := rand.Intn(600/50)+1;
						zombie = newZombie(1215, y*50, 150, this.Textures);
					}
					this.Zombies = append(this.Zombies, zombie);
					fmt.Printf("New Zombie added to game.\n");
					fmt.Printf("Zombies: %v\n", len(this.Zombies));
				}
			}
		} else if this.ZombieKills >= 35 && this.ZombieKills < 55 {
			now := time.Now();
			if len(this.Zombies) < 20 {
				if now.Sub(this.ZombieTimer) > time.Duration(time.Second * 2) {
					this.ZombieTimer = time.Now();
					w := rand.Intn(3);
					var zombie *zombie;
					if w == 0 {
						x := rand.Intn(1150/50)+1;
						zombie = newZombie(x*50, 50, 150, this.Textures);
					} else if w == 1 {
						y := rand.Intn(600/50)+1;
						zombie = newZombie(50, y*50, 150, this.Textures);
					} else if w == 2 {
						x := rand.Intn(1150/50)+1;
						zombie = newZombie(x*50, 650, 150, this.Textures);
					} else if w == 3 {
						y := rand.Intn(600/50)+1;
						zombie = newZombie(1215, y*50, 150, this.Textures);
					}
					this.Zombies = append(this.Zombies, zombie);
					fmt.Printf("New Zombie added to game.\n");
					fmt.Printf("Zombies: %v\n", len(this.Zombies));
				}
			}
		} else if this.ZombieKills >= 55 && this.ZombieKills < 80 {
			now := time.Now();
			if len(this.Zombies) < 25 {
				if now.Sub(this.ZombieTimer) > time.Duration(time.Second * 1) {
					this.ZombieTimer = time.Now();
					w := rand.Intn(3);
					var zombie *zombie;
					if w == 0 {
						x := rand.Intn(1150/50)+1;
						zombie = newZombie(x*50, 50, 150, this.Textures);
					} else if w == 1 {
						y := rand.Intn(600/50)+1;
						zombie = newZombie(50, y*50, 150, this.Textures);
					} else if w == 2 {
						x := rand.Intn(1150/50)+1;
						zombie = newZombie(x*50, 650, 150, this.Textures);
					} else if w == 3 {
						y := rand.Intn(600/50)+1;
						zombie = newZombie(1215, y*50, 150, this.Textures);
					}
					this.Zombies = append(this.Zombies, zombie);
					fmt.Printf("New Zombie added to game.\n");
					fmt.Printf("Zombies: %v\n", len(this.Zombies));
				}
			}
		} else if this.ZombieKills >= 80 && this.ZombieKills < 110 {
			now := time.Now();
			if len(this.Zombies) < 30 {
				if now.Sub(this.ZombieTimer) > time.Duration(time.Millisecond * 500) {
					this.ZombieTimer = time.Now();
					w := rand.Intn(3);
					var zombie *zombie;
					if w == 0 {
						x := rand.Intn(1150/50)+1;
						zombie = newZombie(x*50, 50, 150, this.Textures);
					} else if w == 1 {
						y := rand.Intn(600/50)+1;
						zombie = newZombie(50, y*50, 150, this.Textures);
					} else if w == 2 {
						x := rand.Intn(1150/50)+1;
						zombie = newZombie(x*50, 650, 150, this.Textures);
					} else if w == 3 {
						y := rand.Intn(600/50)+1;
						zombie = newZombie(1215, y*50, 150, this.Textures);
					}
					this.Zombies = append(this.Zombies, zombie);
					fmt.Printf("New Zombie added to game.\n");
					fmt.Printf("Zombies: %v\n", len(this.Zombies));
				}
			}
		} else if this.ZombieKills >= 110 && this.ZombieKills < 160 {
			now := time.Now();
			if len(this.Zombies) < 50 {
				if now.Sub(this.ZombieTimer) > time.Duration(time.Millisecond * 250) {
					this.ZombieTimer = time.Now();
					w := rand.Intn(3);
					var zombie *zombie;
					if w == 0 {
						x := rand.Intn(1150/50)+1;
						zombie = newZombie(x*50, 50, 150, this.Textures);
					} else if w == 1 {
						y := rand.Intn(600/50)+1;
						zombie = newZombie(50, y*50, 150, this.Textures);
					} else if w == 2 {
						x := rand.Intn(1150/50)+1;
						zombie = newZombie(x*50, 650, 150, this.Textures);
					} else if w == 3 {
						y := rand.Intn(600/50)+1;
						zombie = newZombie(1215, y*50, 150, this.Textures);
					}
					this.Zombies = append(this.Zombies, zombie);
					fmt.Printf("New Zombie added to game.\n");
					fmt.Printf("Zombies: %v\n", len(this.Zombies));
				}
			}
		} else {
			this.ZombieNoMore = true;
		}

		if this.ZombieNoMore == true && len(this.Zombies) == 0 {
			this.Status = 3;
			this.EndingText.SetString(fmt.Sprintf("You won the zombies!\nYou killed %v zombies.\nPress Enter to replay.", this.ZombieKills));
		}

		now := time.Now();
		if now.Sub(this.ZombieFreeze) > time.Duration(time.Millisecond * 500) {
			for zombie, _ := range this.Zombies {
				if this.Player.Health > 0 {
					this.Map, this.Player = this.Zombies[zombie].Update(this.Player.Shape.GetPosition(), this.Map, this.Player);
				} else {
					for object, _ := range this.Map.Objects {
						if this.Map.Objects[object].Sprite == "house.png" {
							this.Map, this.Player = this.Zombies[zombie].Update(this.Map.Objects[object].Shape.GetPosition(), this.Map, this.Player);
							break;
						}
					}
				}
			}
		}
		this.PlayerHealth.SetString(fmt.Sprintf("Player: %.1f %% health", this.Player.Health));
		for object, _ := range this.Map.Objects {
			if this.Map.Objects[object].Sprite == "house.png" {
				this.HutHealth.SetString(fmt.Sprintf("Hut: %.1f %% health", this.Map.Objects[object].Health));
				break;
			}
		}
		if this.Player.Health == 0 {
			this.Status = 3;
			this.EndingText.SetString(fmt.Sprintf("You lost to the zombies!\nYou killed %v zombies.\nPress Enter to replay.", this.ZombieKills));
		}
		//this.Map.Update();
	}
}

func (this *game) Draw() {
	this.RenderWindow.Clear(sf.Color{255,255,255,255});
	renderStates := sf.DefaultRenderStates();
	this.Map.DrawGround(this.RenderWindow, renderStates);
	this.Player.Draw(this.RenderWindow, renderStates);
	for zombie, _ := range this.Zombies {
		this.Zombies[zombie].Draw(this.RenderWindow, renderStates);
	}
	this.Map.DrawObjects(this.RenderWindow, renderStates);
	if this.Status == 0 {
		this.Welcome.Draw(this.RenderWindow, renderStates);
		this.WelcomeText.Draw(this.RenderWindow, renderStates);
	} else if this.Status == 1 {
		this.Dialogue.Draw(this.RenderWindow, renderStates);
		this.PlayerHealth.Draw(this.RenderWindow, renderStates);
		this.HutHealth.Draw(this.RenderWindow, renderStates);
	} else if this.Status == 3 {
		this.EndingText.Draw(this.RenderWindow, renderStates);
	}
	this.RenderWindow.Display();
}
