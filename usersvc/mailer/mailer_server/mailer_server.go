package mailerserver

import "userService/usersvc/mailer/mailer_grpc"

type MailerServer struct {
	mailer_grpc.UnimplementedUserServer
}
