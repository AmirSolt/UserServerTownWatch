package notif

import (
	"basedpocket/base"
	"basedpocket/cmodels"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
)

type NotifCreateManyParams struct {
	params []NotifCreateParams
}

type NotifCreateParams struct {
	UserID   string `db:"user_id" json:"user_id"`
	Subject  string `db:"subject" json:"subject"`
	BodyHTML string `db:"body_html" json:"body_html"`
}

func handleCreateNotifs(app core.App, ctx echo.Context, env *base.Env) error {

	var notifCreateManyParams NotifCreateManyParams
	err := ctx.Bind(&notifCreateManyParams)
	if err != nil {
		eventID := sentry.CaptureException(err)
		cerr := &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
		return ctx.String(http.StatusInternalServerError, cerr.Error.Error())
	}

	var notifs []cmodels.Notif
	for _, param := range notifCreateManyParams.params {
		notifs = append(notifs, cmodels.Notif{
			User:     param.UserID,
			Subject:  param.Subject,
			BodyHTML: param.BodyHTML,
		})

	}

	if err := cmodels.Save(app, notifs); err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error.Error())
	}

	return ctx.NoContent(http.StatusOK)
}
