package v1

import (
	"github.com/icrowley/fake"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
)

type Env struct {
	Meta
	cobra.Command
}

type Meta struct {
	Instance                             string
	ServiceVersion                       string
	BuildNumber                          string
	KafkaRetryTopic                      string
	KafkaEingangsTopic                   string
	KafkaDeadLetterTopic                 string
	KafkaBrokers                         string
	KafkaUser                            string
	KafkaPassword                        string
	KafkaConsumerGroupID                 string
	KafkaRetrysHeaderName                string
	KafkaFailedConsumerGroupIdHeaderName string
	KafkaBatchSize                       int
	KafkaRetryLimit                      int
	KafkaAuthEnabled                     bool
}

func (env *Env) Cobra() {
	env.Command = cobra.Command{
		Use:   "RETRY",
		Short: "RETRY",
		Long:  "RETRY is a microservice for consuming and publishing failed kafka messages",
		Run: func(cmd *cobra.Command, args []string) {
			env.Meta = Meta{
				Instance:                             fake.City(),
				BuildNumber:                          viper.GetString("build_number"),
				KafkaRetryTopic:                      viper.GetString("kafka_retry_topic"),
				KafkaEingangsTopic:                   viper.GetString("kafka_eingangs_topic"),
				KafkaDeadLetterTopic:                 viper.GetString("kafka_deadletter_topic"),
				KafkaBrokers:                         viper.GetString("kafka_brokers"),
				KafkaUser:                            viper.GetString("kafka_user"),
				KafkaPassword:                        viper.GetString("kafka_password"),
				KafkaConsumerGroupID:                 viper.GetString("kafka_consumer_group_id"),
				KafkaRetrysHeaderName:                viper.GetString("kafka_retrys_header_name"),
				KafkaFailedConsumerGroupIdHeaderName: viper.GetString("kafka_failed_consumer_group_id_header_name"),
				KafkaBatchSize:                       viper.GetInt("kafka_batch_size"),
				KafkaRetryLimit:                      viper.GetInt("kafka_retry_limit"),
				KafkaAuthEnabled:                     viper.GetBool("kafka_auth_enabled"),
			}
		},
	}
	env.Init()
}

func (env *Env) Init() {

	viper.SetEnvPrefix("RETRY")
	viper.AutomaticEnv()

	flags := env.Flags()
	flags.String("build_number", "0000", "Build Number")
	flags.String("kafka_retry_topic", "", "Kafka Retry Topic")
	flags.String("kafka_eingangs_topic", "", "Kafka Eingangs Topic")
	flags.String("kafka_deadletter_topic", "", "Kafka DeadLetter Topic")
	flags.String("kafka_brokers", "", "Kafka Brokers")
	flags.String("kafka_user", "", "Kafka Username")
	flags.String("kafka_password", "", "Kafka Password")
	flags.String("kafka_consumer_group_id", "", "Kafka Consumer Group ID")
	flags.String("kafka_retrys_header_name", "retrys", "Kafka Retrys Header Name / Anzahl der Verarbeitungsversuche")
	flags.String("kafka_failed_consumer_group_id_header_name", "failed.consumer.group.id", "Failed Consumer Group ID / Consumer Group in der die Messageverarbeitung gefailt ist")
	flags.Int("kafka_batch_size", 5, "Kafka Batch Size")
	flags.Int("kafka_retry_limit", 5, "Kafka Retry Limit")
	flags.Bool("kafka_auth_enabled", false, "Kafka Batch Size")

	_ = viper.BindPFlag("build_number", flags.Lookup("build_number"))
	_ = viper.BindPFlag("kafka_retry_topic", flags.Lookup("kafka_retry_topic"))
	_ = viper.BindPFlag("kafka_eingangs_topic", flags.Lookup("kafka_eingangs_topic"))
	_ = viper.BindPFlag("kafka_deadletter_topic", flags.Lookup("kafka_deadletter_topic"))
	_ = viper.BindPFlag("kafka_brokers", flags.Lookup("kafka_brokers"))
	_ = viper.BindPFlag("kafka_user", flags.Lookup("kafka_user"))
	_ = viper.BindPFlag("kafka_password", flags.Lookup("kafka_password"))
	_ = viper.BindPFlag("kafka_consumer_group_id", flags.Lookup("kafka_consumer_group_id"))
	_ = viper.BindPFlag("kafka_retrys_header_name", flags.Lookup("kafka_retrys_header_name"))
	_ = viper.BindPFlag("kafka_failed_consumer_group_id_header_name", flags.Lookup("kafka_failed_consumer_group_id_header_name"))
	_ = viper.BindPFlag("kafka_batch_size", flags.Lookup("kafka_batch_size"))
	_ = viper.BindPFlag("kafka_retry_limit", flags.Lookup("kafka_retry_limit"))
	_ = viper.BindPFlag("kafka_auth_enabled", flags.Lookup("kafka_auth_enabled"))

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/go/bin")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Panic(err)
	}

	_ = env.Execute()
}

func (meta *Meta) Log() {
	logrus.Printf("-----------------------------------------------------------")
	logrus.Printf("Microservice Parameter")
	logrus.Printf("-----------------------------------------------------------")
	logrus.Printf("Instance => %s", meta.Instance)
	logrus.Printf("Version => %s", meta.ServiceVersion)
	logrus.Printf("BuildNumber => %s", meta.BuildNumber)
	logrus.Printf("-----------------------------------------------------------")
	logrus.Printf("Kafka Cluster Parameter")
	logrus.Printf("-----------------------------------------------------------")
	logrus.Printf("Brokers => %s", meta.KafkaBrokers)
	logrus.Printf("Retry Topic => %s", meta.KafkaRetryTopic)
	logrus.Printf("Eingangs Topic => %s", meta.KafkaEingangsTopic)
	logrus.Printf("DeadLetter Topic => %s", meta.KafkaDeadLetterTopic)
	logrus.Printf("Auth Enabled => %s", strconv.FormatBool(meta.KafkaAuthEnabled))
	logrus.Printf("Consumer Group ID => %s", meta.KafkaConsumerGroupID)
	logrus.Printf("Consumer Group ID => %s", meta.KafkaConsumerGroupID)
	logrus.Printf("Kafka Retrys Header Name => %s", meta.KafkaRetrysHeaderName)
	logrus.Printf("Kafka Failed Consumer Group ID Name => %s", meta.KafkaFailedConsumerGroupIdHeaderName)
	logrus.Printf("Batch Size => %s", strconv.Itoa(meta.KafkaBatchSize))
	logrus.Printf("Retry Limit => %s", strconv.Itoa(meta.KafkaRetryLimit))
	logrus.Printf("-----------------------------------------------------------")
	logrus.Printf("System Messages")
	logrus.Printf("-----------------------------------------------------------")
}
