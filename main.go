package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

type IconImage struct {
	Size     string `json:"size"`
	Idiom    string `json:"idiom"`
	Filename string `json:"filename"`
	Scale    string `json:"scale"`
}

type Contents struct {
	Images []IconImage `json:"images"`
	Info   struct {
		Version int    `json:"version"`
		Author  string `json:"author"`
	} `json:"info"`
}

func main() {
	// CLI引数
	input := flag.String("input", "", "Path to the input 1024x1024 PNG file")
	output := flag.String("output", "AppIcon.appiconset", "Output directory for the .appiconset")
	flag.Parse()

	if *input == "" {
		log.Fatal("❌ --input is required (e.g., --input icon.png)")
	}

	// 出力先フォルダを作成
	err := os.MkdirAll(*output, 0755)
	if err != nil {
		log.Fatalf("Failed to create output dir: %v", err)
	}

	// 画像読み込み
	img, err := imaging.Open(*input)
	if err != nil {
		log.Fatalf("Failed to open input image: %v", err)
	}

	// 必要なアイコンサイズ一覧
	sizes := []struct {
		pointSize int
		scale     int
	}{
		{16, 1}, {16, 2},
		{32, 1}, {32, 2},
		{128, 1}, {128, 2},
		{256, 1}, {256, 2},
		{512, 1}, {512, 2},
	}

	var images []IconImage

	for _, s := range sizes {
		width := s.pointSize * s.scale
		filename := fmt.Sprintf("icon_%dx%d", s.pointSize, s.pointSize)
		if s.scale == 2 {
			filename += "@2x"
		}
		filename += ".png"

		dst := imaging.Resize(img, width, width, imaging.Lanczos)
		savePath := filepath.Join(*output, filename)
		if err := imaging.Save(dst, savePath); err != nil {
			log.Fatalf("Failed to save resized image: %v", err)
		}

		images = append(images, IconImage{
			Size:     fmt.Sprintf("%dx%d", s.pointSize, s.pointSize),
			Idiom:    "mac",
			Filename: filename,
			Scale:    fmt.Sprintf("%dx", s.scale),
		})
	}

	contents := Contents{
		Images: images,
	}
	contents.Info.Author = "xcode"
	contents.Info.Version = 1

	jsonBytes, _ := json.MarshalIndent(contents, "", "  ")
	jsonPath := filepath.Join(*output, "Contents.json")
	os.WriteFile(jsonPath, jsonBytes, 0644)

	fmt.Printf("✅ .appiconset created at: %s\n", *output)
}
