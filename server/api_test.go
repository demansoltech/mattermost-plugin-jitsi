package main

import (
	"strings"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/require"
)

func TestRedirect(t *testing.T) {
	p := Plugin{
		configuration: &configuration{
			JitsiURL:          "http://test",
			JitsiEmbedded:     false,
			JitsiNamingScheme: "mattermost",
			JitsiAppSecret:    "test-secret",
			JitsiAppID:        "test-appidjitsi",
		},
		botID: "test-bot-id",
	}
	site := "http://mattermostserver"
	servicesetting := model.ServiceSettings{SiteURL: &site}

	t.Run("Redirect with Valid JWT", func(t *testing.T) {
		apiMock := plugintest.API{}
		defer apiMock.AssertExpectations(t)
		p.SetAPI(&apiMock)
		p.callback = CallbackValidation{
			room:   "validroom",
			UserID: "test-user"}

		apiMock.On("GetUser", "test-user").Return(&model.User{Id: "test-user", Locale: "en"}, nil)
		apiMock.On("GetConfig").Return(&model.Config{ServiceSettings: servicesetting})

		result := p.handleRedirectURL("validroom")
		require.True(t, strings.Contains(result, "validroom"), "Check Path should have room name")
		require.True(t, strings.Contains(result, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"), "url has JWT header")
	})

	t.Run("Redirect to site url for login authentication", func(t *testing.T) {
		apiMock := plugintest.API{}
		defer apiMock.AssertExpectations(t)
		p.SetAPI(&apiMock)
		p.callback = CallbackValidation{
			room:   "validroom",
			UserID: "test-user"}

		apiMock.On("GetUser", "test-user").Return(&model.User{Id: "test-user", Locale: "en"}, nil)
		apiMock.On("GetConfig").Return(&model.Config{ServiceSettings: servicesetting})

		result := p.handleRedirectURL("invalidroom")
		require.Equal(t, site, result)
	})

	t.Run("Redirect with invalid JWT, error internal", func(t *testing.T) {
		apiMock := plugintest.API{}
		defer apiMock.AssertExpectations(t)
		p.SetAPI(&apiMock)
		p.callback = CallbackValidation{
			room:   "validroom",
			UserID: "test-user"}
		apiMock.On("GetUser", "test-user").Return(&model.User{Id: "test-user", Locale: "en"}, nil)
		apiMock.On("GetConfig").Return(&model.Config{ServiceSettings: servicesetting})

		p.configuration.JitsiAppSecret = "" // error if app secret is empty string

		result := p.handleRedirectURL("validroom")
		require.True(t, strings.Contains(result, "invalidjwttoken"), "error internal, should send invalid jwt token")
	})
}
