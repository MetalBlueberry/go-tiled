package tiled

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testLoadTilesetFile = &Tileset{
	baseDir:  ".",
	Columns:  64,
	FirstGID: 0,
	Image: &Image{
		Format: "",
		Data:   nil,
		Height: 3040,
		Width:  2048,
		Source: "ProjectUtumno_full.png",
		Trans:  nil,
	},
	Margin:       0,
	Name:         "ProjectUtumno_full",
	Properties:   nil,
	Source:       "",
	SourceLoaded: true,
	Spacing:      0,
	TerrainTypes: nil,
	TileCount:    6080,
	TileHeight:   32,
	TileOffset:   nil,
	TileWidth:    32,
	TiledVersion: "1.2.3",
	Tiles: []*TilesetTile{
		{
			ID:           116,
			Type:         "door",
			Animation:    nil,
			Image:        nil,
			ObjectGroups: nil,
			Probability:  0,
			Properties:   nil,
			Terrain:      "",
		},
	},
	Version: "1.2",
}

var testLoadTileFile = &TilesetTile{
	ID: 464,
	Animation: &Animation{
		Frame: []*AnimationFrame{
			{
				Duration: 500,
				TileID:   75,
			},
			{
				Duration: 500,
				TileID:   76,
			},
		},
	},
	Image: nil,
	ObjectGroups: []*ObjectGroup{
		{
			ID:        0,
			Color:     nil,
			DrawOrder: "index",
			Name:      "",
			Objects: []*Object{
				{
					GID:        0,
					Ellipses:   nil,
					Height:     6.125,
					ID:         1,
					Name:       "",
					PolyLines:  nil,
					Polygons:   nil,
					Properties: nil,
					Rotation:   0,
					Text:       nil,
					Type:       "",
					Visible:    nil,
					Width:      32.375,
					X:          -0.25,
					Y:          17.75,
				},
			},
			OffsetX:    0,
			OffsetY:    0,
			Opacity:    0,
			Properties: nil,
			Visible:    nil,
		},
	},
}

func TestLoadTileset(t *testing.T) {
	tsxFile, err := os.Open(filepath.Join(GetAssetsDirectory(), "tilesets/testLoadTileset.tsx"))
	assert.Nil(t, err)
	defer tsxFile.Close()

	tsx, err := LoadTilesetFromReader(".", tsxFile)
	assert.Nil(t, err)

	assert.Equal(t, testLoadTilesetFile, tsx)
}

func TestSaveTileset(t *testing.T) {
	tsxFile, err := os.Open(filepath.Join(GetAssetsDirectory(), "tilesets/testLoadTileset.tsx"))
	assert.Nil(t, err)
	defer tsxFile.Close()

	buffer := &bytes.Buffer{}
	err = SaveTilesetToWriter(testLoadTilesetFile, buffer)
	assert.Nil(t, err)

	assertXMLEqual(t, tsxFile, buffer)
}

func TestLoadTile(t *testing.T) {
	tsxFile, err := os.Open(filepath.Join(GetAssetsDirectory(), "tilesets/testLoadTile.tsx"))
	assert.Nil(t, err)
	defer tsxFile.Close()

	tsx, err := LoadTilesetFromReader(".", tsxFile)
	assert.Nil(t, err)
	assert.Len(t, tsx.Tiles, 1)

	tile := tsx.Tiles[0]
	assert.Equal(t, testLoadTileFile, tile)
}

func TestSaveTile(t *testing.T) {
	tsxFile, err := os.Open(filepath.Join(GetAssetsDirectory(), "tilesets/testLoadTile.tsx"))
	assert.Nil(t, err)
	defer tsxFile.Close()

	tsx, err := LoadTilesetFromReader(".", tsxFile)
	assert.Nil(t, err)

	buffer := &bytes.Buffer{}
	xml.NewEncoder(buffer).Encode(tsx)

	tsxFile.Seek(0, 0)
	assertXMLEqual(t, tsxFile, buffer)

}

func assertXMLEqual(t *testing.T, expected io.Reader, obtained io.Reader) {
	var expec node
	var obt node
	var err error

	err = xml.NewDecoder(expected).Decode(&expec)
	assert.Nil(t, err)
	err = xml.NewDecoder(obtained).Decode(&obt)
	assert.Nil(t, err)

	assert.Equal(t, expec, obt)
}

type node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",innerxml"`
	Nodes   []node     `xml:",any"`
}

func (n *node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type inNode node

	err := d.DecodeElement((*inNode)(n), &start)
	if err != nil {
		return err
	}

	//Discard content if there are child nodes
	if len(n.Nodes) > 0 {
		n.Content = ""
	}
	return nil
}

func b(v bool) *bool {
	return &v
}
