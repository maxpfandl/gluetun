package alpine

import (
	"github.com/qdm12/golibs/os"
	"github.com/qdm12/golibs/os/user"
)

type Configurator interface {
	CreateUser(username string, uid int) (createdUsername string, err error)
}

type configurator struct {
	openFile os.OpenFileFunc
	osUser   user.OSUser
}

func NewConfigurator(openFile os.OpenFileFunc, osUser user.OSUser) Configurator {
	return &configurator{
		openFile: openFile,
		osUser:   osUser,
	}
}
