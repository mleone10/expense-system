package dynamodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/mleone10/expense-system/domain"
)

func (c *Client) GetOrg(ctx context.Context, orgId domain.OrgId) (domain.Organization, error) {
	return domain.Organization{}, nil
}

// GetOrgsForUser first retrieves the IDs of all organizations that the given user is a member of.  Then it calls GetOrg for each ID to retrieve further organization details.
func (c *Client) GetOrgsForUser(ctx context.Context, userId domain.UserId) ([]domain.Organization, error) {
	res, err := c.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              c.table,
		IndexName:              c.reverseLookupGsi,
		KeyConditionExpression: aws.String("sk = :userId and begins_with(pk, :orgPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId":    &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", prefixUser, userId)},
			":orgPrefix": &types.AttributeValueMemberS{Value: prefixOrg},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orgs from dynamodb: %w", err)
	}

	orgUserRecords := []orgUserRecord{}
	err = attributevalue.UnmarshalListOfMaps(res.Items, &orgUserRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dynamodb response: %w", err)
	}

	orgs := []domain.Organization{}
	for _, orgUser := range orgUserRecords {
		orgId := domain.OrgId(strings.Split(orgUser.OrgId, "#")[1])

		org, err := c.GetOrg(ctx, orgId)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve org %v: %w", orgId, err)
		}

		orgs = append(orgs, org)
	}

	return orgs, nil
}
