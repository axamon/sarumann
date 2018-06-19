// Copyright © 2018 Alberto Bregliano <alberto.bregliano@pm.me>
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

	"github.com/spf13/viper"

	"github.com/axamon/sarumann/client"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "cliemt sends data to sarumann server",
	Long:  `sarumann client hostname servizio piattaforma cellreperibile cellulare2(non usato) messaggio`,
	Run: func(cmd *cobra.Command, args []string) {
		endpoint := "http://" + viper.GetString("Server") + "/create"
		fmt.Println(endpoint)
		if len(args) != 6 {
			log.Println("numero di parametri errato")
			os.Exit(1)
		}
		err := client.SendPost(endpoint, args[0], args[1], args[2], args[3], args[4], args[5])
		if err != nil {
			log.Panic(err.Error())
		}
		fmt.Println("client called")
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
