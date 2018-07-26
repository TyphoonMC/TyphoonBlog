package main

import (
	"bytes"
	"fmt"
	t "github.com/TyphoonMC/TyphoonCore"
	"github.com/a8m/mark"
	"io/ioutil"
	"log"
)

func main() {
	core := t.Init()
	core.SetBrand("typhoonblog")

	core.On(func(e *t.PlayerJoinEvent) {
		e.Player.SendMessage(t.ChatMessage("Welcome to my blog !"))
	})

	core.DeclareCommand(t.CommandNodeLiteral(
		"article",
		[]*t.CommandNode{
			t.CommandNodeLiteral(
				"read",
				[]*t.CommandNode{
					t.CommandNodeLiteral(
						"all",
						nil,
						func(player *t.Player, args []string) {
							SendArticles(player)
						},
					),
				},
				nil,
			),
		},
		nil,
	))

	core.Start()
}

func SendArticles(player *t.Player) {
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
		player.SendBukkitMessage(MinecraftRender(string(data)))
	}
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
		case 1:
			format = "&6&l&n"
		case 2:
			format = "&6&l"
		case 3:
			format = "&6&n"
		case 4:
			format = "&e&n"
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
