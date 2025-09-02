package auth

import (
	"context"
	"fmt"
)

/*
type IAuthentication interface {
	Register(login, password string) (token string, err error)
	Login(login, password string) (token string, err error)
}
*/

type AuthenticationProxy struct {
	remoteUnavailable bool
	localAuth         IAuthentication
	remoteAuth        IAuthentication
	ctx               context.Context
}

func (auth AuthenticationProxy) Register(login, password string) (token string, err error) {
	if !auth.remoteUnavailable {
		token, err = auth.remoteAuth.Register(login, password)
	}
	if l.remoteUnavailable || err != nil {
		l.localLogger.WriteWarning(fmt.Sprintf("%s: %v", "WriteWarning(): remote logger unavailable", err))
		l.localLogger.WriteWarning(msg)
	}

	return
}

func (auth AuthenticationProxy) Login(login, password string) (token string, err error) {

	return
}

/*
func qwwerty(url string) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()

	client := gen.NewReaderClient(conn)

	ctx, c := context.WithCancel(context.Background())
	defer c()

	client.GetAnnouncements()

	stream, err := client.Subscribe(ctx, &subpub_service.SubscribeRequest{Key: "Огонь"})
	if err != nil {
		fmt.Println(err)
	}

	msg, err := stream.Recv()
	if err == io.EOF {
		stream.CloseSend()
		fmt.Println(err)
		break
	}
	if err != nil {
		fmt.Println(err)
		break
	}
}
*/
