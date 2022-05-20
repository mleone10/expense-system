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

type orgRepo struct {
	db    *dynamodb.Client
	table *string
}

type orgRecord struct {
	Pk   string `dynamodbav:"pk"`
	Sk   string `dynamodbav:"sk"`
	Name string `dynamodbav:"name"`
}

type userOrgRecord struct {
	Pk    string `dynamodbav:"pk"`
	Sk    string `dynamodbav:"sk"`
	Admin bool   `dynamodbav:"admin"`
}

func NewOrgRepo() (orgRepo, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		return orgRepo{}, fmt.Errorf("failed to create session config: %w", err)
	}

	return orgRepo{
		db:    dynamodb.NewFromConfig(cfg),
		table: aws.String("expense-system-records"),
	}, nil
}

type UserOrg struct {
	Name  string
	Id    string
	Admin bool
}

func (o orgRepo) getOrgsForUser(userId string) ([]UserOrg, error) {
	res, err := o.db.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              o.table,
		KeyConditionExpression: aws.String("pk = :userId and begins_with(sk, :membershipPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId":           &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
			":membershipPrefix": &types.AttributeValueMemberS{Value: "ORG#"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orgs from dynamodb for user %v: %w", userId, err)
	}

	uors := []userOrgRecord{}
	attributevalue.UnmarshalListOfMaps(res.Items, &uors)

	orgs := []UserOrg{}
	for _, uor := range uors {
		orgId := strings.Split(uor.Sk, "#")[1]
		orgName, err := o.getOrgName(orgId)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve name for org %v: %w", orgId, err)
		}
		orgs = append(orgs, UserOrg{
			Id:    orgId,
			Name:  orgName,
			Admin: uor.Admin,
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

func (o orgRepo) createOrg(name, admin string) (string, error) {
	id, err := newId()
	if err != nil {
		return "", fmt.Errorf("failed to generate new org id: %w", err)
	}

	orgItem, err := attributevalue.MarshalMap(orgRecord{
		Pk:   fmt.Sprintf("ORG#%v", id),
		Sk:   "ORG",
		Name: name,
	})

	if err != nil {
		return "", fmt.Errorf("failed to marshal new org record: %w", err)
	}

	orgAdminItem, err := attributevalue.MarshalMap(userOrgRecord{
		Pk:    fmt.Sprintf("USER#%v", admin),
		Sk:    fmt.Sprintf("ORG#%v", id),
		Admin: true,
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

	return id, nil
}

func newId() (string, error) {
	id, err := uuid.NewV4()
	return id.String(), err
}
