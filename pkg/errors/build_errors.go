// Package errors contains errors for the domain "nt".
//
// This file is automatically generated by errawr-gen. Do not modify it.
package errors

import (
	errawr "github.com/puppetlabs/errawr-go/v2/pkg/errawr"
	impl "github.com/puppetlabs/errawr-go/v2/pkg/impl"
)

// Error is the type of all errors generated by this package.
type Error interface {
	errawr.Error
}

// External contains methods that can be used externally to help consume errors from this package.
type External struct{}

// API is a singleton instance of the External type.
var API External

// Domain is the general domain in which all errors in this package belong.
var Domain = &impl.ErrorDomain{
	Key:   "nt",
	Title: "Nebula Tasks",
}

// SecretsSection defines a section of errors with the following scope:
// Secrets errors
var SecretsSection = &impl.ErrorSection{
	Key:   "secrets",
	Title: "Secrets errors",
}

// SecretsGetErrorCode is the code for an instance of "get_error".
const SecretsGetErrorCode = "nt_secrets_get_error"

// IsSecretsGetError tests whether a given error is an instance of "get_error".
func IsSecretsGetError(err errawr.Error) bool {
	return err != nil && err.Is(SecretsGetErrorCode)
}

// IsSecretsGetError tests whether a given error is an instance of "get_error".
func (External) IsSecretsGetError(err errawr.Error) bool {
	return IsSecretsGetError(err)
}

// SecretsGetErrorBuilder is a builder for "get_error" errors.
type SecretsGetErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "get_error" from this builder.
func (b *SecretsGetErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "failed to get the secret at key {{key}}",
		Technical: "failed to get the secret at key {{key}}",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "get_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     SecretsSection,
		ErrorSensitivity: errawr.ErrorSensitivityNone,
		ErrorTitle:       "Get error",
		Version:          1,
	}
}

// NewSecretsGetErrorBuilder creates a new error builder for the code "get_error".
func NewSecretsGetErrorBuilder(key string) *SecretsGetErrorBuilder {
	return &SecretsGetErrorBuilder{arguments: impl.ErrorArguments{"key": impl.NewErrorArgument(key, "the key that had the failure")}}
}

// NewSecretsGetError creates a new error with the code "get_error".
func NewSecretsGetError(key string) Error {
	return NewSecretsGetErrorBuilder(key).Build()
}

// SecretsK8sServiceAccountTokenReadErrorCode is the code for an instance of "k8s_service_account_token_read_error".
const SecretsK8sServiceAccountTokenReadErrorCode = "nt_secrets_k8s_service_account_token_read_error"

// IsSecretsK8sServiceAccountTokenReadError tests whether a given error is an instance of "k8s_service_account_token_read_error".
func IsSecretsK8sServiceAccountTokenReadError(err errawr.Error) bool {
	return err != nil && err.Is(SecretsK8sServiceAccountTokenReadErrorCode)
}

// IsSecretsK8sServiceAccountTokenReadError tests whether a given error is an instance of "k8s_service_account_token_read_error".
func (External) IsSecretsK8sServiceAccountTokenReadError(err errawr.Error) bool {
	return IsSecretsK8sServiceAccountTokenReadError(err)
}

// SecretsK8sServiceAccountTokenReadErrorBuilder is a builder for "k8s_service_account_token_read_error" errors.
type SecretsK8sServiceAccountTokenReadErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "k8s_service_account_token_read_error" from this builder.
func (b *SecretsK8sServiceAccountTokenReadErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "there was an error reading the service account token from the secret file",
		Technical: "there was an error reading the service account token from the secret file",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "k8s_service_account_token_read_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     SecretsSection,
		ErrorSensitivity: errawr.ErrorSensitivityNone,
		ErrorTitle:       "K8s service account token read error",
		Version:          1,
	}
}

// NewSecretsK8sServiceAccountTokenReadErrorBuilder creates a new error builder for the code "k8s_service_account_token_read_error".
func NewSecretsK8sServiceAccountTokenReadErrorBuilder() *SecretsK8sServiceAccountTokenReadErrorBuilder {
	return &SecretsK8sServiceAccountTokenReadErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewSecretsK8sServiceAccountTokenReadError creates a new error with the code "k8s_service_account_token_read_error".
func NewSecretsK8sServiceAccountTokenReadError() Error {
	return NewSecretsK8sServiceAccountTokenReadErrorBuilder().Build()
}

// SecretsKeyNotFoundCode is the code for an instance of "key_not_found".
const SecretsKeyNotFoundCode = "nt_secrets_key_not_found"

// IsSecretsKeyNotFound tests whether a given error is an instance of "key_not_found".
func IsSecretsKeyNotFound(err errawr.Error) bool {
	return err != nil && err.Is(SecretsKeyNotFoundCode)
}

// IsSecretsKeyNotFound tests whether a given error is an instance of "key_not_found".
func (External) IsSecretsKeyNotFound(err errawr.Error) bool {
	return IsSecretsKeyNotFound(err)
}

// SecretsKeyNotFoundBuilder is a builder for "key_not_found" errors.
type SecretsKeyNotFoundBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "key_not_found" from this builder.
func (b *SecretsKeyNotFoundBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "key {{key}} not found",
		Technical: "key {{key}} not found",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "key_not_found",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     SecretsSection,
		ErrorSensitivity: errawr.ErrorSensitivityNone,
		ErrorTitle:       "Key not found",
		Version:          1,
	}
}

// NewSecretsKeyNotFoundBuilder creates a new error builder for the code "key_not_found".
func NewSecretsKeyNotFoundBuilder(key string) *SecretsKeyNotFoundBuilder {
	return &SecretsKeyNotFoundBuilder{arguments: impl.ErrorArguments{"key": impl.NewErrorArgument(key, "the key that is missing")}}
}

// NewSecretsKeyNotFound creates a new error with the code "key_not_found".
func NewSecretsKeyNotFound(key string) Error {
	return NewSecretsKeyNotFoundBuilder(key).Build()
}

// SecretsMalformedValueCode is the code for an instance of "malformed_value".
const SecretsMalformedValueCode = "nt_secrets_malformed_value"

// IsSecretsMalformedValue tests whether a given error is an instance of "malformed_value".
func IsSecretsMalformedValue(err errawr.Error) bool {
	return err != nil && err.Is(SecretsMalformedValueCode)
}

// IsSecretsMalformedValue tests whether a given error is an instance of "malformed_value".
func (External) IsSecretsMalformedValue(err errawr.Error) bool {
	return IsSecretsMalformedValue(err)
}

// SecretsMalformedValueBuilder is a builder for "malformed_value" errors.
type SecretsMalformedValueBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "malformed_value" from this builder.
func (b *SecretsMalformedValueBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "the value for secretRef is not a string",
		Technical: "the value for secretRef is not a string",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "malformed_value",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     SecretsSection,
		ErrorSensitivity: errawr.ErrorSensitivityNone,
		ErrorTitle:       "Malformed value",
		Version:          1,
	}
}

// NewSecretsMalformedValueBuilder creates a new error builder for the code "malformed_value".
func NewSecretsMalformedValueBuilder() *SecretsMalformedValueBuilder {
	return &SecretsMalformedValueBuilder{arguments: impl.ErrorArguments{}}
}

// NewSecretsMalformedValue creates a new error with the code "malformed_value".
func NewSecretsMalformedValue() Error {
	return NewSecretsMalformedValueBuilder().Build()
}

// SecretsMissingSecretRefCode is the code for an instance of "missing_secret_ref".
const SecretsMissingSecretRefCode = "nt_secrets_missing_secret_ref"

// IsSecretsMissingSecretRef tests whether a given error is an instance of "missing_secret_ref".
func IsSecretsMissingSecretRef(err errawr.Error) bool {
	return err != nil && err.Is(SecretsMissingSecretRefCode)
}

// IsSecretsMissingSecretRef tests whether a given error is an instance of "missing_secret_ref".
func (External) IsSecretsMissingSecretRef(err errawr.Error) bool {
	return IsSecretsMissingSecretRef(err)
}

// SecretsMissingSecretRefBuilder is a builder for "missing_secret_ref" errors.
type SecretsMissingSecretRefBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "missing_secret_ref" from this builder.
func (b *SecretsMissingSecretRefBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "secretRef is missing",
		Technical: "secretRef is missing",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "missing_secret_ref",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     SecretsSection,
		ErrorSensitivity: errawr.ErrorSensitivityNone,
		ErrorTitle:       "Missing secretRef",
		Version:          1,
	}
}

// NewSecretsMissingSecretRefBuilder creates a new error builder for the code "missing_secret_ref".
func NewSecretsMissingSecretRefBuilder() *SecretsMissingSecretRefBuilder {
	return &SecretsMissingSecretRefBuilder{arguments: impl.ErrorArguments{}}
}

// NewSecretsMissingSecretRef creates a new error with the code "missing_secret_ref".
func NewSecretsMissingSecretRef() Error {
	return NewSecretsMissingSecretRefBuilder().Build()
}

// SecretsSessionSetupErrorCode is the code for an instance of "session_setup_error".
const SecretsSessionSetupErrorCode = "nt_secrets_session_setup_error"

// IsSecretsSessionSetupError tests whether a given error is an instance of "session_setup_error".
func IsSecretsSessionSetupError(err errawr.Error) bool {
	return err != nil && err.Is(SecretsSessionSetupErrorCode)
}

// IsSecretsSessionSetupError tests whether a given error is an instance of "session_setup_error".
func (External) IsSecretsSessionSetupError(err errawr.Error) bool {
	return IsSecretsSessionSetupError(err)
}

// SecretsSessionSetupErrorBuilder is a builder for "session_setup_error" errors.
type SecretsSessionSetupErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "session_setup_error" from this builder.
func (b *SecretsSessionSetupErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "there was an error setting up the secret engine session",
		Technical: "there was an error setting up the secret engine session",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "session_setup_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     SecretsSection,
		ErrorSensitivity: errawr.ErrorSensitivityNone,
		ErrorTitle:       "Session setup error",
		Version:          1,
	}
}

// NewSecretsSessionSetupErrorBuilder creates a new error builder for the code "session_setup_error".
func NewSecretsSessionSetupErrorBuilder() *SecretsSessionSetupErrorBuilder {
	return &SecretsSessionSetupErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewSecretsSessionSetupError creates a new error with the code "session_setup_error".
func NewSecretsSessionSetupError() Error {
	return NewSecretsSessionSetupErrorBuilder().Build()
}

// SecretsVaultLoginErrorCode is the code for an instance of "vault_login_error".
const SecretsVaultLoginErrorCode = "nt_secrets_vault_login_error"

// IsSecretsVaultLoginError tests whether a given error is an instance of "vault_login_error".
func IsSecretsVaultLoginError(err errawr.Error) bool {
	return err != nil && err.Is(SecretsVaultLoginErrorCode)
}

// IsSecretsVaultLoginError tests whether a given error is an instance of "vault_login_error".
func (External) IsSecretsVaultLoginError(err errawr.Error) bool {
	return IsSecretsVaultLoginError(err)
}

// SecretsVaultLoginErrorBuilder is a builder for "vault_login_error" errors.
type SecretsVaultLoginErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "vault_login_error" from this builder.
func (b *SecretsVaultLoginErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "there was an error logging into the vault server",
		Technical: "there was an error logging into the vault server",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "vault_login_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     SecretsSection,
		ErrorSensitivity: errawr.ErrorSensitivityNone,
		ErrorTitle:       "Vault login error",
		Version:          1,
	}
}

// NewSecretsVaultLoginErrorBuilder creates a new error builder for the code "vault_login_error".
func NewSecretsVaultLoginErrorBuilder() *SecretsVaultLoginErrorBuilder {
	return &SecretsVaultLoginErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewSecretsVaultLoginError creates a new error with the code "vault_login_error".
func NewSecretsVaultLoginError() Error {
	return NewSecretsVaultLoginErrorBuilder().Build()
}

// SecretsVaultSetupErrorCode is the code for an instance of "vault_setup_error".
const SecretsVaultSetupErrorCode = "nt_secrets_vault_setup_error"

// IsSecretsVaultSetupError tests whether a given error is an instance of "vault_setup_error".
func IsSecretsVaultSetupError(err errawr.Error) bool {
	return err != nil && err.Is(SecretsVaultSetupErrorCode)
}

// IsSecretsVaultSetupError tests whether a given error is an instance of "vault_setup_error".
func (External) IsSecretsVaultSetupError(err errawr.Error) bool {
	return IsSecretsVaultSetupError(err)
}

// SecretsVaultSetupErrorBuilder is a builder for "vault_setup_error" errors.
type SecretsVaultSetupErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "vault_setup_error" from this builder.
func (b *SecretsVaultSetupErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "there was an error setting up the Vault client",
		Technical: "there was an error setting up the Vault client",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "vault_setup_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     SecretsSection,
		ErrorSensitivity: errawr.ErrorSensitivityNone,
		ErrorTitle:       "Vault client setup error",
		Version:          1,
	}
}

// NewSecretsVaultSetupErrorBuilder creates a new error builder for the code "vault_setup_error".
func NewSecretsVaultSetupErrorBuilder() *SecretsVaultSetupErrorBuilder {
	return &SecretsVaultSetupErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewSecretsVaultSetupError creates a new error with the code "vault_setup_error".
func NewSecretsVaultSetupError() Error {
	return NewSecretsVaultSetupErrorBuilder().Build()
}

// ServerSection defines a section of errors with the following scope:
// Sever errors
var ServerSection = &impl.ErrorSection{
	Key:   "server",
	Title: "Sever errors",
}

// ServerConfigMapJSONErrorCode is the code for an instance of "config_map_json_error".
const ServerConfigMapJSONErrorCode = "nt_server_config_map_json_error"

// IsServerConfigMapJSONError tests whether a given error is an instance of "config_map_json_error".
func IsServerConfigMapJSONError(err errawr.Error) bool {
	return err != nil && err.Is(ServerConfigMapJSONErrorCode)
}

// IsServerConfigMapJSONError tests whether a given error is an instance of "config_map_json_error".
func (External) IsServerConfigMapJSONError(err errawr.Error) bool {
	return IsServerConfigMapJSONError(err)
}

// ServerConfigMapJSONErrorBuilder is a builder for "config_map_json_error" errors.
type ServerConfigMapJSONErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "config_map_json_error" from this builder.
func (b *ServerConfigMapJSONErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "error when parsing the \"spec.json\" field of the {{namespace}}/{{name}} config map data",
		Technical: "error when parsing the \"spec.json\" field of the {{namespace}}/{{name}} config map data",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "config_map_json_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     ServerSection,
		ErrorSensitivity: errawr.ErrorSensitivityNone,
		ErrorTitle:       "ConfigMap spec.json parse error",
		Version:          1,
	}
}

// NewServerConfigMapJSONErrorBuilder creates a new error builder for the code "config_map_json_error".
func NewServerConfigMapJSONErrorBuilder(name string, namespace string) *ServerConfigMapJSONErrorBuilder {
	return &ServerConfigMapJSONErrorBuilder{arguments: impl.ErrorArguments{
		"name":      impl.NewErrorArgument(name, "The name of the config map"),
		"namespace": impl.NewErrorArgument(namespace, "The namespace of the config map"),
	}}
}

// NewServerConfigMapJSONError creates a new error with the code "config_map_json_error".
func NewServerConfigMapJSONError(name string, namespace string) Error {
	return NewServerConfigMapJSONErrorBuilder(name, namespace).Build()
}

// ServerGetConfigMapErrorCode is the code for an instance of "get_config_map_error".
const ServerGetConfigMapErrorCode = "nt_server_get_config_map_error"

// IsServerGetConfigMapError tests whether a given error is an instance of "get_config_map_error".
func IsServerGetConfigMapError(err errawr.Error) bool {
	return err != nil && err.Is(ServerGetConfigMapErrorCode)
}

// IsServerGetConfigMapError tests whether a given error is an instance of "get_config_map_error".
func (External) IsServerGetConfigMapError(err errawr.Error) bool {
	return IsServerGetConfigMapError(err)
}

// ServerGetConfigMapErrorBuilder is a builder for "get_config_map_error" errors.
type ServerGetConfigMapErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "get_config_map_error" from this builder.
func (b *ServerGetConfigMapErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "error when getting config map {{namespace}}/{{name}}",
		Technical: "error when getting config map {{namespace}}/{{name}}",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "get_config_map_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     ServerSection,
		ErrorSensitivity: errawr.ErrorSensitivityNone,
		ErrorTitle:       "ConfigMap GET failed",
		Version:          1,
	}
}

// NewServerGetConfigMapErrorBuilder creates a new error builder for the code "get_config_map_error".
func NewServerGetConfigMapErrorBuilder(name string, namespace string) *ServerGetConfigMapErrorBuilder {
	return &ServerGetConfigMapErrorBuilder{arguments: impl.ErrorArguments{
		"name":      impl.NewErrorArgument(name, "The name of the config map"),
		"namespace": impl.NewErrorArgument(namespace, "The namespace of the config map"),
	}}
}

// NewServerGetConfigMapError creates a new error with the code "get_config_map_error".
func NewServerGetConfigMapError(name string, namespace string) Error {
	return NewServerGetConfigMapErrorBuilder(name, namespace).Build()
}

// ServerInClusterConfigErrorCode is the code for an instance of "in_cluster_config_error".
const ServerInClusterConfigErrorCode = "nt_server_in_cluster_config_error"

// IsServerInClusterConfigError tests whether a given error is an instance of "in_cluster_config_error".
func IsServerInClusterConfigError(err errawr.Error) bool {
	return err != nil && err.Is(ServerInClusterConfigErrorCode)
}

// IsServerInClusterConfigError tests whether a given error is an instance of "in_cluster_config_error".
func (External) IsServerInClusterConfigError(err errawr.Error) bool {
	return IsServerInClusterConfigError(err)
}

// ServerInClusterConfigErrorBuilder is a builder for "in_cluster_config_error" errors.
type ServerInClusterConfigErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "in_cluster_config_error" from this builder.
func (b *ServerInClusterConfigErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "error fetching the in cluster config",
		Technical: "error fetching the in cluster config",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "in_cluster_config_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     ServerSection,
		ErrorSensitivity: errawr.ErrorSensitivityNone,
		ErrorTitle:       "Config fetch error",
		Version:          1,
	}
}

// NewServerInClusterConfigErrorBuilder creates a new error builder for the code "in_cluster_config_error".
func NewServerInClusterConfigErrorBuilder() *ServerInClusterConfigErrorBuilder {
	return &ServerInClusterConfigErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewServerInClusterConfigError creates a new error with the code "in_cluster_config_error".
func NewServerInClusterConfigError() Error {
	return NewServerInClusterConfigErrorBuilder().Build()
}

// ServerNewK8sClientErrorCode is the code for an instance of "new_k8s_client_error".
const ServerNewK8sClientErrorCode = "nt_server_new_k8s_client_error"

// IsServerNewK8sClientError tests whether a given error is an instance of "new_k8s_client_error".
func IsServerNewK8sClientError(err errawr.Error) bool {
	return err != nil && err.Is(ServerNewK8sClientErrorCode)
}

// IsServerNewK8sClientError tests whether a given error is an instance of "new_k8s_client_error".
func (External) IsServerNewK8sClientError(err errawr.Error) bool {
	return IsServerNewK8sClientError(err)
}

// ServerNewK8sClientErrorBuilder is a builder for "new_k8s_client_error" errors.
type ServerNewK8sClientErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "new_k8s_client_error" from this builder.
func (b *ServerNewK8sClientErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "error creating a Kubernetes client",
		Technical: "error creating a Kubernetes client",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "new_k8s_client_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     ServerSection,
		ErrorSensitivity: errawr.ErrorSensitivityNone,
		ErrorTitle:       "Kubernetes client create error",
		Version:          1,
	}
}

// NewServerNewK8sClientErrorBuilder creates a new error builder for the code "new_k8s_client_error".
func NewServerNewK8sClientErrorBuilder() *ServerNewK8sClientErrorBuilder {
	return &ServerNewK8sClientErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewServerNewK8sClientError creates a new error with the code "new_k8s_client_error".
func NewServerNewK8sClientError() Error {
	return NewServerNewK8sClientErrorBuilder().Build()
}

// ServerRunErrorCode is the code for an instance of "run_error".
const ServerRunErrorCode = "nt_server_run_error"

// IsServerRunError tests whether a given error is an instance of "run_error".
func IsServerRunError(err errawr.Error) bool {
	return err != nil && err.Is(ServerRunErrorCode)
}

// IsServerRunError tests whether a given error is an instance of "run_error".
func (External) IsServerRunError(err errawr.Error) bool {
	return IsServerRunError(err)
}

// ServerRunErrorBuilder is a builder for "run_error" errors.
type ServerRunErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "run_error" from this builder.
func (b *ServerRunErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "an error occurred which running the server",
		Technical: "an error occurred which running the server",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "run_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     ServerSection,
		ErrorSensitivity: errawr.ErrorSensitivityNone,
		ErrorTitle:       "Run error",
		Version:          1,
	}
}

// NewServerRunErrorBuilder creates a new error builder for the code "run_error".
func NewServerRunErrorBuilder() *ServerRunErrorBuilder {
	return &ServerRunErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewServerRunError creates a new error with the code "run_error".
func NewServerRunError() Error {
	return NewServerRunErrorBuilder().Build()
}
