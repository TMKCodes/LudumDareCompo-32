package main

import (
	"log";
	"io/ioutil";
	//"encoding/json";
	sf "bitbucket.org/krepa098/gosfml2";
)

type textures struct {
	Texture []texture;
}

type texture struct {
	Name string
	Data *sf.Texture;
}

func newTextures(dirname string) *textures {
	dir, err := ioutil.ReadDir(dirname);
	if err != nil {
		log.Fatal(err);
	}
	textures := new(textures);
	textures.Texture = make([]texture, 0);
	for _, fileinfo := range(dir) {
		var tex texture;
		tex.Name = fileinfo.Name();
		tex.Data, err = sf.NewTextureFromFile(dirname + fileinfo.Name(), nil);
		if err != nil {
			log.Fatal(err);
		}
		textures.Texture = append(textures.Texture, tex);
	}
	return textures;
}

