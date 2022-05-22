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
	"io"
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
	Use:   "opactl [rule]... [flags] (ex. opactl rule1 subrule11 ...)",
	Short: "opactl executes your own Rego (OPA) policy as CLI command.",
	Long: `You define a rule in OPA policy, for example "rule1". 
Then, "opactl" detects your rule and turns it into subcommand such as "opactl rule1".`,

	//Args: cobra.MinimumNArgs(1),
	Args: func(cmd *cobra.Command, args []string) error {
    if (len(args) < 1) && !viper.GetBool("all") {
      return errors.New("requires at least one rule (=command) or --all flag")
    }

		return nil
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		initConfig()

		_, policyPath, lastPath := Parse(args)
		strList, _ := GetAllRules(policyPath, lastPath, true, cmd.OutOrStderr())
		if toComplete == "." {
			strList = append(strList, ".", "..")
		}
		strList = append(strList, "help\tDisplay help contents (description of rule)")

		return strList, cobra.ShellCompDirectiveNoFileComp
	},
  Run: func(cmd *cobra.Command, args []string) {
		// Prepare commands
		commands := args[:]
		allFlag := viper.GetBool("all")

		err := execOpa(cmd, commands, allFlag, false, viper.GetString("query"))
		if err != nil {
			log.Fatal(err)
		}
  },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is <current directory>/.opactl)")
	rootCmd.PersistentFlags().StringSliceP("parameter", "p", []string{}, "parameter (key=value)")
	viper.BindPFlag("parameter", rootCmd.PersistentFlags().Lookup("parameter"))
	rootCmd.PersistentFlags().StringSliceP("directory", "d", []string{}, "directories")
	viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))
	rootCmd.PersistentFlags().StringSliceP("bundle", "b", []string{}, "bundles")
	viper.BindPFlag("bundle", rootCmd.PersistentFlags().Lookup("bundle"))

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Toggle verbose mode on/off (display print() output)")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("all", "a", false, "Show all commands")
	viper.BindPFlag("all", rootCmd.Flags().Lookup("all"))
	rootCmd.Flags().BoolP("input", "i", false, "Accept stdin as input.stdin. Multiple lines are stored as array. JSON will be parsed and stored in input.json_stdin as well.")
	viper.BindPFlag("input", rootCmd.Flags().Lookup("input"))

	rootCmd.Flags().StringP("query", "q", "", "Input your own query script (example: { rtn | rtn := 1 }")
	viper.BindPFlag("query", rootCmd.Flags().Lookup("query"))
	rootCmd.Flags().StringP("base", "B", "data.opactl", "OPA base path which will be evaluated")
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
			arrayString = append(arrayString, scanner.Text())
		}
		rtn["stdin"] = arrayString

		allString := strings.Join(arrayString, "")

		var record interface{}
		err := json.Unmarshal([]byte(allString), &record)
		if err != nil  {
			record = map[string]interface{}{ "error": "Failed JSON parsing."}
		}

		rtn["json_stdin"] = record
	}
	return rtn
}


func opaEval(policyPath string, q string, verbose bool, stdout, stderr io.Writer) error {
		// Prepare directories
		directories = viper.GetStringSlice("directory")
		printVerbose(verbose, "directory setting is", viper.GetStringSlice("directory")...)
		
		opts := []string{"eval"}
	
		bundles := viper.GetStringSlice("bundle")
		for _, b := range bundles {
			opts = append(opts, "-b")
			opts = append(opts, b)
		}
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
	
		cmdExec.Stdout = stdout
		cmdExec.Stderr = stderr
		err = cmdExec.Run()
		return err
}

func execOpa(cmd *cobra.Command, commands []string, allFlag bool, desc bool, query string) error {
	out := cmd.OutOrStdout()
	stderr := cmd.OutOrStderr()

	verbose := viper.GetBool("verbose")

	abstractPathArray, policyPath, lastPath := Parse(commands)

	help := (lastPath == "help")

	if help {
		needSubrule := false
		if len(commands) >= 1 && len(abstractPathArray) >= 3 {
			lastPath = abstractPathArray[len(abstractPathArray)-2]
			basePath := abstractPathArray[:len(abstractPathArray)-2]
			policyPath2 := strings.Join(basePath, ".")
			policyPath = strings.Join(abstractPathArray[:len(abstractPathArray)-1], ".")

			_Type, err := GetType(policyPath, lastPath, stderr)
			if err != nil {
				return err
			}
			
			if _Type == "object" {
				needSubrule = true
				cmd.Use = fmt.Sprintf(`opactl %s [subrule]...`, strings.Join(commands[:len(commands)-1], " "))
			}

			help, err2 := GetComment(policyPath2, lastPath, stderr)
			if err2 != nil {
				return err2
			}
			cmd.Long = fmt.Sprintf("%s: %s", lastPath, help)
		}

		subCompletionCmd, _, _ := cmd.Find([]string{"completion"})
		cmd.RemoveCommand(subCompletionCmd)
		cmd.Help()

		if !needSubrule {
			return nil
		}
		fmt.Fprintf(out, `
Subrules:
`)
		subcommands, _ := GetAllRules(policyPath, lastPath, true, stderr)
		for _, v := range subcommands {
			fmt.Fprintf(out, "  %s\n", v)
		}
		return nil
	}

	q := query
	if q == "" {
		q = lastPath
	}
	if allFlag {
		q = getAllQuery(lastPath, desc)
	}

	return opaEval(policyPath, q, verbose, out, stderr)
}

func Parse(commands []string) ([]string, string, string) {
	c := commands

	// Prepare paths and query
	basePath := viper.GetString("base")
	if basePath != "" {
		c = append(strings.Split(basePath, "."), c...)
	}
	abstractPathArray := getAbstractPathArray(c)
	policyPath := strings.Join(abstractPathArray, ".")
	lastPath := abstractPathArray[len(abstractPathArray)-1]
	return abstractPathArray, policyPath, lastPath
}

func GetAllRules(policyPath, lastPath string, desc bool, stderr io.Writer) ([]string, error) {
	var outString []string
	err := GetObject(policyPath, getAllQuery(lastPath, desc), stderr, &outString)
	return outString, err
}

func GetType(policyPath string, lastPath string, stderr io.Writer) (string, error) {
	var outString string
	err := GetObject(policyPath, getTypeQuery(lastPath), stderr, &outString)
	return outString, err
}

func GetComment(policyPath, lastPath string, stderr io.Writer) (string, error) {
	var outString string
	err := GetObject(policyPath, getCommentQuery(policyPath, lastPath), stderr, &outString)
	return outString, err
}

func GetObject(policyPath string, query string, stderr io.Writer, outString interface{}) (error) {
	var err error = nil
	outBytes := bytes.Buffer{}

	if err = opaEval(policyPath, query, false, &outBytes, stderr); err != nil {
		return err
	}

	_ = json.Unmarshal(outBytes.Bytes(), outString)
	return nil
}

func getAllQuery(lastPath string, desc bool) string {
	var q string
	if desc {
		queryTemplate := `
		{	desc |
			%s[key]
			not startswith(key, "__")
			field_comment_path := concat("", ["__", key])
			field_comment := object.get(%s, field_comment_path, "")
			package_comment_path := [key, "__comment"]
			package_comment := object.get(%s, package_comment_path, "")
			desc := concat("", [key, "\t", field_comment, package_comment])
		}
		`
		q = fmt.Sprintf(queryTemplate, lastPath, lastPath, lastPath)
	} else {
		queryTemplate := `
		{	key |
			%s[key]
			not startswith(key, "__")
		}
		`
		q = fmt.Sprintf(queryTemplate, lastPath)

	}
	return q
}

func getTypeQuery(lastPath string) string {
	return fmt.Sprintf(`
		type_name(%s)
	`, lastPath)
}

func getCommentQuery(policyPath, lastPath string) string {
	queryTemplate := `{ desc |
		field_comment_path := "__%s"
		field_comment := object.get(%s, field_comment_path, "")
		package_comment_path := ["%s", "__comment"]
		package_comment := object.get(%s, package_comment_path, "")
		desc := concat("", [field_comment, package_comment])
	}[_]`

	return fmt.Sprintf(queryTemplate, lastPath, policyPath, lastPath, policyPath)
}

func printVerbose(verbose bool, label string, log ...string) {
	if verbose {
		fmt.Println(label, log)
	}
}

func getAbstractPathArray(words []string) []string {
	rtn := []string{}
	for _, word := range words {
		// fmt.Println("word", word)
		if word == ".." {
			if len(rtn) > 1 {
				rtn = rtn[:len(rtn)-1]
			}
		} else if word == "." {
		} else {
			if strings.Contains(word, ".") {
				log.Fatal(`[ERROR] Arguments cannot contain "." except "." and ".."`)
			}
			rtn = append(rtn, word)
		}
	}
	return rtn
}
