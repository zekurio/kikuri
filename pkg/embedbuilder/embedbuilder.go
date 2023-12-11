package embedbuilder

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// EmbedBuilder is a helper struct for quickly creating embeds.
type EmbedBuilder struct {
	emb *discordgo.MessageEmbed
}

// New creates a new EmbedBuilder instance.
func New() *EmbedBuilder {
	return &EmbedBuilder{
		emb: &discordgo.MessageEmbed{},
	}
}

func (q *EmbedBuilder) SetTitle(title string) *EmbedBuilder {
	q.emb.Title = title
	return q
}

func (q *EmbedBuilder) SetDescription(description string) *EmbedBuilder {
	q.emb.Description = description
	return q
}

func (q *EmbedBuilder) SetURL(url string) *EmbedBuilder {
	q.emb.URL = url
	return q
}

func (q *EmbedBuilder) SetColor(color int) *EmbedBuilder {
	q.emb.Color = color
	return q
}

func (q *EmbedBuilder) SetAuthor(name, iconURL, pURL, url string) *EmbedBuilder {
	q.emb.Author = &discordgo.MessageEmbedAuthor{
		Name:         name,
		IconURL:      iconURL,
		URL:          url,
		ProxyIconURL: pURL,
	}
	return q
}

func (q *EmbedBuilder) SetThumbnail(url, pURL string, width, height int) *EmbedBuilder {
	q.emb.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:      url,
		ProxyURL: pURL,
		Width:    width,
		Height:   height,
	}
	return q
}

func (q *EmbedBuilder) SetImage(url, pURL string, width, height int) *EmbedBuilder {
	q.emb.Image = &discordgo.MessageEmbedImage{
		URL:      url,
		ProxyURL: pURL,
		Width:    width,
		Height:   height,
	}
	return q
}

func (q *EmbedBuilder) SetFooter(text, iconURL, pURL string) *EmbedBuilder {
	q.emb.Footer = &discordgo.MessageEmbedFooter{
		Text:         text,
		IconURL:      iconURL,
		ProxyIconURL: pURL,
	}
	return q
}

func (q *EmbedBuilder) SetTimestamp(timestamp ...time.Time) *EmbedBuilder {

	if len(timestamp) > 0 {
		q.emb.Timestamp = timestamp[0].Format("02/01/2006 15:04:05")
	} else {
		q.emb.Timestamp = time.Now().Format("02/01/2006 15:04:05")
	}

	return q
}

func (q *EmbedBuilder) Build() *discordgo.MessageEmbed {
	return q.emb
}

func (q *EmbedBuilder) AddField(name, value string, inline ...bool) *EmbedBuilder {

	if value == "" {
		value = "`nil`"
	}

	q.emb.Fields = append(q.emb.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: len(inline) > 0 && inline[0],
	})
	return q
}
