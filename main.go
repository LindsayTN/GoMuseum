package main

/*Import all packages*/
import (
  "html/template"
  "log"
  "net/http"
  "time"
  "fmt"
  "io/ioutil"
  "os"
  "encoding/json"
)

// Create structs for different page types

type Exhibits struct {
    Exhibits []Exhibit `json:"exhibits"`
}

type Exhibit struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Thumbnail   string `json:"thumbnail"`
}

type MembershipForm struct {
    FirstName    string
    LastName     string
    Email        string
    MembershipLevel string
}

type Visit struct {
    TicketPrice  int
    Address      string
}

type PageVariables struct {
	Date                string
	Time                string
    PageContent         string
    PageName            string
    OpenToday           string
    Exhibits            []Exhibit
    MembershipVariables MembershipForm
    VisitVariables      Visit
}


func main() {
    // Establish routes for the webapp/website
	http.HandleFunc("/", HomePage)
    http.HandleFunc("/visit", VisitPage)
    http.HandleFunc("/exhibits", ExhibitsPage)
    http.HandleFunc("/membership", MembershipPage)
    http.HandleFunc("/confirmation", ConfirmationPage)

    //Ask the server to serve static and CSS files 
    http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
    //http.HandleFunc("/css", serveResource)

    //*Start server
	fmt.Println(http.ListenAndServe(":8000", nil))
}

func HomePage(w http.ResponseWriter, r *http.Request){
    now := time.Now() // find the time right now
    
    //Display if the museum is open today
    DayOfWeek := time.Now().Weekday().String()
    var openMessage string = "open"

    if DayOfWeek == "Monday" {
        openMessage = "closed"
    } else if DayOfWeek == "Tuesday" {
        openMessage = "closed"
    }

    HomeVars := PageVariables{ //store page variables
      Date: now.Format("January 02, 2006"),
      OpenToday: openMessage,
      PageName: "Welcome to the Museum",
      PageContent: "The Metropolitan Museum of Art collects, studies, conserves, and presents significant works of art across time and cultures in order to connect all people to creativity, knowledge, ideas, and one another.",
    }



     t, err := template.ParseFiles("template/layout.html", "template/homepage.html") //parse the html files
    if err != nil { // if there is an error
  	  log.Print("template parsing error: ", err) // log it
  	}
    err = t.Execute(w, HomeVars) //execute the template and pass it the HomeVars struct to fill in the gaps
    if err != nil { // if there is an error
  	  log.Print("template executing error: ", err) //log it
  	}
}

func VisitPage(w http.ResponseWriter, r *http.Request){

    InfoVars := PageVariables{ //store page variables
      PageName: "Visitor's Info",
      PageContent: "The Metropolitan Museum of Art presents over 5,000 years of art from around the world for everyone to experience and enjoy.",
      VisitVariables: Visit{TicketPrice: 5, Address: "123 Main St, Los Angeles, CA 90000"},
    }

    t, err := template.ParseFiles("template/layout.html", "template/visit.html") //parse the html files
    if err != nil { // if there is an error
  	  log.Print("template parsing error: ", err) // log it
  	}
    err = t.Execute(w, InfoVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
    if err != nil { // if there is an error
  	  log.Print("template executing error: ", err) //log it
  	}
}

func ExhibitsPage(w http.ResponseWriter, r *http.Request){

        //Open and read JSON file
        exhibitsData, err := os.Open("data/data.json")
        // if we os.Open returns an error then handle it
        if err != nil {
            fmt.Println(err)
        }
        fmt.Println("Successfully Opened data.json")
        defer exhibitsData.Close()

    byteValue, _ := ioutil.ReadAll(exhibitsData)

    var exhibits Exhibits 

    json.Unmarshal(byteValue, &exhibits)

    for i := 0; i < len(exhibits.Exhibits); i++ {
        fmt.Println("Name: " + exhibits.Exhibits[i].Name)
        fmt.Println("Description: " + exhibits.Exhibits[i].Description)
        fmt.Println("Thumbnail: " + exhibits.Exhibits[i].Thumbnail)
    }

    
    InfoVars := PageVariables{ //store page data
      PageName: "Exhibit Information",
      PageContent: "Here you'll find all of our current exhibits",
      Exhibits: exhibits.Exhibits,
    }

     t, err := template.ParseFiles("template/layout.html", "template/exhibits.html") //parse the html files
    if err != nil { // if there is an error
  	  log.Print("template parsing error: ", err) // log it
  	}
    err = t.Execute(w, InfoVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
    if err != nil { // if there is an error
  	  log.Print("template executing error: ", err) //log it
  	}
} 
func MembershipPage(w http.ResponseWriter, r *http.Request) {
	 MembershipVars := PageVariables{ //store page data
      PageName: "Membership",
      PageContent: "Please sign up for a membership below",
    }

    t, err := template.ParseFiles("template/layout.html", "template/membership.html") //parse the html files
    if err != nil { // if there is an error
  	  log.Print("template parsing error: ", err) // log it
  	}
    err = t.Execute(w, MembershipVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
    if err != nil { // if there is an error
  	  log.Print("template executing error: ", err) //log it
  	}
}
func ConfirmationPage(w http.ResponseWriter, r *http.Request) {
    //Retrieve form values
	 ConfirmationVars := PageVariables{ //store page data
      PageName: "Membership",
      PageContent: "Thank you for signing up for a membership! Here's the information we received'",
      MembershipVariables: MembershipForm{FirstName: r.FormValue("FirstName"), LastName: r.FormValue("LastName"), Email: r.FormValue("Email"), MembershipLevel: r.FormValue("MembershipLevel")},
    }

    /*Extra step: add to a database*/

    t, err := template.ParseFiles("template/layout.html", "template/confirmation.html") //parse the html files
    if err != nil { // if there is an error
  	  log.Print("template parsing error: ", err) // log it
  	}
    err = t.Execute(w, ConfirmationVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
    if err != nil { // if there is an error
  	  log.Print("template executing error: ", err) //log it
  	}
}