package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/nsf/jsondiff"
	"github.com/spf13/cobra"
)

var (
	inputFilePath string
	useMsi        bool
)

var matchCmd = &cobra.Command{
	Use:   "match",
	Short: "Matching check the difference between desired and actual",
	Long: `Matching check difference between desired and actual.

You can verify if there is a difference between your desired properties
and actual with this CLI. This CLI read your desired properties as JSON files,
and query to Azure Resource Graph API, then check the difference.`,
	Run: func(cmd *cobra.Command, args []string) {
		if inputFilePath != "" {
			log.Printf("[INFO] input file path: %s\n", inputFilePath)
		} else {
			fmt.Fprintln(os.Stderr, "[ERROR] input file path with --file/-f")
			os.Exit(1)
		}
		filenames, err := filepath.Glob(inputFilePath)
		checkErr(err)
		if len(filenames) == 0 {
			fmt.Fprintf(os.Stderr, "[ERROR] no files matched the file path/pattern: %s\n", inputFilePath)
			os.Exit(1)
		}
		log.Printf("[INFO] input filenames: %v\n", filenames)

		log.Printf("[INFO] use-msi: %v\n", useMsi)
		log.Printf("[DEBUG] client config: %#v\n", config)

		b := &authentication.Builder{
			SubscriptionID:                 config.Subscription.ID,
			ClientID:                       config.Client.ID,
			ClientSecret:                   config.Client.Secret,
			TenantID:                       config.Tenant.ID,
			Environment:                    config.Environment,
			ClientCertPath:                 config.Cert.Path,
			ClientCertPassword:             config.Cert.Password,
			SupportsManagedServiceIdentity: useMsi,
		}

		c, err := buildClient(*b)
		checkErr(err)
		log.Printf("[DEBUG] client: %#v\n", c)

		failFlag := false

		for _, v := range filenames {
			rs, err := getDesiredResources(v)
			checkErr(err)

			for _, r := range rs {
				log.Printf("[DEBUG] desired resource: %#v\n", r)
				id := r.(map[string]interface{})["id"].(string)
				log.Printf("[INFO] target resource id: %#v\n", id)
				idf := strings.Split(strings.Trim(id, "/"), "/")
				sub := idf[1]
				log.Printf("[DEBUG] target subscription id: %#v\n", sub)

				query := fmt.Sprintf("where id == '%s'", id)
				res, err := rgQuery(c, sub, query)
				checkErr(err)
				log.Printf("[DEBUG] query result: %#v\n", res)
				log.Printf("[DEBUG] query result (Data): %#v\n", res.Data.([]interface{})[0])

				switch *res.TotalRecords {
				case int64(0):
					log.Printf("[ERROR] not found the resource. id: %s\n", id)
					failFlag = true
					continue
				case int64(1):
					log.Printf("[DEBUG] found the resource. id: %s\n", id)
				default:
					log.Printf("[ERROR] the resource is not unique. id: %s\n", id)
					failFlag = true
					continue
				}

				act, err := json.Marshal(&res.Data.([]interface{})[0])
				checkErr(err)
				des, err := json.Marshal(&r)
				checkErr(err)
				d, s := diff(act, des)
				log.Printf("[DEBUG] difference: %#v\n", d)
				if d == "FullMatch" || d == "SupersetMatch" {
					log.Printf("[INFO] match. id: %s\n", id)
				} else {
					log.Printf("[ERROR] unmatch. id: %s\n", id)
					log.Printf("[ERROR] difference details: %s\n", s)
					failFlag = true
				}

				qr := res.Header.Get("x-ms-user-quota-remaining")
				log.Printf("[DEBUG] api quota remaining: %s\n", qr)
				ra := res.Header.Get("x-ms-user-quota-resets-after")
				log.Printf("[DEBUG] api quota reset after: %s\n", ra)
				q, err := strconv.Atoi(qr)
				checkErr(err)
				td1, err := time.Parse("15:04:05", ra)
				checkErr(err)
				td2, err := time.Parse("15:04:05", "00:00:00")
				checkErr(err)
				if q == 0 {
					time.Sleep(td1.Sub(td2))
				}
			}
		}

		if failFlag {
			log.Println("[ERROR] matching(s) was not successful")
			os.Exit(1)
		}
		log.Println("[INFO] matching(s) was successful")
	},
}

func init() {
	matchCmd.PersistentFlags().StringVarP(&inputFilePath, "file", "f", "", "path(glob) of the file(s) where the desired resources are written")
	matchCmd.PersistentFlags().BoolVarP(&useMsi, "use-msi", "", false, "flag for using Managed Identity to auth (defalut: false)")
	rootCmd.AddCommand(matchCmd)
}

func getDesiredResources(filename string) ([]interface{}, error) {
	df, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer df.Close()

	d, err := os.ReadFile(df.Name())
	if err != nil {
		return nil, err
	}

	var j []interface{}
	err = json.Unmarshal(d, &j)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func diff(actual, desired []byte) (string, string) {
	opts := jsondiff.DefaultJSONOptions()
	d, s := jsondiff.Compare(actual, desired, &opts)
	return d.String(), s
}
