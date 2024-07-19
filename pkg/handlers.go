package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/billzayy/url-shorten/db"
	"github.com/billzayy/url-shorten/types"
)

func GenerateURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]any{
			"statusCode": http.StatusMethodNotAllowed,
			"data":       errors.New("method not allowed"),
			"message":    "Failure",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req types.Request
	body, err := io.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &req)

	if err != nil {
		panic(err)
	}

	rdb, err := db.InitRedisClient()

	if err != nil {
		fmt.Println(err.Error())
	}
	defer rdb.Close()

	hashUrl := Base62Encode(rand.Uint64())

	shortedUrl := fmt.Sprintf("http://localhost:8080/redirect/%s", hashUrl)

	err = rdb.Set(hashUrl, req.Url, 40*time.Second).Err()

	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"statusCode": http.StatusInternalServerError,
			"data":       err.Error(),
			"message":    "Failure",
		})
		return
	}

	t := time.Now().Add(40 * time.Second)

	json.NewEncoder(w).Encode(map[string]any{
		"statusCode": http.StatusOK,
		"data": types.Response{
			ShortenUrl: shortedUrl,
			ExpireTime: t.Format(time.RFC822),
		},
		"message": "Successful",
	})
}

func RedirectURL(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len("/redirect/"):]

	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	rdb, err := db.InitRedisClient()

	if err != nil {
		panic(err)
	}

	originalURL, err := rdb.Get(shortKey).Result()

	if err != nil {
		http.Error(w, "Error URL", http.StatusBadRequest)
		// panic(err)
	}

	// Redirect the user to the original URL
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
