package jiracli

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"github.com/jinzhu/copier"
	shellquote "github.com/kballard/go-shellquote"
	"github.com/tidwall/gjson"
	"gopkg.in/AlecAivazis/survey.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/coryb/yaml.v2"
	logging "gopkg.in/op/go-logging.v1"
)

type Exit struct {
	Code int
}

type GlobalOptions struct {
	// AuthenticationMethod is the method we use to authenticate with the jira serivce. Possible values are "api-token" or "session".
	// The default is "api-token" when the service endpoint ends with "atlassian.net", otherwise it "session".  Session authentication
	// will promt for user password and use the /auth/1/session-login endpoint.
	AuthenticationMethod figtree.StringOption `yaml:"authentication-method,omitempty" json:"authentication-method,omitempty"`

	// Endpoint is the URL for the Jira service.  Something like: https://go-jira.atlassian.net
	Endpoint figtree.StringOption `yaml:"endpoint,omitempty" json:"endpoint,omitempty"`

	// Insecure will allow you to connect to an https endpoint with a self-signed SSL certificate
	Insecure figtree.BoolOption `yaml:"insecure,omitempty" json:"insecure,omitempty"`

	// Login is the id used for authenticating with the Jira service.  For "api-token" AuthenticationMethod this is usually a
	// full email address, something like "user@example.com".  For "session" AuthenticationMethod this will be something
	// like "user", which by default will use the same value in the `User` field.
	Login figtree.StringOption `yaml:"login,omitempty" json:"login,omitempty"`

	// PasswordSource specificies the method that we fetch the password.  Possible values are "keyring" or "pass".
	// If this is unset we will just prompt the user.  For "keyring" this will look in the OS keychain, if missing
	// then prompt the user and store the password in the OS keychain.  For "pass" this will look in the PasswordDirectory
	// location using the `pass` tool, if missing prompt the user and store in the PasswordDirectory
	PasswordSource figtree.StringOption `yaml:"password-source,omitempty" json:"password-source,omitempty"`

	// PasswordDirectory is only used for the "pass" PasswordSource.  It is the location for the encrypted password
	// files used by `pass`.  Effectively this overrides the "PASSWORD_STORE_DIR" environment variable
	PasswordDirectory figtree.StringOption `yaml:"password-directory,omitempty" json:"password-directory,omitempty"`

	// PasswordName is the the name of the password key entry stored used with PasswordSource `pass`.
	PasswordName figtree.StringOption `yaml:"password-name,omitempty" json:"password-name,omitempty"`

	// Quiet will lower the defalt log level to suppress the standard output for commands
	Quiet figtree.BoolOption `yaml:"quiet,omitempty" json:"quiet,omitempty"`

	// SocksProxy is used to configure the http client to access the Endpoint via a socks proxy.  The value
	// should be a ip address and port string, something like "127.0.0.1:1080"
	SocksProxy figtree.StringOption `yaml:"socksproxy,omitempty" json:"socksproxy,omitempty"`

	// UnixProxy is use to configure the http client to access the Endpoint via a local unix domain socket used
	// to proxy requests
	UnixProxy figtree.StringOption `yaml:"unixproxy,omitempty" json:"unixproxy,omitempty"`

	// User is use to represent the user on the Jira service.  This can be different from the username used to
	// authenticate with the service.  For example when using AuthenticationMethod `api-token` the Login is
	// typically an email address like `username@example.com` and the User property would be someting like
	// `username`  The User property is used on Jira service API calls that require a user to associate with
	// an Issue (like assigning a Issue to yourself)
	User figtree.StringOption `yaml:"user,omitempty" json:"user,omitempty"`
}

func RegisterCommand(regEntry CommandRegistry) {
	globalCommandRegistry = append(globalCommandRegistry, regEntry)
}

func (o *GlobalOptions) AuthMethod() string {
	if strings.Contains(o.Endpoint.Value, ".atlassian.net") && o.AuthenticationMethod.Source == "default" {
		return "api-token"
	}
	return o.AuthenticationMethod.Value
}

func register(app *kingpin.Application, o *oreo.Client, fig *figtree.FigTree) {
	globals := GlobalOptions{
		User:                 figtree.NewStringOption(os.Getenv("USER")),
		AuthenticationMethod: figtree.NewStringOption("session"),
	}
	app.Flag("endpoint", "Base URI to use for Jira").Short('e').SetValue(&globals.Endpoint)
	app.Flag("insecure", "Disable TLS certificate verification").Short('k').SetValue(&globals.Insecure)
	app.Flag("quiet", "Suppress output to console").Short('Q').SetValue(&globals.Quiet)
	app.Flag("unixproxy", "Path for a unix-socket proxy").SetValue(&globals.UnixProxy)
	app.Flag("socksproxy", "Address for a socks proxy").SetValue(&globals.SocksProxy)
	app.Flag("user", "user name used within the Jira service").Short('u').SetValue(&globals.User)
	app.Flag("login", "login name that corresponds to the user used for authentication").SetValue(&globals.Login)

	o = o.WithPreCallback(func(req *http.Request) (*http.Request, error) {
		if globals.AuthMethod() == "api-token" {
			// need to set basic auth header with user@domain:api-token
			token := globals.GetPass()
			authHeader := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", globals.Login.Value, token))))
			req.Header.Add("Authorization", authHeader)
		}
		return req, nil
	})

	o = o.WithPostCallback(func(req *http.Request, resp *http.Response) (*http.Response, error) {
		if globals.AuthMethod() == "session" {
			authUser := resp.Header.Get("X-Ausername")
			if authUser == "" || authUser == "anonymous" {
				// preserve the --quiet value, we need to temporarily disable it so
				// the normal login output is surpressed
				defer func(quiet bool) {
					globals.Quiet.Value = quiet
				}(globals.Quiet.Value)
				globals.Quiet.Value = true

				// we are not logged in, so force login now by running the "login" command
				app.Parse([]string{"login"})

				// rerun the original request
				return o.Do(req)
			}
		} else if globals.AuthMethod() == "api-token" && resp.StatusCode == 401 {
			globals.SetPass("")
			return o.Do(req)
		}
		return resp, nil
	})

	for _, command := range globalCommandRegistry {
		copy := command
		commandFields := strings.Fields(copy.Command)
		var appOrCmd kingpinAppOrCommand = app
		if len(commandFields) > 1 {
			for _, name := range commandFields[0 : len(commandFields)-1] {
				tmp := appOrCmd.GetCommand(name)
				if tmp == nil {
					tmp = appOrCmd.Command(name, "")
				}
				appOrCmd = tmp
			}
		}

		cmd := appOrCmd.Command(commandFields[len(commandFields)-1], copy.Entry.Help)
		LoadConfigs(cmd, fig, &globals)
		cmd.PreAction(func(_ *kingpin.ParseContext) error {
			if globals.Insecure.Value {
				transport := &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				}
				o = o.WithTransport(transport)
			}
			if globals.UnixProxy.Value != "" {
				o = o.WithTransport(unixProxy(globals.UnixProxy.Value))
			} else if globals.SocksProxy.Value != "" {
				o = o.WithTransport(socksProxy(globals.SocksProxy.Value))
			}
			if globals.AuthMethod() == "api-token" {
				o = o.WithCookieFile("")
			}
			if globals.Login.Value == "" {
				globals.Login = globals.User
			}
			return nil
		})

		for _, alias := range copy.Aliases {
			cmd = cmd.Alias(alias)
		}
		if copy.Default {
			cmd = cmd.Default()
		}
		if copy.Entry.UsageFunc != nil {
			copy.Entry.UsageFunc(fig, cmd)
		}

		cmd.Action(func(_ *kingpin.ParseContext) error {
			if logging.GetLevel("") > logging.DEBUG {
				o = o.WithTrace(true)
			}
			return copy.Entry.ExecuteFunc(o, &globals)
		})
	}
}

func LoadConfigs(cmd *kingpin.CmdClause, fig *figtree.FigTree, opts interface{}) {
	cmd.PreAction(func(_ *kingpin.ParseContext) error {
		os.Setenv("JIRA_OPERATION", cmd.FullCommand())
		// load command specific configs first
		if err := fig.LoadAllConfigs(strings.Join(strings.Fields(cmd.FullCommand()), "_")+".yml", opts); err != nil {
			return err
		}

		// then load generic configs if not already populated above
		return fig.LoadAllConfigs("config.yml", opts)
	})
}
