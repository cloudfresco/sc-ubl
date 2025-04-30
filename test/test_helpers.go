package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"strings"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

// GetUUIDDateValues -- data for tests
func GetUUIDDateValues(log *zap.Logger) ([]byte, string, *timestamp.Timestamp, *timestamp.Timestamp, string, string, string, string, error) {
	uuid4, err := common.GetUUIDBytes()
	if err != nil {
		log.Error("Error", zap.Error(err))
		return nil, "", nil, nil, "", "", "", "", err
	}

	ttime := common.GetTimeDetails()
	tn := common.TimeToTimestamp(ttime)

	uuid4Str, err := common.UUIDBytesToStr(uuid4)
	if err != nil {
		log.Error("Error", zap.Error(err))
		return nil, "", nil, nil, "", "", "", "", err
	}
	uuid4StrJSON, _ := json.Marshal(uuid4Str)
	uuid4JSON, _ := json.Marshal(uuid4)
	createdAtJSON, _ := json.Marshal(tn)
	updatedAtJSON, _ := json.Marshal(tn)
	return uuid4, uuid4Str, tn, tn, string(uuid4JSON), string(uuid4StrJSON), string(createdAtJSON), string(updatedAtJSON), nil
}

// LoadSQL -- drop db, create db, use db, load data
func LoadSQL(log *zap.Logger, dbService *common.DBService) error {
	var err error
	ctx := context.Background()

	if dbService.DBType == common.DBMysql {
		err = execSQLFile(ctx, log, dbService.MySQLTruncateFilePath, dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/address.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/credit_note_headers.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/credit_note_lines.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/debit_note_headers.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/debit_note_lines.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/despatch_headers.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/despatch_lines.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/invoice_headers.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/invoice_lines.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/location.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/party_contacts.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/party.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/purchase_order_headers.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/purchase_order_lines.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/receipt_advice_headers.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/receipt_advice_lines.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/tax_scheme.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/tax_category.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/tax_subtotal.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.MySQLTestFilePath+"/tax_total.sql", dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

	} else if dbService.DBType == common.DBPgsql {
		err = execSQLFile(ctx, log, dbService.PgSQLTruncateFilePath, dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}

		err = execSQLFile(ctx, log, dbService.PgSQLTestFilePath, dbService.DB)
		if err != nil {
			log.Error("Error", zap.Error(err))
			return err
		}
	}

	return nil
}

func execSQLFile(ctx context.Context, log *zap.Logger, sqlFilePath string, db *sqlx.DB) error {
	content, err := os.ReadFile(sqlFilePath)
	if err != nil {
		log.Error("Error", zap.Error(err))
		return err
	}
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error("Error", zap.Error(err))
	}

	sqlLines := strings.Split(string(content), ";\n")
	for _, sqlLine := range sqlLines {
		if sqlLine != "" {
			_, err := tx.ExecContext(ctx, sqlLine)
			if err != nil {
				log.Error("Error", zap.Error(err))
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Error("Error", zap.Error(rollbackErr))
					return err
				}
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Error("Error", zap.Error(err))
		return err
	}
	return nil
}
