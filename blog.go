package main

import (
	"bytes"
	"fmt"
	t "github.com/TyphoonMC/TyphoonCore"
	"github.com/a8m/mark"
	"io/ioutil"
	"log"
	"strings"
)

type Article struct {
	Title    string
	Author   string
	Date     string
	Content  string
	Metadata map[string]string
}

var (
	articles []Article
)

func main() {
	readArticles()

	core := t.Init()
	core.SetBrand("typhoonblog")

	core.On(func(e *t.PlayerJoinEvent) {
		e.Player.SendMessage(t.ChatMessage("Welcome to my blog !"))
	})

	core.DeclareCommand(t.CommandNodeLiteral("article",
		[]*t.CommandNode{
			t.CommandNodeLiteral("read",
				[]*t.CommandNode{
					t.CommandNodeLiteral("all",
						nil,
						func(player *t.Player, args []string) {
							SendArticles(player)
						},
					),
					t.CommandNodeArgument("Article Id",
						nil,
						&t.CommandParserInteger{
							t.OptInteger{
								true,
								0,
							},
							t.OptInteger{
								true,
								int32(len(articles) - 1),
							},
						},
						func(player *t.Player, args []string) {
							SendArticles(player)
						},
					),
				},
				nil,
			),
			t.CommandNodeLiteral("search",
				[]*t.CommandNode{
					t.CommandNodeArgument("sentence",
						nil,
						&t.CommandParserString{
							t.CommandParserStringFormatGreedyPhrase,
						},
						func(player *t.Player, args []string) {
							fmt.Println(args)
						},
					),
				},
				nil,
			),
			t.CommandNodeLiteral("list",
				[]*t.CommandNode{},
				func(player *t.Player, args []string) {
					player.SendMessage(t.ChatMessage("Articles:"))

					for i, article := range articles {
						m := t.ChatMessage(article.Title)
						m.SetColor(&t.ChatColorGold)
						m.SetClickEvent(
							t.ChatClickRunCommand(fmt.Sprintf("/article read %d", i)),
						)
						m.SetHoverEvent(t.ChatHoverMessage([]t.IChatComponent{
							t.ChatMessage(article.Title),
							t.ChatMessage(article.Author),
							t.ChatMessage(article.Date),
						}))
						player.SendMessage(m)
					}
				},
			),
		},
		nil,
	))

	core.Start()
}

func readArticles() {
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

		meta := generateMetadata(string(data))
		if _, ok := meta["title"]; ok {
			if _, ok = meta["author"]; ok {
				if _, ok = meta["date"]; ok {
					articles = append(articles, Article{
						meta["title"],
						meta["author"],
						meta["date"],
						string(data),
						meta,
					})
				}
			}
		}
	}
}

func SendArticles(player *t.Player) {
	for _, article := range articles {
		player.SendBukkitMessage(MinecraftRender(article.Content))
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

func generateMetadata(article string) map[string]string {
	meta := make(map[string]string)
	m := mark.New(article, &mark.Options{
		Smartypants: true,
		Fractions:   true,
	})
	m.AddRenderFn(mark.NodeHTML, func(node mark.Node) string {
		rnd := node.Render()
		if strings.HasPrefix(rnd, "<!--") &&
			strings.HasSuffix(rnd, "-->") {
			d := strings.Split(rnd[4:len(rnd)-3], ":")
			if len(d) >= 2 {
				meta[strings.Trim(d[0], " ")] = strings.Trim(d[1], " ")
			}
		}
		return ""
	})
	m.Render()
	return meta
}
