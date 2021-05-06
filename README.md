## Usage

The package is comprised of services used to manipulate various
adwords structures.  To access a service you need to create an
gads.Auth and parse it to the service initializer, then can call
the service methods on the service object.

~~~ go
     authConf, err := NewCredentialsFromFile("~/creds.json")
     campaignService := gads.NewCampaignService(&authConf.Auth)

     campaigns, totalCount, err := campaignService.Get(
       gads.Selector{
         Fields: []string{
           "Id",
           "Name",
           "Status",
         },
       },
     )
~~~

