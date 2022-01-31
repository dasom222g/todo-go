package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

const cookieName = "OAUTH"
const googleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

var googleOauthConfig = oauth2.Config{
	// 로그인 완료 후 콜백받을 주소
	RedirectURL:  "http://localhost:3000/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// 구글로그인 사이트로 리다이렉트
	state := generateStateOauthCookie(w)
	url := googleOauthConfig.AuthCodeURL(state) // CSRF공격을 막기위한 인자
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleLoginCallback(w http.ResponseWriter, r *http.Request) {
	// 로그인 완료 후 state 확인, userInfo 가져오기
	cookie, _ := r.Cookie(cookieName)
	state := r.FormValue("state")
	if state != cookie.Value {
		// CSRF 공격시 에러 로그 남기고 페이지에도 에러표시
		errMessage := fmt.Sprintf("invalid goole oauth state. cookie: %s, state: %s\n", cookie.Value, state)
		log.Println(errMessage)
		http.Error(w, errMessage, http.StatusInternalServerError)
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	userInfo, err := getUserInfo(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	// id를 session에 저장 후 index페이지로 redirect
	var info googleInfo
	err = json.Unmarshal(userInfo, &info)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session, err := store.Get(r, "session")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// set session value
	session.Values["id"] = info.ID
	// save session
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	// 쿠키 만료시간을 현 시간부터 24시간 후로 설정
	expiration := time.Now().Add(1 * 24 * time.Hour)
	byteSlice := make([]byte, 16)
	rand.Read(byteSlice) // 랜덤한 수를 byteSlice에 채워춤
	// 메모리를 덜 차치하며 string으로 변환
	state := base64.URLEncoding.EncodeToString(byteSlice)

	cookie := &http.Cookie{
		Name:    cookieName,
		Value:   state,
		Expires: expiration,
	}
	http.SetCookie(w, cookie)
	return state
}

func getUserInfo(code string) ([]byte, error) {
	// 코드로 토큰 발행
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to exchange %s\n", err.Error())
	}
	// 토큰으로 googleAPI로 요청하여 사용자 정보 get
	res, err := http.Get(googleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to get userInfo %s\n", err.Error())
	}
	return ioutil.ReadAll(res.Body)
}
