package http

import (
	"email-svc/src/application/consumers"
	"email-svc/src/application/services"
	"email-svc/src/infrastructure/configuration"
	"email-svc/src/infrastructure/http/controllers"
	"email-svc/src/infrastructure/rabbitmq"
	"email-svc/src/infrastructure/rabbitmq/queues"
	sendgridDriver "email-svc/src/infrastructure/sendgrid"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sendgrid/sendgrid-go"
	"net/http"
)

type Server struct {
	config          configuration.Config
	logger          configuration.Logger
	rabbitmqChannel *amqp.Channel
	emailConsumer   *consumers.EmailConsumer
	inboundEmail    *controllers.InboundEmailController
}

func NewServer() *Server {
	config := configuration.LoadConfig()
	loggerLevel := config.GetString("logger.level")
	logger := configuration.NewLogger(loggerLevel)

	sendgridKey := config.GetString("sendgrid.clientKey")
	sendgridClient := sendgrid.NewSendClient(sendgridKey)
	outboundEmail := sendgridDriver.NewSendgridOutboundEmail(sendgridClient)

	channel := rabbitmq.CreateConnection(config.GetString("rabbitmq.url"))
	sendEmail := services.NewSendEmail(outboundEmail)
	emailsToSendConsumer := consumers.NewEmailConsumer(logger, sendEmail)

	receivedEmailsTopic := config.GetString("rabbitmq.queues.emailsReceived")
	emailsReceivedPublisher := queues.NewPublisher(channel, receivedEmailsTopic)
	receiveEmail := services.NewReceiveEmail(emailsReceivedPublisher)
	inboundEmail := controllers.NewInboundEmailController(receiveEmail, logger)

	return &Server{
		config:          config,
		logger:          logger,
		rabbitmqChannel: channel,
		emailConsumer:   emailsToSendConsumer,
		inboundEmail:    inboundEmail,
	}
}

func (s Server) Run() {
	s.logger.Info("Starting up server")
	s.bindConsumers()
	s.bindRoutes()
}

func (s Server) bindConsumers() {
	emailsToSend := s.config.GetString("rabbitmq.queues.emailsToSend")
	queues.ConsumeQueue(s.rabbitmqChannel, emailsToSend, s.emailConsumer)
}

func (s Server) bindRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	})

	router.Handle("/emails", s.inboundEmail).Methods(http.MethodPost)

	s.logger.Fatal(http.ListenAndServe(s.config.GetString("http.port"), router))
}
