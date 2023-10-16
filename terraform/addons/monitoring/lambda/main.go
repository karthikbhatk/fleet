/*
This script is intended to be used with AWS Lambda to monitor the various
crons that live inside of Fleet.

We will check to see if there are recent updates from the crons in the
following table:

    - cron_stats

If we have an old/incomplete run in cron_stats or if we are missing a
cron entry entirely, throw an alert to an SNS topic.

Currently tested crons:

    - cleanups_then_aggregation
    - vulnerabilities

*/

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/go-sql-driver/mysql"
	flags "github.com/jessevdk/go-flags"
)

type NullEvent struct{}

type OptionsStruct struct {
	LambdaExecutionEnv string `long:"lambda-execution-environment" env:"AWS_EXECUTION_ENV"`
	SNSTopicArn        string `long:"sns-topic-arn" env:"SNS_TOPIC_ARN" required:"true"`
	MySQLHost          string `long:"mysql-host" env:"MYSQL_HOST" required:"true"`
	MySQLUser          string `long:"mysql-user" env:"MYSQL_USER" required:"true"`
	MySQLPassword      string `long:"mysql-password" env:"MYSQL_PASSWORD" required:"true"`
	MySQLDatabase      string `long:"mysql-database" env:"MYSQL_DATABASE" required:"true"`
	FleetEnv           string `long:"fleet-environment" env:"FLEET_ENV" required:"true"`
}

var options = OptionsStruct{}

func sendSNSMessage(msg string) {
	log.Printf("Sending SNS Message")
	region := "us-east-2"
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config: aws.Config{
				Region: &region,
			},
		},
	))
	svc := sns.New(sess)
	result, err := svc.Publish(&sns.PublishInput{
		Message:  &msg,
		TopicArn: &options.SNSTopicArn,
	})
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf(result.GoString())
}

func checkDB() (err error) {
	cfg := mysql.Config{
		User:                 options.MySQLUser,
		Passwd:               options.MySQLPassword,
		Net:                  "tcp",
		Addr:                 options.MySQLHost,
		DBName:               options.MySQLDatabase,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Printf(err.Error())
		sendSNSMessage("Unable to connect to database. Cron status unknown.")
		return err
	}
	if err = db.Ping(); err != nil {
		log.Printf(err.Error())
		sendSNSMessage("Unable to connect to database. Cron status unknown.")
		return err
	}

	log.Printf("Connected to database!")

	type CronStatsRow struct {
		name       string
		status     string
		updated_at time.Time
	}

	rows, err := db.Query("SELECT b.name,IFNULL(status, 'missing cron'),IFNULL(updated_at, FROM_UNIXTIME(0)) AS updated_at FROM (SELECT 'vulnerabilities' AS name UNION ALL SELECT 'cleanups_then_aggregation' UNION ALL SELECT 'missing') b LEFT JOIN (SELECT name, status, updated_at FROM cron_stats WHERE id IN (SELECT MAX(id) FROM cron_stats WHERE status = 'completed' GROUP BY name)) a ON a.name = b.name;")
	if err != nil {
		log.Printf(err.Error())
		sendSNSMessage("Unable to SELECT cron_stats table.  Unable to continue.")
		return err
	}
	twoHoursAgo := time.Now().Add(time.Duration(-2) * time.Hour)
	for rows.Next() {
		var row CronStatsRow
		if err := rows.Scan(&row.name, &row.status, &row.updated_at); err != nil {
			log.Printf(err.Error())
			sendSNSMessage("Error scanning row in cron_stats table.  Unable to continue.")
			return err
		}
		log.Printf("Row %s last updated at %s", row.name, row.updated_at.String())
		if row.updated_at.Before(twoHoursAgo) {
			log.Printf("*** %s hasn't updated in more than 2 hours, alerting! (status %s)", row.name, row.status)
			// Fire on the first match and return.  We only need to alert that the crons need looked at, not each cron.
			sendSNSMessage(fmt.Sprintf("In the environment '%s', Fleet cron '%s' hasn't updated in more than 2 hours. Last status was '%s' at %s.", options.FleetEnv, row.name, row.status, row.updated_at.String()))
			return nil
		}
	}

	// select name, status, updated_at from cron_stats where id in (select max(id) from cron_stats group by name);

	return nil
}

func handler(ctx context.Context, name NullEvent) error {
	checkDB()
	return nil
}

func main() {
	var err error
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Get config from environment
	parser := flags.NewParser(&options, flags.Default)
	if _, err = parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			return
		} else {
			log.Fatal(err)
		}
	}

	if options.LambdaExecutionEnv == "AWS_Lambda_go1.x" {
		lambda.Start(handler)
	} else {
		if err = handler(context.Background(), NullEvent{}); err != nil {
			log.Fatal(err)
		}
	}
}
