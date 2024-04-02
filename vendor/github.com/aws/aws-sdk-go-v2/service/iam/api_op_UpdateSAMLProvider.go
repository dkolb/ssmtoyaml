// Code generated by smithy-go-codegen DO NOT EDIT.

package iam

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Updates the metadata document for an existing SAML provider resource object.
// This operation requires Signature Version 4 (https://docs.aws.amazon.com/general/latest/gr/signature-version-4.html)
// .
func (c *Client) UpdateSAMLProvider(ctx context.Context, params *UpdateSAMLProviderInput, optFns ...func(*Options)) (*UpdateSAMLProviderOutput, error) {
	if params == nil {
		params = &UpdateSAMLProviderInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "UpdateSAMLProvider", params, optFns, c.addOperationUpdateSAMLProviderMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*UpdateSAMLProviderOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type UpdateSAMLProviderInput struct {

	// An XML document generated by an identity provider (IdP) that supports SAML 2.0.
	// The document includes the issuer's name, expiration information, and keys that
	// can be used to validate the SAML authentication response (assertions) that are
	// received from the IdP. You must generate the metadata document using the
	// identity management software that is used as your organization's IdP.
	//
	// This member is required.
	SAMLMetadataDocument *string

	// The Amazon Resource Name (ARN) of the SAML provider to update. For more
	// information about ARNs, see Amazon Resource Names (ARNs) (https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html)
	// in the Amazon Web Services General Reference.
	//
	// This member is required.
	SAMLProviderArn *string

	noSmithyDocumentSerde
}

// Contains the response to a successful UpdateSAMLProvider request.
type UpdateSAMLProviderOutput struct {

	// The Amazon Resource Name (ARN) of the SAML provider that was updated.
	SAMLProviderArn *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationUpdateSAMLProviderMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsquery_serializeOpUpdateSAMLProvider{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpUpdateSAMLProvider{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "UpdateSAMLProvider"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addOpUpdateSAMLProviderValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opUpdateSAMLProvider(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opUpdateSAMLProvider(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "UpdateSAMLProvider",
	}
}
