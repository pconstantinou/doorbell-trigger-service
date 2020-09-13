package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/pusher/pusher-http-go"

	"encoding/json"
	"io/ioutil"

	"google.golang.org/api/option"
)

const (
	// PusherConfigPath is the location of the pusher secrets and configuration properties
	PusherConfigPath = "./pusher.json"
	// FirebaseConfigPath is the location of the firebase authentication secrets and configuration properties.
	// Note, the ./templates/entercode.html file also needs to be updated.
	FirebaseConfigPath = "./doorbell-520-firebase-adminsdk-mvx33-b87513bbce.json"
)

var formTemplate = loadTemplate("EnterCode", "./templates/entercode.html")
var confirmationTemplate = loadTemplate("Confirmation", "./templates/confirm.html")

type pusherConfig struct {
	AppID   string
	Key     string
	Secret  string
	Cluster string
	Secure  bool
	Channel string
	Event   string
}

var pusherConfigData = loadPusherConfig(PusherConfigPath)

func init() {
	http.HandleFunc("/", handleShowForm)
	http.HandleFunc("/submitCode", handleSendCode)
	http.HandleFunc("/custom_icon.png", handleIcon)
	http.HandleFunc("/fav.ico", handleIcon)
}

const customIconPath = "custom_icon.png"

func handleIcon(w http.ResponseWriter, r *http.Request) {

	file, err := os.Open(customIconPath)
	if err != nil {
		log.Printf("Failed to open %s", customIconPath)
		return
	}
	defer file.Close()

	data := make([]byte, 512)

	for {
		i, err := file.Read(data)
		if err != nil {
			break
		}
		w.Write(data[0:i])
	}
}

func loadTemplate(name, path string) *template.Template {
	templateBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Can't load template path %s for %s", path, name)
		return nil
	}
	defer log.Printf("Loaded template for %s from %s of length %d", name, path, len(templateBytes))
	return template.Must(template.New(name).Parse(string(templateBytes)))
}

func loadPusherConfig(path string) pusherConfig {
	log.Println("Loading ", path)
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Failed to read pusher config:", err)
	}

	p := pusherConfig{}
	if jsonerr := json.Unmarshal(raw, &p); jsonerr != nil {
		log.Fatal(jsonerr)
	}
	log.Print("Loaded app id:", p.AppID)
	return p
}

func handleShowForm(w http.ResponseWriter, r *http.Request) {
	log.Print("Serving the front page.")
	formTemplate.Execute(w, nil)
}

func isAuthorized(r *http.Request) (bool, *auth.UserRecord) {
	idToken := r.Form["idToken"]
	if len(idToken) != 1 || len(idToken[0]) == 0 {
		log.Print("idToken not provided")
		return false, nil
	}

	ctx := r.Context()

	opt := option.WithCredentialsFile(FirebaseConfigPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return false, nil
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing all client: %v\n", err)
		return false, nil
	}

	token, err := client.VerifyIDToken(ctx, idToken[0])
	if err != nil {
		log.Fatalf("error verifying ID token: %v. submitted token: %v\n verified: %v\ncontext: %v\n", err, idToken[0], token, ctx)
		return false, nil
	}

	user, err := client.GetUser(ctx, token.UID)
	if err != nil || user == nil {
		log.Fatalf("Failed to get user verifying ID token: %v\n", err)
		return false, nil
	}

	return true, user
}

func handleSendCode(w http.ResponseWriter, r *http.Request) {
	var inputCode string
	var ip = r.RemoteAddr

	var headers string
	for k, v := range r.Header {
		headers += k + "->" + strings.Join(v, ",")
	}

	log.Print(net.SplitHostPort(r.RemoteAddr))
	log.Print(r.Header)

	r.ParseForm()

	log.Print("Serving the front page.")
	authorized, userRecord := isAuthorized(r)

	if authorized {
		log.Print("Logging in with firebase using email:" + userRecord.Email)
		inputCode = userRecord.Email
	} else if code := r.Form["code"][0]; !strings.ContainsAny(code, ".@-") {
		inputCode = code
		log.Print("Logging in secret passcode")
		authorized = true
	} else {
		authorized = false
		log.Print("Unauthorized access.")
	}

	if authorized {
		data := map[string]string{
			"message":    inputCode,
			"ip":         ip,
			"header":     headers,
			"remoteAddr": r.RemoteAddr,
			"referrer":   r.Referer(),
			"uri":        r.RequestURI,
		}
		log.Print("Sending code to pusher app: ", pusherConfigData.AppID, " from ip ", ip)
		log.Print(data)

		client := pusher.Client{
			AppID:      pusherConfigData.AppID,
			Key:        pusherConfigData.Key,
			Secret:     pusherConfigData.Secret,
			Cluster:    pusherConfigData.Cluster,
			Secure:     pusherConfigData.Secure,
			HTTPClient: http.DefaultClient,
		}

		var err = client.Trigger(pusherConfigData.Channel, pusherConfigData.Event, data)
		if err != nil {
			log.Print("Failed to trigger event :", err)
		} else {
			log.Print("Tigger send successfully")
		}
	}
	confirmationTemplate.Execute(w, nil)
}

func main() {
	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
	// [END setting_port]
}
