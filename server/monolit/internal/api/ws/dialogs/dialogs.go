package dialogs

import (
	"github.com/zamatay/otus/arch/lesson-1/internal/api/dialogs"
)

type Dialogs struct {
	srv dialogs.DialogServiced
}

func NewDialogs(srv dialogs.DialogServiced) *Dialogs {
	return &Dialogs{srv: srv}
}
