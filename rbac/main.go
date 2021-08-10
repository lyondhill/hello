package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	ra "github.com/casbin/redis-adapter/v2"
)

func main() {
	content := `[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act`
	m, err := model.NewModelFromString(content)
	if err != nil {
		panic(err)
	}

	// Initialize a Redis adapter and use it in a Casbin enforcer:
	a := ra.NewAdapter("tcp", "127.0.0.1:6379") // Your Redis network and address.

	// Use the following if Redis has password like "123"
	//a := redisadapter.NewAdapterWithPassword("tcp", "127.0.0.1:6379", "123")
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		panic(err)
	}

	// Load the policy from DB.
	fmt.Println(e.LoadPolicy())

	// Check the permission.
	fmt.Println(e.AddNamedPolicy("lyons", "cow", "domain1", "resource1", "get"))
	fmt.Println(e.Enforce("alice", "data1", "read"))
	fmt.Println(e.Enforce("cow", "data1", "read"))
	fmt.Println(e.Enforce("cow", "domain1", "resource1", "read"))
	fmt.Println("addroll")
	fmt.Println(e.AddRoleForUser("bessy", "cow", "domain1"))
	fmt.Println(e.AddRoleForUser("herferter", "cow", "domain1"))
	fmt.Println(e.AddRoleForUser("glen", "cow", "domain1"))
	fmt.Println(e.Enforce("bessy", "domain1", "resource1", "get"))
	fmt.Println(e.Enforce("herferter", "domain1", "resource2", "get"))

	fmt.Println("actions", e.GetAllActions())
	fmt.Println("subj", e.GetAllSubjects())
	fmt.Println("roles", e.GetAllRoles())

	// fmt.Println(e.GetFilteredNamedPolicy("lyons", 0, "cow"))

	// Modify the policy.
	// e.AddPolicy(...)
	// e.RemovePolicy(...)

	// Save the policy back to DB.
	e.SavePolicy()
}
