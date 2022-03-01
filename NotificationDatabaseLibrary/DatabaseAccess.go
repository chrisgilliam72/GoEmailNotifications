package NotificationDatabaseLibrary

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

const connectionString = "Server=localhost;user=COMCORPOFFICE\\Chris;password=June1972+;Database=EmailNotifications;Trusted_Connection=True;"

func openDBConnection() (*sql.DB, error) {
	sqlDB, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, fmt.Errorf(" error opening db connection: %v", err)

	}
	return sqlDB, nil
}

func closeDBConnection(sqlDB *sql.DB) error {
	err := sqlDB.Close()
	if err != nil {
		return fmt.Errorf(" error closing db connection: %v", err)
	}

	return nil
}

func GetEmailAddress(applicantionRef string) (emailAddress string, emailName string, error error) {

	sqlDB, err := openDBConnection()

	if err != nil {
		return "", "", err
	}
	defer closeDBConnection(sqlDB)

	ctx := context.Background()

	tsql := "select EmailAddress,Name from EmailAdressess where ApplicationReference = @ApplicationReference"
	rows, err := sqlDB.QueryContext(ctx, tsql, sql.Named("ApplicationReference", applicantionRef))
	if err != nil {
		return "", "", fmt.Errorf(" error queerying email address: %v", err)
	}
	rows.Next()
	rows.Scan(&emailAddress, &emailName)
	return emailAddress, emailName, nil

}

func AddNotificationMessage(notificationType int, applicationReference string, bankReference string, eventDate time.Time,
	eventType string, eventComment string, requestType string, messageStatus string, message string, eventID int) (int64, error) {

	sqlDB, err := openDBConnection()

	if err != nil {
		return -1, fmt.Errorf(" error opening db connection: %v", err)
	}
	defer closeDBConnection(sqlDB)

	ctx := context.Background()
	tsql := `insert into [NotificationsLog] ([ApplicationReference],[BankReference],[EventType],[DateTime] ,[EventComment],[RequestType],[MessageStatus],[Message],[EventId])
			 values ( @ApplicationReference,@BankReference,@EventType, @EventDate,@EventComment, @RequestType,@MessageStatus,@Message,@EventId) ;select isNull(SCOPE_IDENTITY(), -1) `

	stmt, err := sqlDB.Prepare(tsql)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("ApplicationReference", applicationReference),
		sql.Named("BankReference", bankReference),
		sql.Named("EventType", eventType),
		sql.Named("EventDate", eventDate),
		sql.Named("EventComment", eventComment),
		sql.Named("RequestType", requestType),
		sql.Named("MessageStatus", messageStatus),
		sql.Named("Message", message),
		sql.Named("EventId", eventID))
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return -1, err
	}

	return newID, nil
}
