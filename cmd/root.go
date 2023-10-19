package cmd

import (
	"fmt"
	"github.com/xuchaoi/nginxlog-tools/pkg/file"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var logPath string
var detail bool
var logStartTime string
var logEndTime string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nginxlog",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var qps = &cobra.Command{
	Use:   "qps",
	Short: "get nginx qps",
	Long:  "get nginx qps by analysis nginx log",
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()
		err := file.AnalysisLogByLine(logPath, detail, logStartTime, logEndTime)
		if err != nil {
			fmt.Println(err)
		}
		endTime := time.Now()
		fmt.Println("执行时间: ", endTime.Sub(startTime))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nginxlog-tools.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.AddCommand(qps)
	qps.Flags().StringVarP(&logPath, "log-path", "p", "access.log", "指定nginx日志路径")
	qps.Flags().BoolVarP(&detail, "detail", "d", false, "是否展示详细的每分钟请求数据")
	qps.Flags().StringVarP(&logStartTime, "log-start-time", "s", "", "设置分析日志的起始时间")
	qps.Flags().StringVarP(&logEndTime, "log-end-time", "e", "", "设置分析日志的结束时间")
}
