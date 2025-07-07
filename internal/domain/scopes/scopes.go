package scopes

import (
	"fmt"
	"slices"
	"strings"
)

type Scope struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

var AllScopes = []Scope{
	// Scopes de Usuário
	{Name: "profile:read", Description: "Ler informações públicas do perfil", Category: "Perfil"},
	{Name: "profile:write", Description: "Modificar informações públicas do perfil", Category: "Perfil"},
	{Name: "password:write", Description: "Alterar a própria senha", Category: "Conta"},
	{Name: "account:delete:self", Description: "Apagar a própria conta", Category: "Conta"},

	// Scopes de Notificações
	{Name: "notifications:read", Description: "Ler as notificações do usuário", Category: "Notificações"},
	{Name: "notifications:write", Description: "Gerenciar as notificações (ex: marcar como lida)", Category: "Notificações"},

	// OpenID Connect
	{Name: "openid", Description: "Sinaliza uma requisição de autenticação OpenID Connect", Category: "OpenID"},

	// Scopes Administrativos
	{Name: "users:read", Description: "Visualizar qualquer usuário do sistema", Category: "Admin - Usuários"},
	{Name: "users:create", Description: "Criar novos usuários no sistema", Category: "Admin - Usuários"},
	{Name: "users:update", Description: "Modificar qualquer usuário do sistema", Category: "Admin - Usuários"},
	{Name: "users:delete", Description: "Remover qualquer usuário do sistema", Category: "Admin - Usuários"},
	{Name: "clients:read", Description: "Visualizar clientes OAuth", Category: "Admin - Clientes"},
	{Name: "clients:create", Description: "Criar novos clientes OAuth", Category: "Admin - Clientes"},
	{Name: "clients:update", Description: "Modificar clientes OAuth", Category: "Admin - Clientes"},
	{Name: "clients:delete", Description: "Remover clientes OAuth", Category: "Admin - Clientes"},
}

func GetScopeByName(name string) (Scope, bool) {
	for _, scope := range AllScopes {
		if scope.Name == name {
			return scope, true
		}
	}
	return Scope{}, false
}

func ValidateScopes(scopeNames []string) error {
	for _, scopeName := range scopeNames {
		if _, ok := GetScopeByName(scopeName); !ok {
			return fmt.Errorf("scope inválido: %s", scopeName)
		}
	}
	return nil
}

func HasScope(userScopes []string, requiredScope string) bool {
	return slices.Contains(userScopes, requiredScope)
}

func HasAllScopes(userScopes []string, requiredScopes []string) bool {
	userScopeSet := make(map[string]struct{}, len(userScopes))
	for _, s := range userScopes {
		userScopeSet[s] = struct{}{}
	}

	for _, required := range requiredScopes {
		if _, ok := userScopeSet[required]; !ok {
			return false
		}
	}
	return true
}

func ParseScopes(scopesString string) []string {
	if scopesString == "" {
		return []string{}
	}
	return strings.Fields(scopesString)
}

func JoinScopes(scopes []string) string {
	return strings.Join(scopes, " ")
}

func GetDefaultFirstPartyScopes() []string {
	return []string{
		"openid",
		"profile:read",
	}
}

func GetDefaultThirdPartyScopes() []string {
	return []string{
		"profile:read",
	}
}
