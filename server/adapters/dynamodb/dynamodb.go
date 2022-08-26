package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/mleone10/expense-system/domain"
)

type Client struct {
	db               *dynamodb.Client
	table            *string
	reverseLookupGsi *string
}

func NewClient() (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, fmt.Errorf("failed to create session config: %w", err)
	}

	c := &Client{
		db:               dynamodb.NewFromConfig(cfg),
		table:            aws.String("expense-system-records"),
		reverseLookupGsi: aws.String("reverse-lookup"),
	}

	return c, nil
}

func (c *Client) GetOrg(orgId domain.OrgId) (domain.Organization, error) {
	return domain.Organization{}, nil
}

func (c *Client) GetOrgsForUser(userId domain.UserId) ([]domain.Organization, error) {
	return []domain.Organization{}, nil
}
