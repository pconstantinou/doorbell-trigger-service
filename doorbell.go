package main

import (
	"github.com/pusher/pusher-http-go"
	"html/template"
	"log"
	"net/http"
)

var pusherConfig = GetPusherConfig()

var client = pusher.Client{
	AppId:   pusherConfig.AppId,
	Key:     pusherConfig.Key,
	Secret:  pusherConfig.Secret,
	Cluster: pusherConfig.Cluster,
	Secure:  pusherConfig.Secure,
}

func main() {
	http.HandleFunc("/", handleShowForm)
	http.HandleFunc("/submitCode", handleSendCode)
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
	client.Trigger(pusherConfig.Channel, pusherConfig.Event, data)
	confirmationTemplate.Execute(w, nil)
}

var confirmationTemplate = template.Must(template.New("Response").Parse(`
<html>
<head><title>Send</title></head>
<meta name="viewport" content="width=device-width,initial-scale=1">
<body>
Code send. <a href="/">Send again</a>
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
</script>
<body>
<div class="w3-container">
  <h1 style="text-align: center;">Doorman</h1>
</div>

<div class="w3-cell-row">
  <div class="w3-cell">
    
  </div>
  <div class="w3-cell w3-container">
<form action="/submitCode" method="POST">

<div>
<input id="code" type="tel" name="code" autofocus="autofocus" 
       style="text-align: center;width: 100%; border: 4px; padding: 4px; margin: 4px;border-color: black;background-color: SkyBlue"><p/>
<input type="submit" value="Unlock"  style="width: 100%; border: 4px; padding: 4px; margin: 4px;">
</div>
</form>
  </div>
</div>

</body></html>
`))
