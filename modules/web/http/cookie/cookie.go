package cookie

import (
	"net/http"

	"github.com/710leo/urlooker/modules/web/g"
	"github.com/gorilla/securecookie"
)

var SecureCookie *securecookie.SecureCookie

func Init() {
	var hashKey = []byte(g.Config.Http.Secret)
	var blockKey = []byte(nil)
	SecureCookie = securecookie.New(hashKey, blockKey)
}

func ReadUser(r *http.Request) (id int64, name string, found bool) {
	if cookie, err := r.Cookie("u"); err == nil {
		value := make(map[string]interface{})
		if err = SecureCookie.Decode("u", cookie.Value, &value); err == nil {
			id = value["id"].(int64)
			name = value["name"].(string)
			if id == 0 || name == "" {
				return
			} else {
				found = true
				return
			}
		}
	}
	return
}

func WriteUser(w http.ResponseWriter, id int64, name string) error {
	value := make(map[string]interface{})
	value["id"] = id
	value["name"] = name
	encoded, err := SecureCookie.Encode("u", value)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "u",
		Value:    encoded,
		Path:     "/",
		MaxAge:   3600 * 24 * 7,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	return nil
}

func RemoveUser(w http.ResponseWriter) error {
	value := make(map[string]interface{})
	value["id"] = ""
	value["name"] = ""
	encoded, err := SecureCookie.Encode("u", value)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:   "u",
		Value:  encoded,
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	return nil
}
