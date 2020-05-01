package uuid

import (
	"fmt"
	"github.com/flxxyz/ServerStatus/cmd"
	uuid "github.com/satori/go.uuid"
)

func Run(p *cmd.Cmd) {
	_ = p
	id := uuid.NewV4()
	fmt.Println(id.String())
}
