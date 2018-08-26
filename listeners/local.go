package listeners

import (
	"context"
	"time"

	"github.com/NightWolf007/rclip/pb"
	"github.com/atotto/clipboard"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RunLocalListener() error {
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

	var prevData string
	for {
		data, err := clipboard.ReadAll()
		if err != nil {
			return err
		}

		if data != prevData {
			log.Debug().Str("message", data).Msgf("New message from local clipboard")

			_, err = client.Push(context.Background(), &pb.PushRequest{Data: []byte(data)})
			if err != nil {
				return err
			}

			prevData = data
		}
		time.Sleep(viper.GetDuration("update_period"))
	}
}
