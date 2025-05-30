/**
 * Copyright 2025 Confluent Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Examples of using bearer authentication with schema registry
package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
)

var srURL = "https://psrc-1234.us-east-1.aws.confluent.cloud"
var tokenURL = "your-token-url"
var clientID = "your-client-id"
var clientSecret = "your-client-secret"
var scopes = []string{"schema_registry"}
var identityPoolID = "pool-1234"
var schemaRegistryLogicalCluster = "lsrc-abcd"

// CustomHeaderProvider is a custom header provider that implements the AuthenticationHeaderProvider interface
type CustomHeaderProvider struct {
	token                        string
	schemaRegistryLogicalCluster string
	identityPoolID               string
}

// GetAuthenticationHeader returns the authentication header for the custom header provider
func (p *CustomHeaderProvider) GetAuthenticationHeader() (string, error) {
	return "Bearer " + p.token, nil
}

// GetLogicalCluster returns the logical cluster for the custom header provider
func (p *CustomHeaderProvider) GetLogicalCluster() (string, error) {
	return p.schemaRegistryLogicalCluster, nil
}

// GetIdentityPoolID returns the identity pool ID for the custom header provider
func (p *CustomHeaderProvider) GetIdentityPoolID() (string, error) {
	return p.identityPoolID, nil
}

func main() {
	// Static token
	staticConf := schemaregistry.NewConfigWithBearerAuthentication(srURL, "token", schemaRegistryLogicalCluster, identityPoolID)
	staticClient, _ := schemaregistry.NewClient(staticConf)

	subjects, err := staticClient.GetAllSubjects()
	if err != nil {
		fmt.Println("Error fetching subjects:", err)
		return
	}
	fmt.Println("Static token subjects:", subjects)

	//OAuthBearer
	ClientCredentialsConf := schemaregistry.NewConfig(srURL)
	ClientCredentialsConf.BearerAuthCredentialsSource = "OAUTHBEARER"
	ClientCredentialsConf.BearerAuthToken = "token"
	ClientCredentialsConf.BearerAuthIdentityPoolID = identityPoolID
	ClientCredentialsConf.BearerAuthLogicalCluster = schemaRegistryLogicalCluster
	ClientCredentialsConf.BearerAuthIssuerEndpointURL = tokenURL
	ClientCredentialsConf.BearerAuthClientID = clientID
	ClientCredentialsConf.BearerAuthClientSecret = clientSecret
	ClientCredentialsConf.BearerAuthScopes = scopes

	ClientCredentialsClient, _ := schemaregistry.NewClient(ClientCredentialsConf)
	subjects, err = ClientCredentialsClient.GetAllSubjects()
	if err != nil {
		fmt.Println("Error fetching subjects:", err)
		return
	}
	fmt.Println("OAuthBearer subjects:", subjects)

	// Custom
	customConf := schemaregistry.NewConfig(srURL)
	customConf.BearerAuthCredentialsSource = "CUSTOM"
	customConf.AuthenticationHeaderProvider = &CustomHeaderProvider{
		token:                        "customToken",
		schemaRegistryLogicalCluster: schemaRegistryLogicalCluster,
		identityPoolID:               identityPoolID,
	}
	schemaRegistryClient, err := schemaregistry.NewClient(customConf)

	if err != nil {
		fmt.Println("Error creating schema registry client:", err)
		return
	}

	subjects, err = schemaRegistryClient.GetAllSubjects()
	if err != nil {
		fmt.Println("Error fetching subjects:", err)
		return
	}
	fmt.Println("Custom OAuth subjects:", subjects)
}
