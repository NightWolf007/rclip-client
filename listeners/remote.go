package listeners

import (
	"context"

	"github.com/NightWolf007/rclip/pb"
	"github.com/atotto/clipboard"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RunRemoteListener() error {
	conn, err := grpc.Dial(
		viper.GetString("address"),
		grpc.WithInsecure(),
		grpc.WithTimeout(viper.GetDuration("timeout")),
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

		log.Debug().Bytes("message", clip.Data).Msgf("New message from RClip")

		err = clipboard.WriteAll(string(clip.Data))
		if err != nil {
			return err
		}
	}
}
