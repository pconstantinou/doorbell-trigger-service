import urllib.request, json 
with urllib.request.urlopen("https://www.purpleair.com/json?show=2911") as url:
    data = json.loads(url.read().decode())
    print(data['results'][0]['pm2_5_atm'])