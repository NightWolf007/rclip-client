package listeners

import (
	"context"
	"time"

	"github.com/NightWolf007/rclip/pb"
	"github.com/atotto/clipboard"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func RunRemoteListener(addr string, timeout time.Duration) error {
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithTimeout(timeout),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewClipboardClient(conn)
	stream, err := client.Subscribe(context.Background(), &pb.SubscribeRequest{})
	if err != nil {
		return err
	}

	for {
		clip, err := stream.Recv()
		if err != nil {
			return err
		}

		log.Debug().Bytes("data", clip.Data).Msg("New message from RClip")

		err = clipboard.WriteAll(string(clip.Data))
		if err != nil {
			return err
		}
	}
}
