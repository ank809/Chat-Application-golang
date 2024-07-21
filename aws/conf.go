package aws_conf

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetAwsConf() (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {

		return aws.Config{}, err
	}
	return cfg, nil

}

func GetS3Client() (*s3.Client, error) {
	cfg, err := GetAwsConf()
	if err != nil {
		return nil, err
	}
	s3client := s3.NewFromConfig(cfg)
	return s3client, nil
}
