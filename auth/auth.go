package auth

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/mattn/go-sqlite3"
	zaia_cache "github.com/youyo/zaia/cache"
	zaia_crypt "github.com/youyo/zaia/crypt"
)

func newSession() *session.Session {
	return session.Must(session.NewSession())
}

func newCredential(sess *session.Session, arn string) *credentials.Credentials {
	return stscreds.NewCredentials(sess, arn, func(p *stscreds.AssumeRoleProvider) {
		p.Duration = time.Duration(59) * time.Minute
	})
}

func getCredentialValues(creds *credentials.Credentials) (credentials.Value, error) {
	return creds.Get()
}

func newCredentialsFromCreds(credValues credentials.Value) *credentials.Credentials {
	return credentials.NewStaticCredentialsFromCreds(credValues)
}

func newConfig(creds *credentials.Credentials, region string) *aws.Config {
	return aws.NewConfig().WithRegion(region).WithCredentials(creds)
}

func getNewCredentialValues(sess *session.Session, arn string) (credValues credentials.Value, err error) {
	creds := newCredential(sess, arn)
	credValues, err = getCredentialValues(creds)
	if err != nil {
		return
	}
	encodedCredValues, err := zaia_crypt.Encode(credValues)
	if err != nil {
		return
	}
	err = zaia_cache.WriteCredentialsToCache(arn, encodedCredValues)
	return
}

func Auth(arn, region string) (*session.Session, *aws.Config) {
	// get credentials from sqlite, if record is exist.
	sess := newSession()
	credValues := func(sess *session.Session, arn string) credentials.Value {
		credValues, err := zaia_cache.ReadCredentialsFromCache(arn)
		if err != nil {
			credValues, err = getNewCredentialValues(sess, arn)
			if err != nil {
				log.Fatal(err)
			}
			return credValues
		}
		return credValues
	}(sess, arn)
	creds := newCredentialsFromCreds(credValues)
	config := newConfig(creds, region)
	return sess, config
}
