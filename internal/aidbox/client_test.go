package aidbox

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestAidboxClient_CreateTokenIntrospector(t *testing.T) {
	t.Skip("Implement this as an integration test")
	client := NewClient("http://localhost:48083", "client-a", "secret")
	ti := &TokenIntrospector{
		JWKSURI:                    "http://keycloak:8080/auth/realms/pkb/protocol/openid-connect/certs",
		TokenIntrospectorJWT:       &TokenIntrospectorJWT{
			ISS:    "http://keycloak:8080/auth/realms/pkb",
		},
		Type:                       TokenIntrospectorTypeJWT,
	}
	b, _ := json.Marshal(ti)
	foo := string(b)
	fmt.Printf("%s\n", foo)
    tokenIntrospector, err := client.CreateTokenIntrospector(ti)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if tokenIntrospector == nil {
        t.Errorf("token introspector is nil")
    }
}

func TestAidboxClient_GetTokenIntrospectors(t *testing.T) {
	t.Skip("Implement this as an integration test")
	client := NewClient("http://localhost:48083", "client-a", "secret")
	_, err := client.GetTokenIntrospectors()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	//if len(tis) != 2 {
	//	t.Errorf("Wrong number of token introspectors found")
	//}
}


func TestAidboxClient_GetTokenIntrospector(t *testing.T) {
	t.Skip("Implement this as an integration test")
	client := NewClient("http://localhost:48083", "client-a", "secret")
	ti, err := client.GetTokenIntrospector("216215c0-af77-4949-ab50-eb85f8304577")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ti.ID != "216215c0-af77-4949-ab50-eb85f8304577" {
		t.Errorf("Wrong token introspector found")
	}
	//if len(tis) != 2 {
	//	t.Errorf("Wrong number of token introspectors found")
	//}
}