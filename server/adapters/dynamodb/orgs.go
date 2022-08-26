package dynamodb

// func (o orgRepo) getOrgName(orgId string) (string, error) {
// 	res, err := o.db.GetItem(context.Background(), &dynamodb.GetItemInput{
// 		TableName: o.table,
// 		Key: map[string]types.AttributeValue{
// 			"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ORG#%s", orgId)},
// 			"sk": &types.AttributeValueMemberS{Value: "ORG"},
// 		},
// 	})
// 	if err != nil {
// 		return "", fmt.Errorf("failed to retrieve name for org %v: %w", orgId, err)
// 	}

// 	or := orgRecord{}
// 	attributevalue.UnmarshalMap(res.Item, &or)

// 	return or.Name, nil
// }

// func (o orgRepo) createOrg(name, adminId string) (string, error) {
// 	id, err := newId()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to generate new org id: %w", err)
// 	}

// 	orgItem, err := attributevalue.MarshalMap(orgRecord{
// 		OrgId:         fmt.Sprintf("ORG#%v", id),
// 		OrgPrimaryKey: "ORG",
// 		Name:          name,
// 	})

// 	if err != nil {
// 		return "", fmt.Errorf("failed to marshal new org record: %w", err)
// 	}

// 	orgAdminItem, err := attributevalue.MarshalMap(orgUserRecord{
// 		UserIdKey: fmt.Sprintf("USER#%v", adminId),
// 		OrgId:     fmt.Sprintf("ORG#%v", id),
// 		Role:      RoleAdmin,
// 	})

// 	if err != nil {
// 		return "", fmt.Errorf("failed to marshal new org admin record: %w", err)
// 	}

// 	_, err = o.db.TransactWriteItems(context.Background(), &dynamodb.TransactWriteItemsInput{
// 		TransactItems: []types.TransactWriteItem{
// 			{
// 				Put: &types.Put{
// 					TableName: o.table,
// 					Item:      orgItem,
// 				},
// 			},
// 			{
// 				Put: &types.Put{
// 					TableName: o.table,
// 					Item:      orgAdminItem,
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		return "", fmt.Errorf("failed to save new org and admin records: %w", err)
// 	}

// 	return id, nil
// }

// func newId() (string, error) {
// 	id, err := uuid.NewV4()
// 	return id.String(), err
// }
