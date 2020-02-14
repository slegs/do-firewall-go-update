package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"reflect"

	"gopkg.in/digitalocean/godo.v1"

	"log"
)

type  Firewall struct {
	ID             string
	Name           string
  InboundRules   []godo.InboundRule
  OutboundRules  []godo.OutboundRule
}



func loadIps(filename string) (rules Firewall, err error) {
	var b []byte

	f, ok := os.Open(filename)

	defer f.Close()

	if ok == nil {

		b, err = ioutil.ReadAll(f)

		if err == nil {
			err = json.Unmarshal(b, &rules)

		}
	}

	return rules,err

}

func saveIps(filename string, rules Firewall) (err error) {

	f, err := os.Create(filename)
	if err != nil {
		return
	}

	defer f.Close()

	b, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return
	}

	f.Write(b)

	return
}

func main() {
	var apiToken, firewallName, firewallID, oldipsfile, newipsfile string

	flag.StringVar(&apiToken, "token", "", "The digitalocean api token")
	flag.StringVar(&firewallName, "firewall-name", "", "The name of the firewall")
	flag.StringVar(&firewallID, "firewall-id", "", "The id of the firewall")
	flag.StringVar(&oldipsfile, "old-ips", "old_ips.json", "Old rules file")
	flag.StringVar(&newipsfile, "new-ips", "new_ips.json", "New rules file")

	flag.Parse()

	if apiToken == "" {
		log.Fatal("You must specify the --token")
	}

	if firewallName == "" && firewallID == "" {
		log.Fatal("You must specify the --firewall-name or the --firewall-id")
	}

	oldIps,err := loadIps(oldipsfile)

	if err != nil {
		log.Fatal(err)
	}

	newIps,err := loadIps(newipsfile)

	if err != nil {
		log.Fatal(err)
	}


	newinstallation := false

	if len(newIps.Name) == 0 && len(oldIps.Name) == 0 {

		newinstallation = true

		log.Println("New installation - downloading rules")

	} else {

		if reflect.DeepEqual(oldIps, newIps) {

			log.Println("NO CHANGE")

			return

		}
	}

	client := newClient(apiToken)

	var firewall *godo.Firewall
	if firewallID != "" {
		firewall, err = findFirewallByID(client, firewallID)
	} else {
		firewall, err = findFirewallByName(client, firewallName)
	}

	if err != nil {
		log.Fatal(err)
	}

	// Small hack to fix inconsistencies in the digitalocean api
	firewall.InboundRules, firewall.OutboundRules = fixInboundOutboundRules(firewall.InboundRules, firewall.OutboundRules)

	if newinstallation {

		newIps = Firewall{firewall.ID,firewall.Name,firewall.InboundRules,firewall.OutboundRules}

		err = saveIps(newipsfile,newIps)
		if err != nil {
			log.Fatal(err)
		}

	} else {


		//var newrule bool
		if newIps.Name == firewall.Name && newIps.ID == firewall.ID {
			log.Println("Firewall Found: " + newIps.Name + " / " + newIps.ID)

			fr := &godo.FirewallRequest{
	 	 	Name:          firewall.Name,
	 	 	InboundRules:  newIps.InboundRules,
	 	 	OutboundRules: newIps.OutboundRules,
	 	 	DropletIDs:    firewall.DropletIDs,
	 	 	Tags:          firewall.Tags,
	 	 }

	 	 err = updateFirewall(client, firewall.ID, fr)

	 	 if err != nil {
	 	 	log.Fatal(err)
	 	 }

	 		log.Println("Firewall Name: " + firewall.Name + " ID: " + firewall.ID + " updated successfully")

		}



	}

	err = saveIps(oldipsfile,newIps)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Rules file updated successfully")
}
