package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const migrationTemplate = `package migrations

import (
	"e-shop-api/internal/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func {{.Name}}Migration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "{{.Timestamp}}",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.{{.Name}}{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("{{.TableName}}")
		},
	}
}
`

const registryTemplate = `package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{{.Name}}Migration(),
	})

	return m.Migrate()
}
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/gen/main.go <ModelName>")
		return
	}

	modelName := os.Args[1]
	lowerName := strings.ToLower(modelName)
	
	var tableName string
	if strings.HasSuffix(lowerName, "y") {
		tableName = strings.TrimSuffix(lowerName, "y") + "ies"
	} else if strings.HasSuffix(lowerName, "s") {
		tableName = lowerName
	} else {
		tableName = lowerName + "s"
	}

	ts := time.Now().Format("20060102150405")
	dirPath := "internal/migrations"
	os.MkdirAll(dirPath, 0755)

	fileName := filepath.Join(dirPath, fmt.Sprintf("%s_%s_migration.go", ts, lowerName))

	f, _ := os.Create(fileName)
	t := template.Must(template.New("migration").Parse(migrationTemplate))
	t.Execute(f, map[string]string{
		"Name":      modelName,
		"Timestamp": ts,
		"TableName": tableName,
	})
	f.Close()

	handleRegistry(modelName)

	fmt.Printf("Success! Migration %s created and registered.\n", modelName)
}

func handleRegistry(modelName string) {
	registryPath := "internal/migrations/migration.go"
	newEntry := fmt.Sprintf("\t\t%sMigration(),", modelName)

	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		fmt.Println("Registry not found. Creating new migration.go...")
		f, _ := os.Create(registryPath)
		defer f.Close()

		t := template.Must(template.New("migration").Parse(registryTemplate))
		t.Execute(f, map[string]string{"Name": modelName})
		return
	}

	input, _ := os.ReadFile(registryPath)
	lines := strings.Split(string(input), "\n")
	var newLines []string

	for _, line := range lines {
		if strings.Contains(line, "})") {
			newLines = append(newLines, newEntry)
		}
		newLines = append(newLines, line)
	}

	os.WriteFile(registryPath, []byte(strings.Join(newLines, "\n")), 0644)
}