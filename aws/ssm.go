package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMService struct{}

// AwsSsmGetParam fetches a single SSM parameter value.
func (s *SSMService) AwsSsmGetParam(ctx context.Context, path string) (string, error) {
	ssmClient, err := GetAWSService("ap-northeast-2").GetSSMClient()
	if err != nil {
		return "", fmt.Errorf("failed to initialize SSM client: %w", err)
	}

	param, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(path),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch SSM parameter: %w", err)
	}

	return aws.ToString(param.Parameter.Value), nil
}

// AwsSsmGetParams fetches multiple SSM parameter values.
func (s *SSMService) AwsSsmGetParams(ctx context.Context, paths []string) ([]string, error) {
	ssmClient, err := GetAWSService("ap-northeast-2").GetSSMClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SSM client: %w", err)
	}

	params, err := ssmClient.GetParameters(ctx, &ssm.GetParametersInput{
		Names:          paths,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch SSM parameters: %w", err)
	}

	// Handle invalid parameters
	if len(params.InvalidParameters) > 0 {
		return nil, fmt.Errorf("invalid SSM parameters: %v", params.InvalidParameters)
	}

	var values []string
	for _, param := range params.Parameters {
		values = append(values, aws.ToString(param.Value))
	}
	return values, nil
}
