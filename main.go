package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

type Example struct {
	Sequence int    `json:"id"`
	Message  string `json:"message"`
}

func eventReceiver(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, error) {
	log.Printf("Handling EventID: %s", event.ID())
	fmt.Printf("Got Event Context: %+v\n", event.Context)
	data := &Example{}
	if err := event.DataAs(data); err != nil {
		fmt.Printf("Got Data Error: %s\n", err.Error())
	}
	fmt.Printf("Got Data: %+v\n", data)
	fmt.Printf("----------------------------\n")

	responseEvent := cloudevents.NewEvent()
	responseEvent.SetID(uuid.New().String())
	responseEvent.SetSource("/test")
	responseEvent.SetType("samples.http.test")
	responseEvent.SetSubject(fmt.Sprintf("%s#%d", event.Source(), data.Sequence))
	_ = responseEvent.SetData(cloudevents.ApplicationJSON, Example{
		Sequence: data.Sequence,
		Message:  "test done!",
	})
	log.Printf("Sending event with ID %s", responseEvent.ID())
	return &responseEvent, nil
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	log.Print("Starting CE-Test")
	ctx := context.Background()

	p, err := cloudevents.NewHTTP(cloudevents.WithPort(env.Port), cloudevents.WithPath(env.Path))
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}
	c, err := cloudevents.NewClient(p,
		cloudevents.WithUUIDs(),
		cloudevents.WithTimeNow(),
	)
	if err != nil {
		log.Fatalf("failed to create client: %s", err.Error())
	}

	log.Printf("listening on :%d%s\n", env.Port, env.Path)
	if err := c.StartReceiver(ctx, eventReceiver); err != nil {
		log.Fatalf("failed to start receiver: %s", err.Error())
	}

	<-ctx.Done()
}
