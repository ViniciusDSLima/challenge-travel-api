package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const migrationsDir = "./migrations"
const migrationTable = "schema_migrations"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("Error to connect with your database. Please set the DATABASE_DSN environment variable")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Erro ao verificar conexão com o banco de dados: %v", err)
	}

	createMigrationTable(db)

	switch command {
	case "up":
		migrateUp(db)
	case "down":
		steps := 1
		if len(os.Args) > 2 {
			var err error
			steps, err = strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalf("Número de steps inválido: %v", err)
			}
		}
		migrateDown(db, steps)
	case "status":
		showStatus(db)
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Você deve especificar uma versão para forçar")
		}
		forceVersion(db, os.Args[2])
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Uso: go run cmd/migrate/migrate.go <comando> [opções]")
	fmt.Println("Comandos:")
	fmt.Println("  up                  Executa todas as migrations pendentes")
	fmt.Println("  down [steps]        Reverte a última migration ou N migrations se steps for especificado")
	fmt.Println("  status              Mostra o status das migrations")
	fmt.Println("  force <version>     Define manualmente a versão atual")
}

func createMigrationTable(db *sql.DB) {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);
	`, migrationTable))

	if err != nil {
		log.Fatalf("Erro ao criar tabela de migrations: %v", err)
	}
}

func getMigrationFiles() (map[string]string, []string) {
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Erro ao ler diretório de migrations: %v", err)
	}

	migrationMap := make(map[string]string)
	var versions []string

	re := regexp.MustCompile(`^(\d+)_(.+)\.(up|down)\.sql$`)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matches := re.FindStringSubmatch(file.Name())
		if matches == nil {
			continue
		}

		version := matches[1]
		direction := matches[3]

		if direction == "up" {
			migrationMap[version] = file.Name()
			versions = append(versions, version)
		}
	}

	sort.Strings(versions)

	return migrationMap, versions
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT version FROM %s ORDER BY version", migrationTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, nil
}

func migrateUp(db *sql.DB) {
	migrationMap, versions := getMigrationFiles()
	applied, err := getAppliedMigrations(db)
	if err != nil {
		log.Fatalf("Erro ao obter migrations aplicadas: %v", err)
	}

	for _, version := range versions {
		if !applied[version] {
			upFile := filepath.Join(migrationsDir, migrationMap[version])
			fmt.Printf("Aplicando migration: %s\n", upFile)

			content, err := ioutil.ReadFile(upFile)
			if err != nil {
				log.Fatalf("Erro ao ler arquivo de migration: %v", err)
			}

			tx, err := db.Begin()
			if err != nil {
				log.Fatalf("Erro ao iniciar transação: %v", err)
			}

			_, err = tx.Exec(string(content))
			if err != nil {
				tx.Rollback()
				log.Fatalf("Erro ao executar migration %s: %v", upFile, err)
			}

			_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (version) VALUES ($1)", migrationTable), version)
			if err != nil {
				tx.Rollback()
				log.Fatalf("Erro ao registrar migration %s: %v", version, err)
			}

			if err := tx.Commit(); err != nil {
				log.Fatalf("Erro ao finalizar transação: %v", err)
			}

			fmt.Printf("Migration aplicada com sucesso: %s\n", version)
		}
	}

	fmt.Println("Todas as migrations foram aplicadas!")
}

func migrateDown(db *sql.DB, steps int) {
	if steps <= 0 {
		log.Fatal("O número de steps deve ser maior que zero")
	}

	rows, err := db.Query(fmt.Sprintf("SELECT version FROM %s ORDER BY version DESC LIMIT $1", migrationTable), steps)
	if err != nil {
		log.Fatalf("Erro ao obter migrations aplicadas: %v", err)
	}
	defer rows.Close()

	var versionsToRevert []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			log.Fatalf("Erro ao ler versão: %v", err)
		}
		versionsToRevert = append(versionsToRevert, version)
	}

	if len(versionsToRevert) == 0 {
		fmt.Println("Não há migrations para reverter")
		return
	}

	for _, version := range versionsToRevert {
		downFile := filepath.Join(migrationsDir, fmt.Sprintf("%s_", version))

		matches, err := filepath.Glob(downFile + "*.down.sql")
		if err != nil || len(matches) == 0 {
			log.Fatalf("Arquivo down.sql não encontrado para versão %s", version)
		}

		downFile = matches[0]
		fmt.Printf("Revertendo migration: %s\n", downFile)

		content, err := ioutil.ReadFile(downFile)
		if err != nil {
			log.Fatalf("Erro ao ler arquivo de migration: %v", err)
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("Erro ao iniciar transação: %v", err)
		}

		_, err = tx.Exec(string(content))
		if err != nil {
			tx.Rollback()
			log.Fatalf("Erro ao executar migration %s: %v", downFile, err)
		}

		_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE version = $1", migrationTable), version)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Erro ao remover registro da migration %s: %v", version, err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatalf("Erro ao finalizar transação: %v", err)
		}

		fmt.Printf("Migration revertida com sucesso: %s\n", version)
	}
}

func showStatus(db *sql.DB) {
	migrationMap, versions := getMigrationFiles()
	applied, err := getAppliedMigrations(db)
	if err != nil {
		log.Fatalf("Erro ao obter migrations aplicadas: %v", err)
	}

	fmt.Println("Status das migrations:")
	fmt.Println("======================")

	if len(versions) == 0 {
		fmt.Println("Nenhuma migration encontrada")
		return
	}

	for _, version := range versions {
		status := "Pendente"
		if applied[version] {
			status = "Aplicada"
		}
		name := strings.TrimSuffix(strings.TrimPrefix(migrationMap[version], version+"_"), ".up.sql")
		fmt.Printf("%s: %s (%s)\n", version, name, status)
	}
}

func forceVersion(db *sql.DB, version string) {
	_, versions := getMigrationFiles()
	versionExists := false
	for _, v := range versions {
		if v == version {
			versionExists = true
			break
		}
	}

	if !versionExists {
		log.Fatalf("Versão %s não encontrada nos arquivos de migration", version)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Erro ao iniciar transação: %v", err)
	}

	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s", migrationTable))
	if err != nil {
		tx.Rollback()
		log.Fatalf("Erro ao limpar tabela de migrations: %v", err)
	}

	_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (version, applied_at) VALUES ($1, $2)", migrationTable),
		version, time.Now())
	if err != nil {
		tx.Rollback()
		log.Fatalf("Erro ao inserir versão forçada: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Erro ao finalizar transação: %v", err)
	}

	fmt.Printf("Versão forçada para: %s\n", version)
}
