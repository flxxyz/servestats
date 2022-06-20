package uuid

import (
	"fmt"
	"github.com/flxxyz/servestats/cmd"
	uuid "github.com/satori/go.uuid"
)

func Run(p *cmd.Cmd) {
	_ = p
	id := uuid.NewV4()
	fmt.Println(id.String())
}
