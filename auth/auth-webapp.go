package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "audi/account/proto"
	_ "audi/pkg/logger"
)

var (
	accountServiceClient pb.AccountServiceClient
	port                 string
)

func init() {
	flag.StringVar(&port, "port", "80", "The port of service listen")
	flag.Parse()

	addr, ok := os.LookupEnv("ACCOUNT_SERVICE_ADDR")
	if !ok {
		log.Fatal("ACCOUNT_SERVICE_ADDR not found in enviroment ")
		return
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatal("Failed to connect AccountService server")
	}

	accountServiceClient = pb.NewAccountServiceClient(conn)
	log.WithField("account.service.addr", addr).Info("account grpc service address")
}

func main() {
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.WithField("addr", addr).Info("Server starting")
	log.Fatal(http.ListenAndServe(addr, NewRouter()))
}

func NewRouter() *mux.Router {
	h := NewHandler()
	r := mux.NewRouter()
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/logout", h.Logout).Methods("POST")
	return r
}

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Params(r *http.Request) map[string]string {
	params := make(map[string]string)
	values := r.URL.Query()
	for k, _ := range values {
		if v := values.Get(k); v != "" {
			params[k] = strings.TrimSpace(v)
		}
	}
	return params
}

func (h *Handler) Error(w http.ResponseWriter, body error, statusCode int) {
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		fmt.Fprintln(w, err)
	}
}

func (h *Handler) WriteJson(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-type", "application/json;charset=UTF-8")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.WithError(err).Error("Failed to encode response")
		log.Errorf("%#v \n", err)
		h.Error(w, err, 500)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req pb.LoginReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Error(w, err, 400)
		log.WithError(err).Error("Failed to decode request body")
		log.Errorf("%#v \n", err)
		return
	}

	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		req.Ip = host
	}
	req.Device = r.UserAgent()
	req.Fp = r.Header.Get("Finger-Print")

	resp, err := accountServiceClient.Login(r.Context(), &req)

	if err != nil {
		h.Error(w, err, 401)
	}

	h.WriteJson(w, resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req pb.LogoutReq
	resp, err := accountServiceClient.Logout(r.Context(), &req)
	if err != nil {
		h.Error(w, err, 400)
	}
	h.WriteJson(w, resp)
}
