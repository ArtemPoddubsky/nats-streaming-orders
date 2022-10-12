package subscriber

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"main/internal/config"
	"main/internal/database"
	"main/internal/inmemory"
	"main/internal/log"
	"main/internal/model"
	"time"
)

var (
	connection stan.Conn
)

// Subscriber holds all data needed to restore cache, connect and subscribe to nats-streaming.
type Subscriber struct {
	cfg         *config.Config
	memoryCache *inmemory.Cache
	connection  stan.Conn
	Db          *database.Postgres
}

// NewSubscriber returns new Subscriber instance.
func NewSubscriber(cfg *config.Config, memoryCache *inmemory.Cache) Subscriber {
	return Subscriber{
		cfg:         cfg,
		memoryCache: memoryCache,
		connection:  nil,
		Db:          database.NewDB(cfg),
	}
}

// RestoreCache gets all data from database and loads it to cache.
func (s Subscriber) RestoreCache() {
	log.Logger.Infoln("Restoring cache...")

	records := s.Db.GetRecords()

	s.memoryCache.Mutex.Lock()
	for i := range records {
		s.memoryCache.Storage[records[i].OrderUID] = *records[i]
	}
	s.memoryCache.Mutex.Unlock()
}

// Run connects and subscribes to nats-streaming.
func (s Subscriber) Run() {
	reconnect := make(chan bool)

	go s.connect(reconnect)
	reconnect <- true

	s.subscribe(reconnect)
}

func (s Subscriber) connect(reconnect chan bool) {
	defer close(reconnect)
	for {
		select {
		case <-reconnect:
			var err error
			connection, err = stan.Connect(s.cfg.Nats.ServerID, s.cfg.Nats.ClientID,
				stan.NatsURL(s.cfg.Nats.NatsURL),
				stan.Pings(1, 5),
				stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
					log.Logger.Fatalf("Connection lost, reason: %v", reason)
					connection.Close()
					reconnect <- true
				}))

			if err != nil {
				log.Logger.Fatalln("Connection error: ", err)
			}

			log.Logger.Infoln("Connected to NATS-Streaming")

			reconnect <- true
		}
	}
}

func (s Subscriber) subscribe(reconnect chan bool) {
	defer close(reconnect)
	for {
		select {
		case <-reconnect:
			_, err := connection.Subscribe("WB", func(m *stan.Msg) {
				if err := m.Ack(); err != nil {
					log.Logger.Fatalln(err)
					return
				}

				record := model.RecordModel{}
				err := json.Unmarshal(m.Data, &record)
				if err != nil {
					log.Logger.Errorln(err)
					return
				}

				valid, err := record.ValidateRecord()
				if !valid || err != nil {
					log.Logger.Errorln("Model is not valid")
					return
				}

				log.Logger.Traceln("Received a message with orderUID:", record.OrderUID)

				if err = s.Db.AddRecord(record); err != nil {
					log.Logger.Errorln(err)
					return
				}
				s.memoryCache.Insert(record)

			}, stan.DurableName("durability"),
				stan.AckWait(time.Minute),
				stan.SetManualAckMode())

			if err != nil {
				log.Logger.Errorln(err)
			}
		}
	}
}
