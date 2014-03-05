package main

import (
	"encoding/json"
	"fmt"
	"github.com/Mistobaan/mixpanels-go"
	"github.com/stretchr/jsonblend/blend"
	"io/ioutil"
	"net/http"
)

func writeError(
	mp *mixpanel.Mixpanel,
	ID string,
	w http.ResponseWriter,
	str string,
	args ...interface{}) {
	errStr := fmt.Sprintf(str, args)
	fmt.Fprintf(w, str, args)
	mp.Track(ID, "error", &mixpanel.P{"error message": errStr})
}

func handler(w http.ResponseWriter, r *http.Request) {

	consumer := mixpanel.NewBuffConsumer(8)
	mp := mixpanel.NewMixpanelWithConsumer(
		"c6d3b1ae675719a889a0259abef2bdd5",
		consumer,
	)

	defer consumer.Flush()

	ID := r.RemoteAddr
	cookie, err := r.Cookie("mp_c6d3b1ae675719a889a0259abef2bdd5_mixpanel")
	if err == nil {
		ID = cookie.Value
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(mp, ID, w, "Unable to read request body. Error: %s", err.Error())
		return
	}

	var items []map[string]interface{}
	err = json.Unmarshal(body, &items)
	if err != nil {
		writeError(mp, ID, w, "Unable to unmarshal JSON from body. Error: %s", err.Error())
		return
	}

	dest := map[string]interface{}{}
	count := 0
	for _, item := range items {
		err := blend.Blend(item, dest)
		if err != nil {
			writeError(mp, ID, w, "Unable to blend JSON. Error: %s", err.Error())
			return
		}
		count++
	}

	mp.Track(ID, "lines-posted", &mixpanel.P{"count": count})
	blended, err := json.Marshal(dest)
	if err != nil {
		writeError(mp, ID, w, "Unable to marshal blended map to JSON. Error: %s", err.Error())
		return
	}

	mp.Track(ID, "blend", &mixpanel.P{"result-size": len(blended)})
	fmt.Fprintf(w, "%s", blended)

}

func main() {
	fmt.Println("JSONBlend server started")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", muxWrapper{http.DefaultServeMux})
}
