package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

var groupList []string

func healthz(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok")

func groupz(w http.ResponseWriter, req *http.Request) {
	groupJSON, _ := json.Marshal(groupList)
	fmt.Fprintf(w, string(groupJSON))
}

func main() {
	subID := os.Getenv("SUBSCRIPTION_ID")
	if subID == "" {
		log.Fatalln("Subscription ID empty")
	}

	log.Println("Got subscription. Continuing...")
	groupsClient := resources.NewGroupsClient(subID)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Auth successful...")
	groupsClient.Authorizer = authorizer

	ctx := context.Background()
	log.Println("Getting groups list...")
	groups, err := groupsClient.ListComplete(ctx, "", nil)
	if err != nil {
		log.Println("Error getting groups", err)
	}

	log.Println("Enumerating groups...")
	for groups.NotDone() {
		groupList = append(groupList, *groups.Value().Name)
		log.Println(*groups.Value().Name)
		err := groups.NextWithContext(ctx)
		if err != nil {
			log.Println("error getting next group")
		}
	}

	log.Println("Serving on 8080...")
	http.HandleFunc("/healthz", healtz)
	http.HandleFunc("/groupz", groupz)
	http.ListenAndServe(":8080", nil)

}
