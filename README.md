
To get started, create an account at pusher.com and get an API key

```
cp pusher.json.default to pusher.json
```

Create a Firebase project with a Standard AppEngine configuration. Follow the instructions here: https://firebase.google.com/docs/admin/setup#add_firebase_to_your_app
to create the Firebase project. Download the Private Key for the Firebase Admin SDK and save in with the name:

```
doorbell-dev-firebase-adminsdk.json
```
This file allows the Go App Engine App to access information about the authenticated users.


Obtain the HTML/Javascript Initialize Firebase snippet under Authentication > WEB SETUP. Place it at the top of ./templates/entercode.html.

Enable Google and other desired sign in methods Authentication > Sign-In Moethod.



Test the code locally:
```
make test
```
A local development environment accessible at http://localhost:8080/

Once everything looks good, push it to Google cloud:
```
gcloud app deploy
```
or
```
make deploy
```

