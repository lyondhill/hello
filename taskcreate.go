package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/influxdata/idpe/task/store"
	"github.com/influxdata/idpe/task/store/rpc"
	"github.com/influxdata/platform"
	phttp "github.com/influxdata/platform/http"
	"google.golang.org/grpc"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	taskdEndpoint   = flag.String("taskd", "localhost:8275", "TCP endpoint of taskd server")
	gatewayEndpoint = flag.String("gatewayd", "http://localhost:9999", "HTTP endpoint of gatewayd server")
)

func main() {
	flag.Parse()

	now := fmt.Sprintf("%d", time.Now().UTC().Unix())

	ctx := context.Background()

	// Connect to taskd right away, so we can fail fast if anything goes wrong.
	gc, err := grpc.Dial(*taskdEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial taskd grpc server: %v", err)
	}
	taskStore := rpc.NewGRPCStore(gc)

	users := phttp.UserService{Addr: *gatewayEndpoint}
	u := &platform.User{Name: "task-demo-user-" + now}
	if err := users.CreateUser(ctx, u); err != nil {
		log.Fatalf("failed to create user: %#v", err)
	}
	log.Printf("Created user %q with ID %x", u.Name, []byte(u.ID))

	orgs := phttp.OrganizationService{Addr: *gatewayEndpoint}
	o := &platform.Organization{Name: "task-demo-org-" + now}
	if err := orgs.CreateOrganization(ctx, o); err != nil {
		log.Fatalf("failed to create org: %#v", err)
	}
	log.Printf("Created org %q with ID %x", o.Name, []byte(o.ID))

	buckets := phttp.BucketService{Addr: *gatewayEndpoint}
	bIn := &platform.Bucket{Name: "task-demo-bucket-in-" + now, OrganizationID: o.ID, RetentionPeriod: time.Hour}
	if err := buckets.CreateBucket(ctx, bIn); err != nil {
		log.Fatalf("failed to create bucket: %#v", err)
	}
	log.Printf("Created bucket %q with ID %x", bIn.Name, []byte(bIn.ID))
	bOut := &platform.Bucket{Name: "task-demo-bucket-out-" + now, OrganizationID: o.ID, RetentionPeriod: 24 * time.Hour}
	if err := buckets.CreateBucket(ctx, bOut); err != nil {
		log.Fatalf("failed to create bucket: %#v", err)
	}
	log.Printf("Created bucket %q with ID %x", bOut.Name, []byte(bOut.ID))

	auths := phttp.AuthorizationService{Addr: *gatewayEndpoint}
	a := &platform.Authorization{
		UserID: u.ID,
		Permissions: []platform.Permission{
			platform.WriteBucketPermission(bIn.ID),
		},
	}
	if err := auths.CreateAuthorization(ctx, a); err != nil {
		log.Fatalf("failed to create authorization: %#v", err)
	}
	log.Printf("Created authorization: %#v", a)

	go writePoints(o.Name, bIn.Name, a.Token)
	log.Printf("Writing 10 points every second...")

	log.Printf("Try one of these queries:")
	log.Printf(`curl -v -XPOST localhost:8093/v1/query --data-urlencode orgID=%s --data-urlencode 'q=from(bucket:"%s") |> range(start:-3s)'`,
		o.ID.String(), bIn.Name,
	)
	log.Printf(`curl -v -XPOST localhost:8093/v1/query --data-urlencode orgID=%s --data-urlencode 'q=from(bucket:"%s") |> range(start:-3s)'`,
		o.ID.String(), bOut.Name,
	)

	if err := createTasks(taskStore, o.ID, u.ID, bIn, bOut); err != nil {
		log.Fatalf("error creating tasks: %#v", err)
	}

	// Block forever.
	select {}
}

func writePoints(orgName, bucketName, token string) {
	url := *gatewayEndpoint + "/v1/write?org=" + orgName + "&bucket=" + bucketName

	// Once a second, write 10 random values.
	for t := range time.Tick(time.Second) {
		now := t.UnixNano()

		var lines [10]string
		for i := range lines {
			// Writing the ten points in time-descending order because it's easier.
			timestamp := now - 1e6*int64(i) // Subtract 100ms * i.
			lines[i] = fmt.Sprintf("random_values x=%f %d", rand.Float32(), timestamp)
		}

		req, err := http.NewRequest("POST", url, strings.NewReader(strings.Join(lines[:], "\n")))
		if err != nil {
			log.Fatalf("failed to create request: %#v", err)
		}
		req.Header.Set("User-Agent", "demo/taskwiring")
		req.Header.Set("Authorization", "Token "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("failed to write batch: %#v (%#v)", err, err.Error())
		}
		if resp.StatusCode != 204 {
			log.Fatalf("unexpected response status code from write: %d", resp.StatusCode)
		}
	}
}

func createTasks(taskStore store.Store, orgID, userID platform.ID, in, out *platform.Bucket) error {
	taskFmt := `//! every=%s
	from(bucket:%q) |> range(start:%s) |> mean() |> to(bucket:%q, org: %q)`

	ctx := context.Background()

	rollup10s := fmt.Sprintf(taskFmt, "1s", in.Name, "10s")
	tid, err := taskStore.CreateTask(ctx, orgID, userID, "10s rollup", rollup10s)
	if err != nil {
		log.Fatalf("failed to create 10s rollup task: %#v", err)
	}
	log.Printf("created 10s rollup task %q with ID %x", rollup10s, []byte(tid))

	rollup1m := fmt.Sprintf(taskFmt, "10s", in.Name, "1m")
	tid, err = taskStore.CreateTask(ctx, orgID, userID, "1m rollup", rollup1m)
	if err != nil {
		log.Fatalf("failed to create 1m rollup task: %#v", err)
	}
	log.Printf("created 1m rollup task %q with ID %x", rollup1m, []byte(tid))

	return nil
}
