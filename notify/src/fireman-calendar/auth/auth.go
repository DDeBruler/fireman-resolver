// Package calendarauth contains utility functions for managing OAuth2 tokens
// used for google calendar API interactions
package calendarauth

import (
  "encoding/json"
  "fmt"
  "net/http"
  "log"
  "bytes"

  "golang.org/x/net/context"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "google.golang.org/api/calendar/v3"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Retrieve a token, saves the token, then returns the generated client.
func GoogleClient() *http.Client {
  credentialBuffer := configFromS3()
  config, err := google.ConfigFromJSON(credentialBuffer, calendar.CalendarReadonlyScope)
  if err != nil {
    log.Fatalf("Unable to parse google config: %v", err)
  }
  return getClient(config)
}

func getClient(config *oauth2.Config) *http.Client {
  tokFile := "token.json"
  tok, err := tokenFromFile(tokFile)
  if err != nil {
    tok = getTokenFromWeb(config)
    saveToken(tokFile, tok)
  }
  return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
  authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
  fmt.Printf("Go to the following link in your browser then type the "+
    "authorization code: \n%v\n", authURL)

  var authCode string
  if _, err := fmt.Scan(&authCode); err != nil {
    log.Fatalf("Unable to read authorization code: %v", err)
  }

  tok, err := config.Exchange(oauth2.NoContext, authCode)
  if err != nil {
    log.Fatalf("Unable to retrieve token from web: %v", err)
  }
  return tok
}

// Retrieves a token from an S3 file.
func tokenFromFile(s3ObjectName string) (*oauth2.Token, error) {
  tokenBody, err := downloadToBuffer(s3ObjectName)
  if err != nil {
    return nil, err
  }
  tok := &oauth2.Token{}
  tokenReader := bytes.NewReader(tokenBody)
  err = json.NewDecoder(tokenReader).Decode(tok)
  return tok, err
}

// Saves a token to a file path.
func saveToken(s3ObjectName string, token *oauth2.Token) {
  var buf bytes.Buffer

  json.NewEncoder(&buf).Encode(token)

  // The session the S3 Uploader will use
  sess := session.Must(session.NewSession())

  // Create an uploader with the session and default options
  uploader := s3manager.NewUploader(sess)

  _, err := uploader.Upload(&s3manager.UploadInput{
    Bucket: aws.String("fireman-resolver"),
    Key:    aws.String(s3ObjectName),
    Body:   &buf,
  })

  if err != nil {
    log.Fatalf("Unable to cache oauth token: %v", err)
  }
}

// Reads google oauth config from S3
func configFromS3() ([]byte) {
  config, err := downloadToBuffer("google-calender-credentials.json")
  if err != nil {
    log.Fatalf("Unable to download OAuth credentials: %v", err)
  }
  return config
}

// Downloads S3 object content to in-memory string
func downloadToBuffer(s3ObjectName string) ([]byte, error) {
  buf := aws.NewWriteAtBuffer([]byte{})

  // Initialize a session that the SDK will use to load credentials
  sess, _ := session.NewSession(&aws.Config{
    Region: aws.String("us-east-1")},
  )

  downloader := s3manager.NewDownloader(sess)

  _, err := downloader.Download(buf,
    &s3.GetObjectInput{
      Bucket: aws.String("fireman-resolver"),
      Key:    aws.String(s3ObjectName),
    })
  return buf.Bytes(), err
}
