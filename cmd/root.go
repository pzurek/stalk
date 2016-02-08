// Copyright © 2016 Piotr Zurek <p.zurek@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/pzurek/clearbit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	clearbitKey string
	cb          *clearbit.Client
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "stalk",
	Short: "A little command line stalker using the Clearbit API",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		email, _ := cmd.Flags().GetString("email")
		stalk(email)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stalk/config.yaml)")
	RootCmd.PersistentFlags().StringVar(&clearbitKey, "key", "", "ClearBit API key")

	RootCmd.Flags().StringP("email", "e", "alex@clearbit.com", "Email of the person to find")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("config")       // name of config file (without extension)
	viper.AddConfigPath("$HOME/.stalk") // adding home directory as first search path
	viper.AddConfigPath(".")            // optionally look for config in the working directory
	viper.AutomaticEnv()                // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	clearbitKey = viper.GetString("clearbit_key")
}

func stalk(email string) {
	cb := clearbit.NewClient(clearbitKey, nil)

	enrichment, err := cb.Enrichements.GetCombined(email)
	if err != nil {
		log.Printf("Getting an enrichment failed: %s\n", err)
		return
	}

	if enrichment.Person == nil {
		fmt.Printf("Didn't find a person associated with: %s\n", email)
		return
	}

	person := enrichment.Person
	printName(person)
	printEmployment(person)
	printLinks(person)
}

func printName(p *clearbit.Person) {
	if p.Name.FullName != nil {
		fmt.Printf("This email seems to belong to: %s\n", *p.Name.FullName)
	}
}

func printEmployment(p *clearbit.Person) {
	if p.Employment.Name != nil {
		if p.Employment.Title != nil {
			fmt.Printf("Looks like they are working at %s as a %s\n", *p.Employment.Name, *p.Employment.Title)
		} else {
			fmt.Printf("Looks like they are working at %s\n", *p.Employment.Name)
		}
	}
}

func printLinks(p *clearbit.Person) {
	links := map[string]string{}

	if p.Facebook.Handle != nil {
		links["facebook"] = fmt.Sprintf("Facebook: https://facebook.com/%s", *p.Facebook.Handle)
	}
	if p.Twitter.Handle != nil {
		links["twitter"] = fmt.Sprintf("Twitter:  https://twitter.com/%s", *p.Twitter.Handle)
	}
	if p.Github.Handle != nil {
		links["github"] = fmt.Sprintf("GitHub:   https://github.com/%s", *p.Github.Handle)
	}
	if p.Linkedin.Handle != nil {
		links["linkedin"] = fmt.Sprintf("LinkedIn: https://linkedin.com/%s", *p.Linkedin.Handle)
	}
	if p.Googleplus.Handle != nil {
		links["googleplus"] = fmt.Sprintf("Google+:  https://plus.google.com/%s", *p.Googleplus.Handle)
	}

	if len(links) == 0 {
		fmt.Println("No public links found.")
		return
	}

	fmt.Println("You can follow them at:")
	for _, v := range links {
		fmt.Println(v)
	}
}
