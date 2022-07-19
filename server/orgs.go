package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gofrs/uuid"
)

type role string

const (
	RoleAdmin    role = "ADMIN"
	RoleApprover role = "APPROVER"
	RoleUser     role = "USER"
)

type orgRepo struct {
	db               *dynamodb.Client
	table            *string
	reverseLookupGsi *string
}

type orgRecord struct {
	OrgId         string `dynamodbav:"pk"`
	OrgPrimaryKey string `dynamodbav:"sk"`
	Name          string `dynamodbav:"name"`
}

type orgUserRecord struct {
	OrgId     string `dynamodbav:"pk"`
	UserIdKey string `dynamodbav:"sk"`
	Role      role   `dynamodbav:"role"`
}

func NewOrgRepo() (orgRepo, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		return orgRepo{}, fmt.Errorf("failed to create session config: %w", err)
	}

	return orgRepo{
		db:               dynamodb.NewFromConfig(cfg),
		table:            aws.String("expense-system-records"),
		reverseLookupGsi: aws.String("reverse-lookup"),
	}, nil
}

type UserOrg struct {
	Name string
	Id   string
	Role role
}

func (uo UserOrg) IsAdmin() bool {
	return uo.Role == RoleAdmin
}

func (o orgRepo) getOrgsForUser(userId string) ([]UserOrg, error) {
	res, err := o.db.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              o.table,
		IndexName:              o.reverseLookupGsi,
		KeyConditionExpression: aws.String("sk = :userId and begins_with(pk, :orgPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId":    &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
			":orgPrefix": &types.AttributeValueMemberS{Value: "ORG#"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orgs from dynamodb for user %v: %w", userId, err)
	}

	records := []orgUserRecord{}
	attributevalue.UnmarshalListOfMaps(res.Items, &records)

	orgs := []UserOrg{}
	for _, r := range records {
		orgId := strings.Split(r.OrgId, "#")[1]
		orgName, err := o.getOrgName(orgId)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve name for org %v: %w", orgId, err)
		}
		orgs = append(orgs, UserOrg{
			Id:   orgId,
			Name: orgName,
			Role: r.Role,
		})
	}

	return orgs, nil
}

func (o orgRepo) getOrgName(orgId string) (string, error) {
	res, err := o.db.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: o.table,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ORG#%s", orgId)},
			"sk": &types.AttributeValueMemberS{Value: "ORG"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to retrieve name for org %v: %w", orgId, err)
	}

	or := orgRecord{}
	attributevalue.UnmarshalMap(res.Item, &or)

	return or.Name, nil
}

func (o orgRepo) createOrg(name, adminId string) (string, error) {
	id, err := newId()
	if err != nil {
		return "", fmt.Errorf("failed to generate new org id: %w", err)
	}

	orgItem, err := attributevalue.MarshalMap(orgRecord{
		OrgId:         fmt.Sprintf("ORG#%v", id),
		OrgPrimaryKey: "ORG",
		Name:          name,
	})

	if err != nil {
		return "", fmt.Errorf("failed to marshal new org record: %w", err)
	}

	orgAdminItem, err := attributevalue.MarshalMap(orgUserRecord{
		UserIdKey: fmt.Sprintf("USER#%v", adminId),
		OrgId:     fmt.Sprintf("ORG#%v", id),
		Role:      RoleAdmin,
	})

	if err != nil {
		return "", fmt.Errorf("failed to marshal new org admin record: %w", err)
	}

	_, err = o.db.TransactWriteItems(context.Background(), &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName: o.table,
					Item:      orgItem,
				},
			},
			{
				Put: &types.Put{
					TableName: o.table,
					Item:      orgAdminItem,
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to save new org and admin records: %w", err)
	}

	return id, nil
}

func newId() (string, error) {
	id, err := uuid.NewV4()
	return id.String(), err
}
