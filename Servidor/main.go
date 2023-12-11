package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hackthon/routersdb"
	"io"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type Categoria struct {
	ID   int64
	Nome string
}

type Empresa struct {
	CNPJ        int    `json:cnpj"`
	NomeEmpresa string `json:nome_empresa"`
	Senha       string `json:senha"`
	Alarms      string `json:alarms"`
}

type Cliente struct {
	Nome           string `json:nome_client"`
	EmpresaCLiente string `json:empresa_client"`
	Email          string `json:email"`
}

func Openconn() {
	connStr := "admin:adminroot@tcp(db-hack.czim4vhg4n7p.us-east-1.rds.amazonaws.com:3306)/hacktoon"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("asset/*.html"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Index.html", nil)
}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/cadastroEmpresa", cadastroEmpresa)
	http.HandleFunc("/cadastroCliente", cadastroCliente)
	http.HandleFunc("/dadosEmpresa", dadosEmpresa)
	http.HandleFunc("/dadosClient", dadosClient)
	// http.HandleFunc("/retornarEmpresa", dadosEmpresa)

	if err := http.ListenAndServe(":5050", nil); err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func cadastroEmpresa(w http.ResponseWriter, r *http.Request) {
	corpoRequest, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println("[cadastroEmpresa ] - ERRO 1")
		log.Fatal(err)
	}

	var empresaToCadastro Empresa

	dataErr := json.Unmarshal(corpoRequest, &empresaToCadastro)

	if dataErr != nil {
		fmt.Println("[espHandler] - ENTROU NO ERRO 2")
		http.Error(w, dataErr.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(empresaToCadastro)

	if routersdb.EmpresasInsert(routersdb.Empresa(empresaToCadastro)) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func cadastroCliente(w http.ResponseWriter, r *http.Request) {
	corpoRequest, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println("[cadastroEmpresa ] - ERRO 1")
		log.Fatal(err)
	}

	var clienteToCadastro Cliente

	dataErr := json.Unmarshal(corpoRequest, &clienteToCadastro)
	if dataErr != nil {
		fmt.Println("[espHandler] - ENTROU NO ERRO 2")
		http.Error(w, dataErr.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(clienteToCadastro)
	if routersdb.ClienteInsert(routersdb.Cliente(clienteToCadastro)) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func dadosEmpresa(w http.ResponseWriter, r *http.Request) {
	corpoRequest, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println("[dadosEmpresa ] - ERRO 1")
	}

	resposta := routersdb.ReturnEmpresaDado(string(corpoRequest))

	w.Write([]byte(resposta))
}

func dadosClient(w http.ResponseWriter, r *http.Request) {
	corpoRequest, err := io.ReadAll(r.Body)

	fmt.Println(string(corpoRequest))

	if err != nil {
		fmt.Println("[dadosClient ] - ERRO 1")
	}

	resposta := routersdb.ReturnSingleClientDado(string(corpoRequest))

	w.Write([]byte(resposta))
}
