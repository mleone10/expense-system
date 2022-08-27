package dynamodb

const (
	prefixOrg  = "ORG#"
	prefixUser = "USER#"
)

type orgRecord struct {
	OrgId   string `dynamodbav:"pk"`
	OrgFlag string `dynamodbav:"sk"`
	Name    string `dynamodbav:"name"`
}

type orgUserRecord struct {
	OrgId     string `dynamodbav:"pk"`
	UserIdKey string `dynamodbav:"sk"`
	Admin     bool   `dynamodbav:"admin"`
}
