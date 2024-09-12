package email

import (
	"fmt"
	"mime"
	"strings"
)

type Composer struct {
	isClosed bool
	content  strings.Builder
}

func NewComposer() *Composer {
	return &Composer{}
}

func (c *Composer) AddHeader(name, value string) {
	//TODO: fix line folding >80 characters
	if c.isClosed {
		panic("composer had already written body!")
	}
	header := fmt.Sprintf("%s: %s\n", name, mime.QEncoding.Encode("utf-8", value))
	c.content.WriteString(header)
}

func (c *Composer) AddHeaderRaw(name, value string) {
	//TODO: fix line folding >80 characters
	if c.isClosed {
		panic("composer had already written body!")
	}
	header := fmt.Sprintf("%s: %s\n", name, value)
	c.content.WriteString(header)
}

func (c *Composer) Body(body string) {
	c.content.WriteString("\n")
	c.content.WriteString(body)
	c.isClosed = true
}

func (c *Composer) String() string {
	return c.content.String()
}
