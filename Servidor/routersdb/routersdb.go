package routersdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
)

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

type ClientDB struct {
	ClientID      int     `json:"client_id"`
	NomeClient    string  `json:"nome_client"`
	EmpresaClient string  `json:"empresa_client"`
	Email         string  `json:"email"`
	PowerCap      float64 `json:"power_cap"`
}

func EmpresasInsert(empresa Empresa) bool {
	connStr := "admin:adminroot@tcp(db-hack.czim4vhg4n7p.us-east-1.rds.amazonaws.com:3306)/hacktoon"
	db, err := sql.Open("mysql", connStr)

	if err != nil {
		fmt.Println(err)
		return false
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Connected!")

	ctx := context.Background() // Create a context
	query := "INSERT INTO empresas(cnpj, nome_empresa, senha) VALUES (?,?,?);"

	cnpj := empresa.CNPJ
	nome_empresa := empresa.NomeEmpresa
	senha := empresa.Senha

	res, err := db.ExecContext(ctx, query, cnpj, nome_empresa, senha)

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println(res)

	lastId, err := res.LastInsertId()

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Printf("The last inserted row id: %d\n", lastId)
	return true
}

func ClienteInsert(cliente Cliente) bool {
	connStr := "admin:adminroot@tcp(db-hack.czim4vhg4n7p.us-east-1.rds.amazonaws.com:3306)/hacktoon"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return false
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Connected!")

	ctx := context.Background() // Create a context
	query := "INSERT INTO client(nome_client, empresa_client, email, power_cap) VALUES (?,?,?,?);"

	nome_client := cliente.Nome
	empresa_client := cliente.EmpresaCLiente
	email := cliente.Email
	power_cap := rand.Intn(100)

	res, err := db.ExecContext(ctx, query, nome_client, empresa_client, email, power_cap)

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println(res)

	lastId, err := res.LastInsertId()

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Printf("The last inserted row id: %d\n", lastId)
	return true
}

func ReturnEmpresaDado(cnpj string) string {
	connStr := "admin:adminroot@tcp(db-hack.czim4vhg4n7p.us-east-1.rds.amazonaws.com:3306)/hacktoon"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		fmt.Println(err)
	}

	cnpjInt, err := strconv.Atoi(cnpj)

	defer db.Close()

	if err != nil {
		fmt.Println(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connected!")

	ctx := context.Background() // Create a context
	query := "SELECT * FROM empresas WHERE cnpj = ?;"

	rows, err := db.QueryContext(ctx, query, cnpjInt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var empresas []Empresa
	for rows.Next() {
		var empresa Empresa
		err := rows.Scan(&empresa.CNPJ, &empresa.NomeEmpresa, &empresa.Senha, &empresa.Alarms)
		if err != nil {
			log.Fatal(err)
		}
		empresas = append(empresas, empresa)
	}

	jsonData, err := json.Marshal(empresas)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))
	return string(jsonData)
}

func ReturnSingleClientDado(empresaNome string) string {
	connStr := "admin:adminroot@tcp(db-hack.czim4vhg4n7p.us-east-1.rds.amazonaws.com:3306)/hacktoon"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Connected!")
	ctx := context.Background() // Create a context
	query := "SELECT * FROM client WHERE empresa_client = ?;"

	rows, err := db.QueryContext(ctx, query, empresaNome)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var clients []ClientDB
	rows.Next()

	var client ClientDB
	err = rows.Scan(&client.ClientID, &client.NomeClient, &client.EmpresaClient, &client.Email, &client.PowerCap)
	if err != nil {
		fmt.Println(err)
	}
	clients = append(clients, client)

	jsonData, err := json.Marshal(clients)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(jsonData))
	return string(jsonData)
}
