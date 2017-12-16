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
<body>
Code send. <a href="/">Send again</a>
</body></html>
`))

var formTemplate = template.Must(template.New("EnterCode").Parse(`
<html>
<head><title>Enter Code</title></head>
<script>
window.onload = function() {
  document.getElementById("code").focus();
};
</script>
<body>
<form action="/submitCode" method="POST">
<center>
<lable>Code</label>
<input id="code" type="text" name="code" autofocus="autofocus"><br/>
<input type="submit" value="Unlock">
</center>
</form>
</body></html>
`))
