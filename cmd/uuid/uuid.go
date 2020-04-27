package uuid

import (
    "ServerStatus/cmd"
    "fmt"
    uuid "github.com/satori/go.uuid"
)

func Run(p *cmd.Cmd) {
    _ = p
    id := uuid.NewV4()
    fmt.Println(id.String())
}
