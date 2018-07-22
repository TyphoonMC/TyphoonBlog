package main

import (
	t "github.com/TyphoonMC/TyphoonCore"
	"log"
	"io/ioutil"
	"github.com/a8m/mark"
	"fmt"
	"bytes"
)

func main() {
	core := t.Init()
	core.SetBrand("typhoonblog")

	core.On(func(e *t.PlayerJoinEvent) {
		articlesDir := "./articles"

		files, err := ioutil.ReadDir(articlesDir)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			data, err := ioutil.ReadFile(articlesDir + "/" + file.Name())
			if err != nil {
				panic(err)
			}
			fmt.Println(MinecraftRender(string(data)))
			e.Player.SendBukkitMessage(MinecraftRender(string(data)))
		}
	})

	core.Start()
}

func MinecraftRender(article string) string {
	m := mark.New(article, &mark.Options{
		Smartypants: true,
		Fractions:   true,
	})
	m.AddRenderFn(mark.NodeHeading, func(node mark.Node) string {
		h, _ := node.(*mark.HeadingNode)
		format := ""
		switch h.Level {
		case 1: format = "&6&l&n"
		case 2: format = "&6&l"
		case 3: format = "&6&n"
		case 4: format = "&e&n"
		}
		return fmt.Sprintf("%s%s&r\n", format, h.Text)
	})
	m.AddRenderFn(mark.NodeParagraph, func(node mark.Node) string {
		p, _ := node.(*mark.ParagraphNode)
		buff := bytes.NewBufferString("&f")
		for _, n := range p.Nodes {
			buff.WriteString(n.Render() + "\n")
		}
		buff.WriteString("&r\n")
		return buff.String()
	})
	m.AddRenderFn(mark.NodeText, func(node mark.Node) string {
		p, _ := node.(*mark.TextNode)
		return p.Text
	})
	m.AddRenderFn(mark.NodeHTML, func(node mark.Node) string {
		return ""
	})
	return m.Render()
}