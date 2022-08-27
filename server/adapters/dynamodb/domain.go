package dynamodb

const (
	prefixOrg  = "ORG#"
	prefixUser = "USER#"
)

type orgRecord struct {
	orgId   string `dynamodbav:"pk"`
	orgFlag string `dynamodbav:"sk"`
	name    string `dynamodbav:"name"`
}

type orgUserRecord struct {
	orgId     string `dynamodbav:"pk"`
	userIdKey string `dynamodbav:"sk"`
	admin     bool   `dynamodbav:"admin"`
}
