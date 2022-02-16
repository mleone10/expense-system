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
	db *dynamodb.Client
}

type org struct {
	Name string
	Id   string
}

type tableRecord struct {
	Pk string
	Sk string
}

type orgRecord struct {
	Pk   string `dynamodbav:"pk"`
	Sk   string `dynamodbav:"sk"`
	Name string `dynamodbav:"name"`
}

const tableName string = "expense-system-records"

func NewOrgRepo() (orgRepo, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		return orgRepo{}, fmt.Errorf("failed to create session config: %w", err)
	}

	return orgRepo{
		db: dynamodb.NewFromConfig(cfg),
	}, nil
}

func (o orgRepo) getOrgsForUser(userId string) ([]org, error) {
	res, err := o.db.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("pk = :userId and begins_with(sk, :membershipPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId":           &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
			":membershipPrefix": &types.AttributeValueMemberS{Value: "ORG#"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orgs from dynamodb for user %v: %w", userId, err)
	}

	trs := []tableRecord{}
	attributevalue.UnmarshalListOfMaps(res.Items, &trs)

	orgs := []org{}
	for _, tr := range trs {
		orgId := strings.Split(tr.Sk, "#")[1]
		orgName, err := o.getOrgName(orgId)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve name for org %v: %w", orgId, err)
		}
		orgs = append(orgs, org{
			Id:   orgId,
			Name: orgName,
		})
	}

	return orgs, nil
}

func (o orgRepo) getOrgName(orgId string) (string, error) {
	res, err := o.db.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
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

func (o orgRepo) createOrg(name string) (string, error) {
	id, err := newId()
	if err != nil {
		return "", fmt.Errorf("failed to generate new org id: %w", err)
	}

	item, err := attributevalue.MarshalMap(orgRecord{
		Pk:   fmt.Sprintf("ORG#%v", id),
		Sk:   "ORG",
		Name: name,
	})

	// TODO: Org creation is working, just need to link to active user with a second write
	_, err = o.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return "", fmt.Errorf("failed to write new organization to database: %w", err)
	}

	return id, nil
}

func newId() (string, error) {
	id, err := uuid.NewV4()
	return id.String(), err
}
