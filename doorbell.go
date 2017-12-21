package main

import (
	"github.com/pusher/pusher-http-go"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"html/template"
	"log"
	"net/http"
)

var pusherConfig = GetPusherConfig()

func init() {
	http.HandleFunc("/", handleShowForm)
	http.HandleFunc("/submitCode", handleSendCode)
}

func handleShowForm(w http.ResponseWriter, r *http.Request) {
	log.Print("Serving the front page.")
	formTemplate.Execute(w, nil)
}

func handleSendCode(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	log.Print("Serving the front page.")

	inputCode := r.Form["code"][0]
	data := map[string]string{"message": inputCode}
	log.Print("Sending code")

	ctx := appengine.NewContext(r)
	var client = pusher.Client{
		AppId:      pusherConfig.AppId,
		Key:        pusherConfig.Key,
		Secret:     pusherConfig.Secret,
		Cluster:    pusherConfig.Cluster,
		Secure:     pusherConfig.Secure,
		HttpClient: urlfetch.Client(ctx),
	}

	var _, err = client.Trigger(pusherConfig.Channel, pusherConfig.Event, data)
	if err != nil {
		log.Print("Failed to trigger event :", err)
	} else {
		log.Print("Tigger send successfully")
	}
	var appError = confirmationTemplate.Execute(w, nil)
	if appError != nil {
		log.Print("Failed to write template:", appError)
	}
}

var confirmationTemplate = template.Must(template.New("Response").Parse(`
<html>
<head><title>Enter Code</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="stylesheet" href="https://www.w3schools.com/w3css/4/w3mobile.css">
</head>
<script>
window.onload = function() {
  document.getElementById("code").focus();
};
</script>
<body>
<div class="w3-container">
  <h1 style="text-align: center;">Doorman</h1>
</div>

<div class="w3-cell-row">
  <div class="w3-cell">
    
  </div>
  <div class="w3-cell w3-container">
<form action="/" method="GET">
<div>
<div style="text-align: center;">Code sent.</div>
<input type="submit" value="Send Again"  style="width: 100%; border: 4px; padding: 4px; margin: 4px;">
</div>
</form>
  </div>
</div>

</body></html>

`))

var formTemplate = template.Must(template.New("EnterCode").Parse(`
<html>
<head><title>Enter Code</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="stylesheet" href="https://www.w3schools.com/w3css/4/w3mobile.css">
</head>
<script>
window.onload = function() {
  document.getElementById("code").focus();
};

function validateForm() {
    var x = document.forms["codeForm"]["code"].value;
    if (x == "") {
        alert("Provide a passcode.");
        return false;
    }
}

</script>
<body>
<div class="w3-container">
  <h1 style="text-align: center;">Doorman</h1>
</div>

<div class="w3-cell-row">
  <div class="w3-cell">
    
  </div>
  <div class="w3-cell w3-container">
<form name="codeForm" action="/submitCode" method="POST" onsubmit="return validateForm()">

<div>
<input id="code" type="tel" name="code" autofocus="autofocus" 
	   placeholder="Enter Passcode"
       style="text-align: center;width: 100%; border: 4px; padding: 4px; margin: 4px;border-color: black;background-color: SkyBlue"><p/>
<input type="submit" value="Unlock"  style="width: 100%; border: 4px; padding: 4px; margin: 4px;">
</div>
</form>
  </div>
</div>

</body></html>
`))
