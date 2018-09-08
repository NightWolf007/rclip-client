package listeners

import (
	"context"
	"time"

	"github.com/NightWolf007/rclip/pb"
	"github.com/atotto/clipboard"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func RunLocalListener(addr string, timeout time.Duration, updateDelay time.Duration) error {
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

	var prevData string
	for {
		data, err := clipboard.ReadAll()
		if err != nil {
			return err
		}

		if data != prevData {
			log.Debug().Str("data", data).Msg("New message from local clipboard")

			_, err = client.Push(context.Background(), &pb.PushRequest{Data: []byte(data)})
			if err != nil {
				return err
			}

			prevData = data
		}
		time.Sleep(updateDelay)
	}
}
