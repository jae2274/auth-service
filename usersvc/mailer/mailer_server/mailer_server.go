package mailerserver

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"userService/usersvc/mailer/mailer_grpc"
	"userService/usersvc/mailer/mapper"
)

type MailerServer struct {
	db *sql.DB
	mailer_grpc.UnimplementedUserServer
}

func NewMailerServer(db *sql.DB) *MailerServer {
	return &MailerServer{
		db: db,
	}
}

func (s *MailerServer) GetUserEmails(ctx context.Context, req *mailer_grpc.GetUserEmailsRequest) (*mailer_grpc.UserEmails, error) {
	intUserIds := make([]int, len(req.UserIds))
	for i, userId := range req.UserIds {
		intUserId, err := strconv.Atoi(userId)
		if err != nil {
			return nil, err
		}
		intUserIds[i] = intUserId
	}

	users, err := mapper.GetUserEMails(ctx, s.db, intUserIds)
	if err != nil {
		return nil, err
	}

	userEmails := make([]*mailer_grpc.UserEmail, len(users))
	for i, user := range users {
		userEmails[i] = &mailer_grpc.UserEmail{
			UserId: fmt.Sprintf("%d", user.UserID),
			Email:  user.Email,
		}
	}

	return &mailer_grpc.UserEmails{
		UserEmails: userEmails,
	}, nil
}
