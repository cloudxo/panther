package table

/**
 * Panther is a Cloud-Native SIEM for the Modern Security Team.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/panther-labs/panther/api/lambda/alerts/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// UpdateAlertStatus - updates the alert details and returns the updated item
func (table *AlertsTable) UpdateAlertStatus(input *models.UpdateAlertStatusInput) (*AlertItem, error) {
	// Create the dynamo key we want to update
	var alertKey = DynamoItem{AlertIDKey: {S: aws.String(*input.AlertID)}}

	// Create the update builder
	updateBuilder := createUpdateBuilder(input)

	// Create the condition builder
	conditionBuilder := createConditionBuilder(input)

	// Build an expression from our builders
	expression, err := buildExpression(updateBuilder, conditionBuilder)
	if err != nil {
		return nil, err
	}

	// Create our dynamo update item
	updateItem := dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expression.Names(),
		ExpressionAttributeValues: expression.Values(),
		Key:                       alertKey,
		ReturnValues:              aws.String("ALL_NEW"),
		TableName:                 &table.AlertsTableName,
		UpdateExpression:          expression.Update(),
		ConditionExpression:       expression.Condition(),
	}

	// Run the update query and marshal
	updatedAlert := &AlertItem{}
	if err = table.update(updateItem, &updatedAlert); err != nil {
		return nil, err
	}

	return updatedAlert, nil
}

// createUpdateBuilder - creates an update builder
func createUpdateBuilder(input *models.UpdateAlertStatusInput) expression.UpdateBuilder {
	// When settig an "open" status we actually remove the attribute
	// for uniformity against previous items in the database
	// which also do not have a status attribute.
	if *input.Status == models.OpenStatus {
		return expression.
			Remove(expression.Name(StatusKey)).
			Set(expression.Name(LastUpdatedByKey), expression.Value(input.UserID)).
			Set(expression.Name(LastUpdatedByTimeKey), expression.Value(aws.Time(time.Now().UTC())))
	}

	return expression.
		Set(expression.Name(StatusKey), expression.Value(input.Status)).
		Set(expression.Name(LastUpdatedByKey), expression.Value(input.UserID)).
		Set(expression.Name(LastUpdatedByTimeKey), expression.Value(aws.Time(time.Now().UTC())))
}

// createConditionBuilder - creates a condition builder
func createConditionBuilder(input *models.UpdateAlertStatusInput) expression.ConditionBuilder {
	return expression.Equal(expression.Name(AlertIDKey), expression.Value(input.AlertID))
}

// buildExpression - builds an expression
func buildExpression(
	updateBuilder expression.UpdateBuilder,
	conditionBuilder expression.ConditionBuilder,
) (expression.Expression, error) {

	expr, err := expression.
		NewBuilder().
		WithUpdate(updateBuilder).
		WithCondition(conditionBuilder).
		Build()
	if err != nil {
		return expr, &genericapi.InternalError{
			Message: "failed to build update expression: " + err.Error()}
	}
	return expr, nil
}

// table.update - runs an update query
func (table *AlertsTable) update(
	item dynamodb.UpdateItemInput,
	newItem interface{},
) error {

	response, err := table.Client.UpdateItem(&item)

	if err != nil {
		return &genericapi.AWSError{Method: "dynamodb.UpdateItem", Err: err}
	}

	if err = dynamodbattribute.UnmarshalMap(response.Attributes, newItem); err != nil {
		return &genericapi.InternalError{Message: "failed to unmarshal dynamo item: " + err.Error()}
	}
	return nil
}
