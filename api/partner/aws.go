package partner

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

var s3Session *s3.S3

const (
	BUCKET_NAME = "nft-image"
	REGION      = "kr-standard"
)

func init() {
	s3Session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(REGION),
		Endpoint:    aws.String("https://kr.object.ncloudstorage.com"),
		Credentials: credentials.NewSharedCredentials("./configs/partner/ncp/credentials", "default"),
	})))

}

func UploadObject(chunk []byte) (string, error) {
	fileId := uuid.NewString()
	buffer := bytes.NewReader(chunk)
	_, err := s3Session.PutObject(&s3.PutObjectInput{
		Body:        buffer,
		Bucket:      aws.String(BUCKET_NAME),
		Key:         aws.String(fileId),
		ACL:         aws.String(s3.BucketCannedACLPublicRead),
		ContentType: aws.String(http.DetectContentType(chunk)),
	})

	if err != nil {
		return "", errors.Wrap(err, "Upload Object")
	}

	return fileId, nil
}
