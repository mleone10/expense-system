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

// GetOrg retrieves all details of the organization with the given ID.  If no org exists with that ID, an error is returned.
func (c *Client) GetOrg(ctx context.Context, orgId domain.OrgId) (domain.Organization, error) {
	orgRes, err := c.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: c.table,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", prefixOrg, orgId)},
			"sk": &types.AttributeValueMemberS{Value: "ORG"},
		},
	})
	if err != nil {
		return domain.Organization{}, fmt.Errorf("failed to retrieve org %v: %w", orgId, err)
	}
	if orgRes.Item == nil {
		return domain.Organization{}, fmt.Errorf("no org found with id %v", orgId)
	}

	or := orgRecord{}
	attributevalue.UnmarshalMap(orgRes.Item, &or)

	usersRes, err := c.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              c.table,
		KeyConditionExpression: aws.String("pk = :orgId and begins_with(sk, :userPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":orgId":      &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", prefixOrg, orgId)},
			":userPrefix": &types.AttributeValueMemberS{Value: prefixUser},
		},
	})
	if err != nil {
		return domain.Organization{}, fmt.Errorf("failed to query dynamodb for members")
	}

	orgUserRecords := []orgUserRecord{}
	err = attributevalue.UnmarshalListOfMaps(usersRes.Items, &orgUserRecords)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("failed to parse dynamodb response: %w", err)
	}

	members := []domain.Member{}
	for _, orgUser := range orgUserRecords {
		userId := strings.Split(orgUser.UserIdKey, "#")[1]
		members = append(members, domain.Member{
			Id:    domain.UserId(userId),
			Admin: orgUser.Admin,
		})
	}

	return domain.Organization{
		Id:      domain.OrgId(or.OrgId),
		Name:    or.Name,
		Members: members,
	}, nil
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

func (c *Client) CreateOrg(ctx context.Context, name string, adminId domain.UserId) (domain.OrgId, error) {
	orgId, err := domain.NewOrgId()
	if err != nil {
		return "", fmt.Errorf("failed to generate new org id: %w", err)
	}

	orgItem, err := attributevalue.MarshalMap(orgRecord{
		OrgId:   fmt.Sprintf("%s%s", prefixOrg, orgId),
		OrgFlag: "ORG",
		Name:    name,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal new org record to dynamodb item: %w", err)
	}

	orgAdminItem, err := attributevalue.MarshalMap(orgUserRecord{
		OrgId:     fmt.Sprintf("%s%s", prefixOrg, orgId),
		UserIdKey: fmt.Sprintf("%s%s", prefixUser, adminId),
		Admin:     true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal new org admin record to dynamodb item: %w", err)
	}

	_, err = c.db.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName: c.table,
					Item:      orgItem,
				},
			},
			{
				Put: &types.Put{
					TableName: c.table,
					Item:      orgAdminItem,
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to save new org and admin records to dynamodb: %w", err)
	}

	// TODO: Return JSON from http server when errors occur
	return orgId, nil
}
