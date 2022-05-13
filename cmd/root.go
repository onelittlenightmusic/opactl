/*
Copyright © 2022 Roy Hiroyuki Osaki <hiroyuki.osaki@hal.hitachi.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var directories []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "opactl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	//Args: cobra.MinimumNArgs(1),
	Args: func(cmd *cobra.Command, args []string) error {
    if (len(args) < 1) && !viper.GetBool("all") {
      return errors.New("requires at least one command or --all flag")
    }
		return nil
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		initConfig()
		// fmt.Println("args", args)
		// fmt.Println("toComplete", toComplete)
		strList := execAll(args)
		// fmt.Println("strList", strList)
		return strList, cobra.ShellCompDirectiveNoFileComp
	},
	// ValidArgs: validArgs,
  Run: func(cmd *cobra.Command, args []string) {
		// Prepare commands
		commands := args[:]
		allFlag := viper.GetBool("all")
		var out bytes.Buffer

		err := execOpa(commands, allFlag, viper.GetString("query"), &out)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(out.String())
  },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.opactl.yaml)")
	rootCmd.PersistentFlags().StringSliceP("parameter", "p", []string{}, "parameter (key=value)")
	viper.BindPFlag("parameter", rootCmd.PersistentFlags().Lookup("parameter"))
		rootCmd.PersistentFlags().StringSliceP("directory", "d", []string{}, "directories")
	viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))

	// rootCmd.PersistentFlags().StringSliceP("precommand", "P", []string{}, "precommand")
	// viper.BindPFlag("precommand", rootCmd.PersistentFlags().Lookup("precommand"))

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Toggle verbose mode on/off")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("all", "a", false, "Show all commands")
	viper.BindPFlag("all", rootCmd.Flags().Lookup("all"))
	rootCmd.Flags().BoolP("input", "i", false, "Accept stdin as input.stdin")
	viper.BindPFlag("input", rootCmd.Flags().Lookup("input"))
	rootCmd.Flags().StringP("query", "q", "", "Input your own query script (example: { rtn | rtn := 1 }")
	viper.BindPFlag("query", rootCmd.Flags().Lookup("query"))
	rootCmd.Flags().StringP("base", "b", "data.opactl", "OPA base path which will be evaluated")
	viper.BindPFlag("base", rootCmd.Flags().Lookup("base"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("opactl")
	viper.AutomaticEnv() 
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".opactl" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".opactl")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		verbose := viper.GetBool("verbose")
		printVerbose(verbose, "Using config file:", viper.ConfigFileUsed())
	}
}

func parseParam(params []string, stdin bool) map[string]interface{} {
	rtn := map[string]interface{}{}
	for _, param := range params {
		kv := strings.Split(param, "=")
		if len(kv) <= 1 {
			continue
		}
		rtn[kv[0]] = kv[1]
	}
	if stdin {
		scanner := bufio.NewScanner(os.Stdin)
		arrayString := []string{}
		for scanner.Scan() {
			// texts := strings.Split(scanner.Text(), " ")
			arrayString = append(arrayString, scanner.Text())
		}
		rtn["stdin"] = arrayString
	}
	return rtn
}

func execOpa(commands []string, allFlag bool, query string, out *bytes.Buffer) error {
	verbose := viper.GetBool("verbose")
	c := commands

	// Prepare paths and query
	basePath := viper.GetString("base")
	if basePath != "" {
		c = append([]string{basePath}, c...)
	}
	policyPath := strings.Join(c, ".")
	lastPath := c[len(c)-1]
	q := query
	if q == "" {
		q = lastPath
	}
	if allFlag {
		q = "{key|"+lastPath+"[key]; not startswith(key, \"__\")}"
	}

	// Prepare directories
	directories = viper.GetStringSlice("directory")
	printVerbose(verbose, "directory setting is", viper.GetStringSlice("directory")...)
	
	opts := []string{"eval"}
	for _, d := range directories {
		opts = append(opts, "-d")
		opts = append(opts, d)
	}
	opts = append(opts, "--import", policyPath, q, "--format=pretty", "-I")

	printVerbose(verbose, "opts:", opts...)
	cmdExec := exec.Command("opa", opts...)
	paramMap := parseParam(viper.GetStringSlice("parameter"), viper.GetBool("input"))
	jsonStr, err := json.Marshal(paramMap)
	cmdExec.Stdin = strings.NewReader(string(jsonStr))
	if err != nil {
		return err
	}

	cmdExec.Stdout = out

	err = cmdExec.Run()

	return err
}

func execAll(commands []string) []string {
	var out bytes.Buffer

	err := execOpa(commands, true, "", &out)
	if err != nil {
		log.Fatal(err)
	}

	var str []string
	// fmt.Println("out:", out.String())
	err = json.Unmarshal(out.Bytes(), &str)

	if err != nil {
		log.Fatal(err)
	}
	return str
}

func printVerbose(verbose bool, label string, log ...string) {
	if verbose {
		fmt.Println(label, log)
	}
}