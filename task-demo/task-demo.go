package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/influxdata/platform"
	phttp "github.com/influxdata/platform/http"
)

var (
	querydEndpoint  = flag.String("queryd", "http://localhost:8093", "HTTP endpoint of queryd server")
	gatewayEndpoint = flag.String("gatewayd", "http://localhost:9999", "HTTP endpoint of gatewayd server")

	namespace string
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	flag.Usage = func() {
		base := path.Base(os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s namespace [bootstrap|write|create-task|read-in|read-out|downsample-once]\n", base)
		fmt.Fprintf(flag.CommandLine.Output(), "\tbootstrap: create org, user, buckets using the given namespace\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\twrite: write to the input bucket forever\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\tcreate-task: create a downsample task, using the buckets from the namespace\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\tread-in: read the past 5s from the input bucket\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\tread-out: read the past 5s from the output bucket\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\tdownsample-once: downsample the past 5s from the input bucket, once\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\tlist-tasks: list tasks for the namespace's user\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\tremove-tasks: remove tasks for the namespace's user\n")
		flag.PrintDefaults()

		fmt.Fprintf(flag.CommandLine.Output(), "Typical workflow:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\t1. run `bootstrap`\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\t2. run `write` in another window and leave it running\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\t3. show off `downsample-once` and `read-out`\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\t4. bring up taskd logs in another window\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\t5. show off `create-task` and `list-tasks`\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\t6. run `read-out` some more\n")
	}
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	namespace = flag.Arg(0)

	switch flag.Arg(1) {
	case "bootstrap":
		bootstrap()
	case "write":
		write()
	case "read-in":
		readOnce(bucketInName(), "-5s")
	case "read-out":
		readOnce(bucketOutName(), "-30s")
	case "downsample-once":
		downsampleOnce("-5s")
	case "create-task":
		createTask()
	case "list-tasks":
		listTasks()
	case "remove-tasks":
		removeTasks()
	default:
		flag.Usage()
		os.Exit(1)
	}
	os.Exit(0)

}

func userName() string {
	return "demo-user-" + namespace
}
func orgName() string {
	return "demo-org-" + namespace
}
func bucketInName() string {
	return "demo-bucket-in-" + namespace
}
func bucketOutName() string {
	return "demo-bucket-out-" + namespace
}

func bootstrap() {
	ctx := context.Background()

	users := phttp.UserService{Addr: *gatewayEndpoint}
	u := &platform.User{Name: userName()}
	if err := users.CreateUser(ctx, u); err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	log.Printf("Created user %q with ID %x", u.Name, []byte(u.ID))

	orgs := phttp.OrganizationService{Addr: *gatewayEndpoint}
	o := &platform.Organization{Name: orgName()}
	if err := orgs.CreateOrganization(ctx, o); err != nil {
		log.Fatalf("failed to create org: %v", err)
	}
	log.Printf("Created org %q with ID %x", o.Name, []byte(o.ID))

	buckets := phttp.BucketService{Addr: *gatewayEndpoint}
	bIn := &platform.Bucket{Name: bucketInName(), OrganizationID: o.ID, RetentionPeriod: time.Hour}
	if err := buckets.CreateBucket(ctx, bIn); err != nil {
		log.Fatalf("failed to create bucket: %v", err)
	}
	log.Printf("Created bucket %q with ID %x", bIn.Name, []byte(bIn.ID))
	bOut := &platform.Bucket{Name: bucketOutName(), OrganizationID: o.ID, RetentionPeriod: 24 * time.Hour}
	if err := buckets.CreateBucket(ctx, bOut); err != nil {
		log.Fatalf("failed to create bucket: %v", err)
	}
	log.Printf("Created bucket %q with ID %x", bOut.Name, []byte(bOut.ID))

	log.Printf("Try one of these queries:")
	log.Printf(`curl -v -XPOST localhost:8093/v1/query --data-urlencode orgID=%s --data-urlencode 'q=from(bucket:"%s") |> range(start:-5m)'`,
		o.ID.String(), bIn.Name,
	)
	log.Printf(`curl -v -XPOST localhost:8093/v1/query --data-urlencode orgID=%s --data-urlencode 'q=from(bucket:"%s") |> range(start:-5m)'`,
		o.ID.String(), bOut.Name,
	)
}

func write() {
	users := phttp.UserService{Addr: *gatewayEndpoint}
	un := userName()
	u, err := users.FindUser(context.Background(), platform.UserFilter{Name: &un})
	if err != nil {
		log.Fatal(err)
	}

	buckets := phttp.BucketService{Addr: *gatewayEndpoint}
	bn := bucketInName()
	on := orgName()
	bIn, err := buckets.FindBucket(context.Background(), platform.BucketFilter{Name: &bn, Organization: &on})
	if err != nil {
		log.Fatal(err)
	}

	auths := phttp.AuthorizationService{Addr: *gatewayEndpoint}
	a := &platform.Authorization{
		UserID: u.ID,
		Permissions: []platform.Permission{
			platform.WriteBucketPermission(bIn.ID),
		},
	}
	if err := auths.CreateAuthorization(context.Background(), a); err != nil {
		log.Fatalf("failed to create authorization: %v", err)
	}
	log.Printf("Created authorization: %#v", a)

	url := *gatewayEndpoint + "/v1/write?org=" + on + "&bucket=" + bn

	for i := 0; ; i++ {
		if i != 0 {
			time.Sleep(100 * time.Millisecond)
		}
		write := fmt.Sprintf("counter n=%d", i)
		req, err := http.NewRequest("POST", url, strings.NewReader(write))
		if err != nil {
			log.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("User-Agent", "demo")
		req.Header.Set("Authorization", "Token "+a.Token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("failed to write batch: %v", err)
		}
		if resp.StatusCode != 204 {
			log.Fatalf("unexpected response status code from write: %d", resp.StatusCode)
		}

		log.Printf("successfully wrote %q to bucket %q in org %q", write, bn, on)
	}
}

func readOnce(bucketName, startRange string) {
	orgs := phttp.OrganizationService{Addr: *gatewayEndpoint}
	on := orgName()
	o, err := orgs.FindOrganization(context.Background(), platform.OrganizationFilter{Name: &on})
	if err != nil {
		log.Fatal(err)
	}

	q := fmt.Sprintf("from(bucket:%q) |> range(start:%s)", bucketName, startRange)
	vals := url.Values{
		"orgID": []string{o.ID.String()},
		"q":     []string{q},
	}

	req, err := http.NewRequest("POST", *querydEndpoint+"/v1/query", strings.NewReader(vals.Encode()))
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "demo")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed on HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read HTTP response: %v", err)
	}
	log.Printf("Executed query: %s", q)
	log.Printf("Response code: %d", resp.StatusCode)
	log.Printf("Response: %s", string(body))
}

func downsampleOnce(startRange string) {
	orgs := phttp.OrganizationService{Addr: *gatewayEndpoint}
	on := orgName()
	o, err := orgs.FindOrganization(context.Background(), platform.OrganizationFilter{Name: &on})
	if err != nil {
		log.Fatal(err)
	}

	q := fmt.Sprintf(
		"from(bucket:%q) |> range(start:%s) |> sum() |> to(bucket:%q, org:%q)",
		bucketInName(), startRange, bucketOutName(), on,
	)
	vals := url.Values{
		"orgID": []string{o.ID.String()},
		"q":     []string{q},
	}

	req, err := http.NewRequest("POST", *querydEndpoint+"/v1/query", strings.NewReader(vals.Encode()))
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "demo")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed on HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read HTTP response: %v", err)
	}
	log.Printf("Executed query: %s", q)
	log.Printf("Response code: %d", resp.StatusCode)
	log.Printf("Response: %s", string(body))
}

func createTask() {
	orgs := phttp.OrganizationService{Addr: *gatewayEndpoint}
	on := orgName()
	o, err := orgs.FindOrganization(context.Background(), platform.OrganizationFilter{Name: &on})
	if err != nil {
		log.Fatal(err)
	}

	users := phttp.UserService{Addr: *gatewayEndpoint}
	un := userName()
	u, err := users.FindUser(context.Background(), platform.UserFilter{Name: &un})
	if err != nil {
		log.Fatal(err)
	}

	var taskJSON struct {
		ID           platform.ID `json:"id,omitempty"`
		Organization platform.ID `json:"organizationId"`
		Name         string      `json:"name"`
		// Status string `json:"status"`
		Owner platform.User `json:"owner"`
		Flux  string        `json:"flux"`
		Every string        `json:"every,omitempty"`
		Cron  string        `json:"cron,omitempty"`
		// Last   Run    `json:"last,omitempty"`
	}

	taskJSON.Name = "downsample-every-1s"
	taskJSON.Organization = o.ID
	taskJSON.Owner = platform.User{
		ID:   u.ID,
		Name: u.Name,
	}
	taskJSON.Flux = fmt.Sprintf(
		`//! every=1s
		from(bucket:%q) |> range(start:%s) |> sum() |> to(bucket:%q, org:%q)`,
		bucketInName(), "-5s", bucketOutName(), on,
	)

	b, err := json.Marshal(taskJSON)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Creating task with JSON: %s", b)

	req, err := http.NewRequest("POST", *gatewayEndpoint+"/v1/tasks", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "demo")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed on HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Task creation status code: %d", resp.StatusCode)
	log.Printf("Task creation headers: %v", resp.Header)
	log.Printf("Response from task creation: %s", string(body))
}

func listTasks() {
	users := phttp.UserService{Addr: *gatewayEndpoint}
	un := userName()
	u, err := users.FindUser(context.Background(), platform.UserFilter{Name: &un})
	if err != nil {
		log.Fatal(err)
	}

	url := *gatewayEndpoint + "/v1/tasks?user=" + u.ID.String()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var tasks []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Script string `json:"flux"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		log.Fatal(err)
	}
	for _, t := range tasks {
		log.Printf("Task: ID=%s Name=%s Script=%s", t.ID, t.Name, t.Script)
	}
}

func removeTasks() {
	fmt.Println("here")
	users := phttp.UserService{Addr: *gatewayEndpoint}
	un := userName()
	u, err := users.FindUser(context.Background(), platform.UserFilter{Name: &un})
	if err != nil {
		fmt.Println("hereerr", err)
		log.Fatal(err)
	}

	fmt.Println("user", u)

	url := *gatewayEndpoint + "/v1/tasks?user=" + u.ID.String()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("hereear", err)
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var tasks []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Script string `json:"flux"`
	}
	fmt.Println("here")
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		log.Fatal(err)
	}
	fmt.Println("here2")

	for _, t := range tasks {
		log.Printf("Task: ID=%s\n", t.ID)

		req, err := http.NewRequest("DELETE", *gatewayEndpoint+"/v1/tasks/"+t.ID, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("User-Agent", "demo")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("Failed on HTTP request: %v", err)
		}
		defer resp.Body.Close()

		log.Printf("Task removed status code: %d", resp.StatusCode)

	}
}
