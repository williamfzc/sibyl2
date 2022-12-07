package server

import (
	"github.com/spf13/cobra"
	"github.com/williamfzc/sibyl2/pkg/server"
	"github.com/williamfzc/sibyl2/pkg/server/object"
)

var serverBackendUrl string
var serverUser string
var serverPwd string
var serverUploadWorkerCount int
var serverUploadQueueSize int

func NewServerCmd() *cobra.Command {
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "sibyl server cmd",
		Long:  `sibyl server cmd`,
		Run: func(cmd *cobra.Command, args []string) {
			config := object.DefaultExecuteConfig()
			if serverBackendUrl != "" {
				config.DbType = object.DtNeo4j
				config.Neo4jUri = serverBackendUrl
			}
			if serverUser != "" {
				config.Neo4jUserName = serverUser
			}
			if serverPwd != "" {
				config.Neo4jPassword = serverPwd
			}
			if serverUploadWorkerCount != 0 {
				config.WorkerCount = serverUploadWorkerCount
			}
			if serverUploadQueueSize != 0 {
				config.WorkerQueueSize = serverUploadQueueSize
			}

			server.Execute(config)
		},
	}
	serverCmd.PersistentFlags().StringVar(&serverBackendUrl, "uri", "", "neo4j backend url")
	serverCmd.PersistentFlags().StringVar(&serverUser, "user", "", "neo4j user")
	serverCmd.PersistentFlags().StringVar(&serverPwd, "pwd", "", "neo4j password")
	serverCmd.PersistentFlags().IntVar(&serverUploadWorkerCount, "workers", 0, "upload worker count")
	serverCmd.PersistentFlags().IntVar(&serverUploadQueueSize, "queueSize", 0, "upload worker count")

	return serverCmd
}
