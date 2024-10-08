package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func createCalendarService(userEmail string) (*calendar.Service, error) {
	credentials := os.Getenv("GOOGLE_CALENDAR_CREDENTIALS")

	var credentialsJSON map[string]interface{}

	// Coba parse sebagai JSON object terlebih dahulu
	err := json.Unmarshal([]byte(credentials), &credentialsJSON)
	if err != nil {
		// Jika gagal, coba parse sebagai string JSON
		var credentialsString string
		err = json.Unmarshal([]byte(credentials), &credentialsString)
		if err != nil {
			return nil, fmt.Errorf("unable to parse GOOGLE_CALENDAR_CREDENTIALS: %v", err)
		}
		// Parse string JSON menjadi object
		err = json.Unmarshal([]byte(credentialsString), &credentialsJSON)
		if err != nil {
			return nil, fmt.Errorf("unable to parse GOOGLE_CALENDAR_CREDENTIALS content: %v", err)
		}
	}

	// Konversi kembali ke JSON bytes untuk ConfigFromJSON
	credentialsBytes, err := json.Marshal(credentialsJSON)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal credentials: %v", err)
	}

	config, err := google.ConfigFromJSON(credentialsBytes, calendar.CalendarScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret to config: %v", err)
	}

	ctx := context.Background()
	client := getClient(config, userEmail)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Calendar client: %v", err)
	}

	return srv, nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "rundll32"
		args = append(args, "url.dll,FileProtocolHandler")
	case "darwin":
		cmd = "open"
	default:
		return fmt.Errorf("unsupported platform")
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Printf("Go to the following link in your browser: \n%v\n", authURL)

	err := openBrowser(authURL)
	if err != nil {
		log.Fatalf("Unable to open browser: %v", err)
	}

	// Start a local server to handle the callback
	codeChan := make(chan string)
	srv := &http.Server{Addr: "localhost:4040"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			fmt.Fprintf(w, "Authorization code received. You can close this window.")
			codeChan <- code
		} else {
			fmt.Fprintf(w, "No authorization code received.")
		}
	})
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Unable to start server: %v", err)
		}
	}()

	code := <-codeChan
	srv.Shutdown(context.TODO())

	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func getClient(config *oauth2.Config, userEmail string) *http.Client {
	tokFile := fmt.Sprintf("token_%s.json", userEmail)
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func CreateGoogleCalendarEvent(senderEmail, summary, description, startDateTime, endDateTime, timeZone string, attendees []string) (*calendar.Event, error) {
	srv, err := createCalendarService(senderEmail)
	if err != nil {
		return nil, err
	}

	event := &calendar.Event{
		Summary:     summary,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: startDateTime,
			TimeZone: timeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: endDateTime,
			TimeZone: timeZone,
		},
		Attendees: make([]*calendar.EventAttendee, 0, len(attendees)),
	}

	// Hanya tambahkan attendee jika bukan pengirim
	for _, email := range attendees {
		if email != senderEmail {
			event.Attendees = append(event.Attendees, &calendar.EventAttendee{Email: email})
		}
	}

	calendarId := "primary"
	event, err = srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to create event: %v", err)
	}

	return event, nil
}
