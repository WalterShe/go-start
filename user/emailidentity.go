package user

import (
	"fmt"
	"github.com/ungerik/go-mail"
	"github.com/ungerik/go-start/mgo/bson"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
	"net/url"
)

///////////////////////////////////////////////////////////////////////////////
// EmailIdentity

type EmailIdentity struct {
	//	mongo.SubDocumentBase
	Address          model.Email
	Description      model.String
	Confirmed        model.DateTime
	ConfirmationCode model.String
}

// EmailIdentity has to be saved after a successful call because the confirmation code could have changed
// confirmationPage needs to be a page with one URL parameter
func (self *EmailIdentity) SendConfirmationEmail(response *view.Response, confirmationURL view.URL) <-chan error {
	errChan := make(chan error, 1)

	confirmationCode := self.ConfirmationCode.Get()
	if confirmationCode == "" {
		confirmationCode = bson.NewObjectId().Hex()
		self.ConfirmationCode.SetString(confirmationCode)
	}

	subject := fmt.Sprintf(Config.ConfirmationMessage.EmailSubject, view.Config.SiteName)
	confirm := confirmationURL.URL(response) + "?code=" + url.QueryEscape(confirmationCode)
	message := fmt.Sprintf(Config.ConfirmationMessage.EmailMessage, view.Config.SiteName, confirm)

	go func() {
		errChan <- email.NewBriefMessage(subject, message, self.Address.Get()).Send()
		close(errChan)
	}()

	return errChan
}

func (self *EmailIdentity) MailtoURL() string {
	return "mailto:" + self.Address.Get()
}

func (self *EmailIdentity) URL(response *view.Response) string {
	return self.MailtoURL()
}

func (self *EmailIdentity) LinkContent(response *view.Response) view.View {
	return view.Escape(self.Address.Get())
}

func (self *EmailIdentity) LinkTitle(response *view.Response) string {
	return self.Address.Get()
}

func (self *EmailIdentity) LinkRel(response *view.Response) string {
	return ""
}
