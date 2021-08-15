package main

import (
	"context"
	"crypto"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	LetsEncryptStagingCA    = "https://acme-staging-v02.api.letsencrypt.org/directory"
	LetsEncryptProductionCA = "https://acme-v02.api.letsencrypt.org/directory"
	ZeroSSLProductionCA     = "https://acme.zerossl.com/v2/DV90"
)

type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

type Domain struct {
	Domain string
}

func (u *Domain) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// UnmarshalBinary decodes the struct into a User
func (u *Domain) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, u); err != nil {
		return err
	}
	return nil
}

func (u *Domain) String() string {
	return "Domain: " + u.Domain
}

func main() {

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalln("Error loading .env")
	}

	eggoConfig, err := buildConfig()
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: eggoConfig.RedisAddress,
	})

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&CertificateStore{})

	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		time.Sleep(3 * time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx := context.Background()
	topic := redisClient.Subscribe(ctx, "cert_issue")
	channel := topic.Channel()

	privateKey, _ := decodeKey(eggoConfig.PrivateKey, eggoConfig.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	myUser := MyUser{
		Email: eggoConfig.AcmeEmail,
		key:   privateKey,
	}

	config := lego.NewConfig(&myUser)

	config.CADirURL = ZeroSSLProductionCA
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Optionally can bind to a non-80 port. If binding to an alternative port number
	// will need a reverse-proxy to forward to the correct port.
	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", eggoConfig.ListenPort))
	if err != nil {
		log.Fatal(err)
	}

	for msg := range channel {
		d := &Domain{}
		err := d.UnmarshalBinary([]byte(msg.Payload))
		if err != nil {
			log.Fatal(err)
		}

		reg, err := handleZSSLAuth(client, eggoConfig.APIKey)
		if err != nil {
			log.Fatal(err)
		}
		myUser.Registration = reg

		request := certificate.ObtainRequest{
			Domains: []string{d.Domain},
			Bundle:  true,
		}
		certificates, err := client.Certificate.Obtain(request)
		if err != nil {
			log.Fatal(err)
		}

		c := CertificateStoreFromResource(certificates)
		db.Create(&c)
	}
}
