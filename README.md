# 2022-10-11-oauth

```bash
# Get an an authorization code
$ curl -i 'http://localhost:8080/oauth2/authorize?response_type=code&client_id=drat' -H 'Cookie: drat_session=matts'
HTTP/1.1 302 Found
Content-Length: 0
Date: Wed, 12 Oct 2022 16:22:56 GMT
Location: http://localhost:3000?code=NGU3NGZJODATZJU0NY0ZNWZJLTLMOWETMGNHZTM5ZTQ4OTIY
Vary: Origin

# Use the authorization code to get an access token
$ curl -i 'http://localhost:8080/oauth2/token?grant_type=authorization_code&client_id=drat&code=NGU3NGZJODATZJU0NY0ZNWZJLTLMOWETMGNHZTM5ZTQ4OTIY&redirect_uri=http://localhost:3000'
HTTP/1.1 200 OK
Cache-Control: no-store
Content-Length: 175
Content-Type: application/json;charset=UTF-8
Date: Wed, 12 Oct 2022 16:26:50 GMT
Pragma: no-cache
Vary: Origin

{"access_token":"NZY0NTVKMDYTMWIYZC0ZYJE3LWFKODATNDBHNTQ5ZTHIZJLJ","expires_in":7200,"refresh_token":"YJNKZTKXYTKTZWY2MC01NWNJLWFHNTETNGU5MDVKZMU4YJNM","token_type":"Bearer"}

# Include state parameter to verify authenticity
curl -i 'http://localhost:8080/oauth2/authorize?response_type=code&client_id=drat&state=asdf' -H 'Cookie: drat_session=matts'
HTTP/1.1 302 Found
Content-Length: 0
Date: Wed, 12 Oct 2022 16:32:28 GMT
Location: http://localhost:3000?code=MTRMYJLIMJETMJUXYS0ZZDFLLWFIMWITODM5OWM4NZU1NWIY&state=asdf
Vary: Origin
```
