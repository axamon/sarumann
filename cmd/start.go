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
	"net/http"
	"os"

	auth "github.com/abbot/go-http-auth"
	"github.com/axamon/sarumann/server"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

//voicecall gestisce l'invio o meno delle notifiche telefoniche
var voicecall bool

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		authenticator := auth.NewBasicAuthenticator("azure.com", Secret)
		//port := flag.String("port", "8080", "porta TCP da usare.")
		router := mux.NewRouter()
		switch {
		case voicecall == true:
			fmt.Println("Notifiche vocali attive!")
			router.HandleFunc("/create", server.CreateNotifica).Methods("POST")
		default:
			fmt.Println("Notifiche vocali disattivate.")
			router.HandleFunc("/create", server.CreateNotificaNoVoiceCall).Methods("POST")
		}
		router.Handle("/stop", authenticator.Wrap(stop))
		router.HandleFunc("/callfile", server.Callfile)
		router.HandleFunc("/reper", server.SetReper).Methods("POST")
		router.HandleFunc("/getreper/{piatta}", server.GetReper)
		port := cmd.Flag("port").Value.String()
		log.Fatal(http.ListenAndServe(":"+port, router))
		fmt.Println("start called")
	},
}

func stop(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	fmt.Fprintf(w, "<html><body><h1>Sto spegnendo tutto, %s!</h1></body></html>", r.Username)
	os.Exit(0)
}

//Secret gestisce le password
func Secret(user, realm string) string {
	if user == "sarumann" {
		// password is "hello"
		return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
	}
	return ""
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().String("port", "8080", "TCP port to use")
	startCmd.Flags().BoolVarP(&voicecall, "voicecall", "v", true, "Attivazione chiamate vocali")

}
