package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"

	"github.com/anvari1313/river/datastore"
)

var streamCMD = &cobra.Command{
	Use:   "stream",
	Short: "Start publishing messages to RabbitMQ",
	Run: func(cmd *cobra.Command, args []string) {
		stream()
	},
}

var (
	dbURI      string
	db         string
	collection string
)

func init() {
	streamCMD.PersistentFlags().StringVar(&dbURI, "db-uri", "", "Data store URI")
	streamCMD.PersistentFlags().StringVar(&db, "db-name", "", "Data store Database")
	streamCMD.PersistentFlags().StringVar(&collection, "ds-collection", "", "Data store URI")
}

func stream() {
	ds, err := datastore.NewMongoDataStore(dbURI, db, collection)
	if err != nil {
		log.WithField("error", err).Fatal("error in data store initiation")
	}

	defer ds.Close()

	conn, err := amqp.Dial(rabbitURI)
	if err != nil {
		log.WithField("error", err).Fatal("error in dialing to rabbit mq")
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.WithField("error", err).Fatal("error in opening the channel")
	}

	defer ch.Close()

	err = ch.ExchangeDeclare(
		rabbitEx,
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.WithField("error", err).Fatal("error in declaring the exchange")
	}

	for ds.HasNext() {
		data, err := ds.Next()
		if err != nil {
			log.WithField("err", err).Error("error in fetching record from data store")
			continue
		}

		err = ch.Publish(
			rabbitEx, // exchange
			"",       // routing key
			false,    // mandatory
			false,    // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        data,
			})
		if err != nil {
			log.WithField("err", err).Error("error in fetching record from data store")
			continue
		}
	}
}
