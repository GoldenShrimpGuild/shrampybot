package space

import (
	"bytes"
	"encoding/json"
	"shrampybot/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/litui/helix/v3"
)

var (
	keyMastodonToTwitch  = "v3/mastodon_to_twitch"
	keyTwitchToMastodon  = "v3/twitch_to_mastodon"
	keyTwitchUsers       = "v3/twitch_users"
	keyTwitchEventMsgIds = "v3/twitch_event_msg_ids"
	keyStreamCache       = "v3/stream_cache"
	keyStreamMeta        = "v3/stream_meta"
	keyChannelCache      = "v3/channel_cache"
	keyCategoryMap       = "v3/category_map"
	keySubsCache         = "v3/subs_cache"
)

type Space struct {
	s3Session *session.Session
	s3Client  *s3.S3
	s3Bucket  string
}

func NewSpace() (Space, error) {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(config.AwsAccessKeyId, config.AwsSecretAccessKey, ""),
		Endpoint:         aws.String(config.AwsEndpointUrl),
		Region:           aws.String("us-east-1"), // nonsense region to handle AWS-specific logic
		S3ForcePathStyle: aws.Bool(false),         // Configures to use subdomain/virtual calling format. Depending on your version, alternatively use o.UsePathStyle = false
	}

	s3Session, s3Error := session.NewSession(s3Config)
	if s3Error != nil {
		return Space{}, s3Error
	}
	s3Client := s3.New(s3Session)

	return Space{
		s3Session: s3Session,
		s3Client:  s3Client,
		s3Bucket:  config.AwsBucket,
	}, nil
}

func (s *Space) GetTwitchUsers() (helix.ManyUsers, error) {
	obj, err := s.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: &s.s3Bucket,
		Key:    &keyTwitchUsers,
	})
	if err != nil {
		return helix.ManyUsers{}, err
	}

	buffer := []byte{}
	byteCount, err := obj.Body.Read(buffer)
	if err != nil || byteCount < 1 {
		return helix.ManyUsers{}, err
	}

	output := helix.ManyUsers{}
	err = json.Unmarshal(buffer, &output)
	if err != nil {
		return helix.ManyUsers{}, err
	}

	return output, nil
}

func (s *Space) PutTwitchUsers(users *helix.ManyUsers) error {
	usersJson, err := json.Marshal(users)
	if err != nil {
		return err
	}

	_, err = s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: &s.s3Bucket,
		Key:    &keyTwitchUsers,
		Body:   bytes.NewReader(usersJson),
	})
	if err != nil {
		return err
	}

	return nil
}
