package querytracker

import (
	"e-shop-api/internal/constants"
	"e-shop-api/internal/pkg/logger"
	"e-shop-api/internal/pkg/utils"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	slowQueryThreshold = utils.GetEnvTime("SLOW_QUERY_THRESHOLD", constants.SlowQueryThreshold)
)

type slowQueryPlugin struct{}

func (p *slowQueryPlugin) Name() string {
	return "slowquery"
}

func (p *slowQueryPlugin) Initialize(db *gorm.DB) error {
	if err := db.Callback().Query().Before("gorm:before_query").Register("slowquery:before", beforeQuery); err != nil {
		return err
	}
	if err := db.Callback().Query().After("gorm:after_query").Register("slowquery:after", afterQuery); err != nil {
		return err
	}

	if err := db.Callback().Create().Before("gorm:before_create").Register("slowquery:before_create", beforeQuery); err != nil {
		return err
	}
	if err := db.Callback().Create().After("gorm:after_create").Register("slowquery:after_create", afterQuery); err != nil {
		return err
	}

	if err := db.Callback().Update().Before("gorm:before_update").Register("slowquery:before_update", beforeQuery); err != nil {
		return err
	}
	if err := db.Callback().Update().After("gorm:after_update").Register("slowquery:after_update", afterQuery); err != nil {
		return err
	}

	if err := db.Callback().Delete().Before("gorm:before_delete").Register("slowquery:before_delete", beforeQuery); err != nil {
		return err
	}
	if err := db.Callback().Delete().After("gorm:after_delete").Register("slowquery:after_delete", afterQuery); err != nil {
		return err
	}

	return nil
}

func beforeQuery(db *gorm.DB) {
	db.Set("start_time", time.Now())
}

func afterQuery(db *gorm.DB) {
	startTime, ok := db.Get("start_time")
	if !ok {
		return
	}

	duration := time.Since(startTime.(time.Time))
	sql := db.Statement.SQL.String()

	logger.L.Info(
		fmt.Sprintf("Query executed: %v - %s", duration.Round(time.Millisecond), extractTableName(sql)),
		zap.String("sql", sql),
		zap.Duration("duration", duration),
	)

	if duration > slowQueryThreshold {
		logger.L.Warn(
			fmt.Sprintf("[SLOW QUERY] Query executed: %v - %s", duration.Round(time.Millisecond), extractTableName(sql)),
			zap.String("sql", sql),
			zap.Duration("duration", duration),
			zap.Duration("threshold", slowQueryThreshold),
		)
	}
}

func extractTableName(sql string) string {
	sql = strings.TrimSpace(sql)
	sqlUpper := strings.ToUpper(sql)

	keywords := []string{"FROM", "INTO", "UPDATE", "JOIN"}
	for _, keyword := range keywords {
		if idx := strings.Index(sqlUpper, keyword+" "); idx != -1 {
			tablePart := strings.TrimSpace(sql[idx+len(keyword)+1:])
			tablePart = regexp.MustCompile(`[\s;,].*`).ReplaceAllString(tablePart, "")
			tablePart = regexp.MustCompile(`\(.*`).ReplaceAllString(tablePart, "")
			if tablePart != "" {
				return tablePart
			}
		}
	}

	if strings.HasPrefix(sqlUpper, "SELECT") {
		return "subquery"
	}

	return "unknown"
}

func Register(db *gorm.DB) error {
	plugin := &slowQueryPlugin{}
	if err := plugin.Initialize(db); err != nil {
		return fmt.Errorf("failed to initialize slow query plugin: %w", err)
	}
	return nil
}
