package aws

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type STSProvisioner struct {
	TOTPCode  string
	MFASerial string
}

func (p STSProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	config := aws.NewConfig()
	config.Credentials = credentials.NewStaticCredentialsProvider(in.ItemFields[fieldname.AccessKeyID], in.ItemFields[fieldname.SecretAccessKey], "")
	config.Region = in.ItemFields[FieldNameDefaultRegion]
	stsProvider := sts.NewFromConfig(*config)
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int32(900), // minimum expiration time - 15 minutes
		SerialNumber:    aws.String(p.MFASerial),
		TokenCode:       aws.String(p.TOTPCode),
	}

	result, err := stsProvider.GetSessionToken(ctx, input)
	if err != nil {
		out.AddError(err)
		return
	}
	out.AddEnvVar("AWS_ACCESS_KEY_ID", *result.Credentials.AccessKeyId)
	out.AddEnvVar("AWS_SECRET_ACCESS_KEY", *result.Credentials.SecretAccessKey)
	out.AddEnvVar("AWS_SESSION_TOKEN", *result.Credentials.SessionToken)
	if region, ok := in.ItemFields[FieldNameDefaultRegion]; ok {
		out.AddEnvVar("AWS_DEFAULT_REGION", region)
	}

}

func (p STSProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p STSProvisioner) Description() string {
	return "Provision environment variables with the temporary credentials AWS_ACCESS_KEY_ID, AWS_ACCESS_KEY_ID, AWS_ACCESS_KEY_ID"
}
