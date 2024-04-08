/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package rest

import (
	"errors"
	"fmt"

	"github.com/cm-mayfly/cm-mayfly/src/cmd"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

var serviceName string
var actionName string
var method string
var isInit bool
var isListMode bool
var isVerbose bool

type ServiceInfo struct {
	BaseURL string `yaml:"baseurl"`
	Auth    struct {
		Type     string `yaml:"type"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"auth"`
}

var serviceInfo ServiceInfo

// apiCmd represents the svc command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Open API calls to the Cloud-Migrator system",
	Long: `Open API calls to the Cloud-Migrator system. For example:

./mayfly api --help
./mayfly api --list
./mayfly api --service spider --list
./mayfly api --service spider --action ListCloudOS
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		isInit = false
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			fmt.Println(cmd.Help())
			return
		}

		//fmt.Println("============ 아규먼트 :  " + strconv.Itoa(len(args)))
		//fmt.Println("============ 플래그 수 :  " + strconv.Itoa(cmd.Flags().NFlag()))

		//viper.AddConfigPath("../conf")
		viper.SetConfigFile(configFile)

		// 설정 파일 읽어오기
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Printf("Error reading config file: %s\n", err)
			return
		}
		isInit = true

		//fmt.Println("cliSpecVersion : ", viper.GetString("cliSpecVersion"))
		//fmt.Println("Loaded configurations:", viper.AllSettings())
		if isVerbose {
			spew.Dump(viper.AllSettings())
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		if !isInit {
			return
		}

		//
		// list 명령어 처리
		//
		if isListMode {
			if isVerbose {
				fmt.Println("List Mode")
			}
			if serviceName != "" {
				showActionList(serviceName)
			} else {
				showServiceList()
			}

			return
		}

		errParse := parseRequestInfo()
		if errParse != nil {
			fmt.Println(errParse)
			return
		}

		/*
					//fmt.Println("ListCloudOS.method:", viper.GetString("serviceActions.spider.ListCloudOS.method"))
			//fmt.Println("ListCloudOS.resourcePath:", viper.GetString("serviceActions.spider.ListCloudOS.resourcePath"))

			// Spider 서비스 정보를 가져오기
			spiderService := viper.GetStringMap("services.spider")

			// Spider 서비스 정보가 존재하는지 확인
			if spiderService != nil {
				// Spider 서비스 정보에서 필요한 값 추출
				baseurl, ok := spiderService["baseurl"].(string)
				if !ok {
					fmt.Println("Error: baseurl is not a string")
					return
				}

				authMap := viper.GetStringMap("services.test.auth")
				// Spider 서비스 정보 출력
				fmt.Println("Spider Service:")
				fmt.Println("Base URL:", baseurl)
				fmt.Println("Auth Type:", authMap)

				if len(authMap) != 0 {
					// auth 맵 내부의 필요한 값 추출
					authType, ok := authMap["type"].(string)
					if !ok {
						fmt.Println("Error: auth type is not a string")
						return
					}
					username, ok := authMap["username"].(string)
					if !ok {
						fmt.Println("Error: username is not a string")
						return
					}
					password, ok := authMap["password"].(string)
					if !ok {
						fmt.Println("Error: password is not a string")
						return
					}

					// 추출한 값 출력
					fmt.Println("Auth Type:", authType)
					fmt.Println("Username:", username)
					fmt.Println("Password:", password)
				}
			} else {
				fmt.Println("Spider service not found")
			}
		*/

		/*
			baseURL := viper.GetString("services.spider.baseurl")
			username := viper.GetString("services.spider.username")
			password := viper.GetString("services.spider.password")

			fmt.Println("Base URL:", baseURL)
			fmt.Println("Username:", username)
			fmt.Println("Password:", password)

			serviceActions := viper.GetStringMap("serviceActions.spider")
			action := serviceActions[actionName]

			if action != nil {
				actionMap := action.(map[string]interface{})
				method := actionMap["method"].(string)
				resourcePath := actionMap["resourcePath"].(string)

				fmt.Println("Method:", method)
				fmt.Println("Resource Path:", resourcePath)
			} else {
				fmt.Println("Action not found for Spider service")
			}
		*/
	},
}

// 서비스 목록 조회
func showServiceList() {
	services := viper.GetStringMap("services")

	fmt.Printf("============\n")
	fmt.Printf("Service list\n")
	fmt.Printf("============\n")

	for serviceName := range services {
		fmt.Println(serviceName)
	}
}

// 서비스 하위의 액션 목록 조회
func showActionList(serviceName string) {
	spiderActions := viper.GetStringMap("serviceActions." + serviceName)

	fmt.Printf("==============================\n")
	fmt.Printf("[%s] Service Actions list\n", serviceName)
	fmt.Printf("==============================\n")
	for actionName := range spiderActions {
		fmt.Println(actionName)
	}
}

// 입력 값 기반으로 호출할 서비스 정보를 정리함.
func parseRequestInfo() error {
	if serviceName == "" {
		return errors.New("no service is specified to call")
	}

	if !viper.IsSet("services." + serviceName) {
		return errors.New("information about the service [" + serviceName + "] you are trying to call does not exist")
	}

	if actionName == "" {
		return errors.New("no action name is specified to call")
	}

	if !viper.IsSet("serviceActions." + serviceName + "." + actionName) {
		return errors.New("the requested action[" + actionName + "] does not exist for the service[" + serviceName + "] you are trying to call")
	}

	err := viper.UnmarshalKey("services."+serviceName, &serviceInfo)
	if err != nil {
		return err
	}

	if serviceInfo.BaseURL == "" {
		return errors.New("couldn't find the BaseURL information for the service to call")
	}

	if isVerbose {
		fmt.Println("Base URL:", serviceInfo.BaseURL)
		fmt.Println("Auth Type:", serviceInfo.Auth.Type)
		fmt.Println("Username:", serviceInfo.Auth.Username)
		fmt.Println("Password:", serviceInfo.Auth.Password)
	}

	// action 정보 검증

	return nil
}

func init() {
	apiCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "../conf/api.yaml", "config file (default is ../conf/api.yaml)")

	apiCmd.PersistentFlags().StringVarP(&serviceName, "service", "s", "", "Service to perform")
	apiCmd.PersistentFlags().StringVarP(&actionName, "action", "a", "", "Action to perform")
	apiCmd.PersistentFlags().StringVarP(&method, "method", "m", "", "HTTP Method")
	apiCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Show more detail information")

	apiCmd.Flags().BoolVarP(&isListMode, "list", "l", false, "Show Service or Action list")

	cmd.RootCmd.AddCommand(apiCmd)
}
