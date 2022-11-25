package apiKey

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"restapi/internal/apiKey"
	"restapi/internal/config"
	"restapi/pkg/client/postgresql"
	"restapi/pkg/logging"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
	cfg    *config.Config
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) FindeApiKey(ctx context.Context, key string) (apiKey.ApiKeyResult, error) {
	var (
		points                                int
		active, is_paid, is_spent, is_blocked uint
		accountName                           string
	)
	rez := apiKey.ApiKeyResult{}
	selectSQL :=
		"SELECT " +
			"billing_periods.points,api_keys.active,billing_periods.is_paid,billing_periods.is_spent,billing_periods.is_blocked, accounts.name" +
			" FROM " +
			"%s.api_keys" +
			" LEFT JOIN" +
			" %s.billing_periods ON api_keys.user_id = billing_periods.user_id " +
			" LEFT JOIN" +
			" bitquery.accounts ON accounts.id = api_keys.user_id WHERE api_keys.`key` = ? " +
			" AND NOW() BETWEEN started_at and ended_at "

	dbName := r.cfg.DBConf.DBName
	selectSQL = fmt.Sprintf(selectSQL, dbName, dbName)
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(selectSQL)))
	rows := r.client.QueryRow(ctx, selectSQL, key)

	err := rows.Scan(&points, &active, &is_paid, &is_spent, &is_blocked, &accountName)

	if err != nil {
		if err == sql.ErrNoRows {
			rez.ResponseBody = "Api key provided is not active, check your API Key in the profile page"
			return rez, nil
		} else {
			// we have error in db , so default response
			rez.Location = r.cfg.Classes[len(r.cfg.Classes)-1].Location
			return rez, err
		}
	}
	if active == 0 || is_spent == 1 || is_blocked == 1 {
		rez.ResponseBody = "Api key provided is not active, check your API Key in the profile page"
		return rez, nil
	}
	rez.Location, err = determineQOS(r.cfg, points, is_paid)
	rez.AccountName = accountName
	return rez, err
}

func determineQOS(conf *config.Config, points int, is_paid uint) (string, error) {
	//  we try to find the first match element from the array
	for i := 0; i < len(conf.Classes); i++ {
		if is_paid == conf.Classes[i].Paid && conf.Classes[i].Max_points <= points {
			return conf.Classes[i].Location, nil
		}
		// if there is no match we take the last default element
		if i == len(conf.Classes)-1 {
			return conf.Classes[len(conf.Classes)-1].Location, nil
		}
	}
	return "", errors.New("Something is wrong with config.")
}
func (r *repository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	return nil
}

func NewRepository(client postgresql.Client, cfg *config.Config, logger *logging.Logger) apiKey.Repository {
	return &repository{
		client: client,
		logger: logger,
		cfg:    cfg,
	}
}
