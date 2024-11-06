package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
)

func main() {
    fmt.Println("Enter your GitHub username:")
    var username string
    fmt.Scanln(&username)

    for {
        fmt.Println("How many events do you want to see? (1-7)")
        var numEvents int
        fmt.Scanln(&numEvents)

        if numEvents < 1 || numEvents > 7 {
            fmt.Println("Please enter a valid number between 1 and 7.")
            continue
        }

        // Construct the URL
        baseURL := "https://api.github.com/users/%s/events"
        fullURL := fmt.Sprintf(baseURL, username)
        parsedUrl, err := url.Parse(fullURL)
        if err != nil {
            fmt.Println("Error parsing the URL:", err)
            continue
        }

        resp, err := http.Get(parsedUrl.String())
        if err != nil {
            fmt.Println("Error making the request:", err)
            continue
        }
        defer resp.Body.Close()

        if resp.StatusCode != 200 {
            fmt.Println("Error: received non-200 response code")
            continue
        }

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            fmt.Println("Error reading the response body:", err)
            continue
        }

        var events []map[string]interface{}
        if err := json.Unmarshal(body, &events); err != nil {
            fmt.Println("Error parsing JSON:", err)
            continue
        }

        if len(events) == 0 {
            fmt.Println("No events found for the user.")
            continue
        } else if len(events) > numEvents {
            events = events[:numEvents]
        }

        fmt.Println("User's recent GitHub events 🚄:")
        for _, event := range events {
            repo := event["repo"].(map[string]interface{})
            payload := event["payload"].(map[string]interface{})

            fmt.Println("Repository ⭐:", repo["name"])
            switch event["type"] {
            case "PushEvent":
                fmt.Println("Event Type: Push ⬆️")
            case "CreateEvent":
                fmt.Println("Event Type: Create ➕")
            case "DeleteEvent":
                fmt.Println("Event Type: Delete 🗑️")
            case "ForkEvent":
                fmt.Println("Event Type: Fork 🍴")
            case "GollumEvent":
                fmt.Println("Event Type: Wiki Edit 📄")
            case "IssueCommentEvent":
                fmt.Println("Event Type: Issue Comment 💬")
            case "IssuesEvent":
                fmt.Println("Event Type: Issue 📝")
            case "MemberEvent":
                fmt.Println("Event Type: Member 👤")
            case "PublicEvent":
                fmt.Println("Event Type: Public 🌍")
            case "PullRequestEvent":
                fmt.Println("Event Type: Pull Request 🔄")
            case "PullRequestReviewEvent":
                fmt.Println("Event Type: PR Review 👀")
            case "PullRequestReviewCommentEvent":
                fmt.Println("Event Type: PR Review Comment 💭")
            case "PullRequestReviewThreadEvent":
                fmt.Println("Event Type: PR Review Thread 🧵")
            case "ReleaseEvent":
                fmt.Println("Event Type: Release 🚀")
            case "SponsorshipEvent":
                fmt.Println("Event Type: Sponsorship 💰")
            case "WatchEvent":
                fmt.Println("Event Type: Watch ⭐")
            default:
                fmt.Println("Event Type: Unknown ❓")
            }

            if commits, ok := payload["commits"].([]interface{}); ok && len(commits) > 0 {
                commit := commits[0].(map[string]interface{})
                fmt.Println("Commit Message ✉️:", commit["message"])
                commitHTMLURL := convertToHTMLURL(commit["url"].(string))
                fmt.Println("Commit URL 🔗:", commitHTMLURL)
            } else {
                fmt.Println("No commits available for this event 🙇🏼‍♂️")
                repoHTMLURL := convertToHTMLURL(repo["url"].(string))
                fmt.Println("Repo URL 🔗:", repoHTMLURL)
            }

            fmt.Println("Created At 🕑:", event["created_at"])
            fmt.Println("-----")
        }

        fmt.Println("Press Enter to continue or Ctrl+C to exit.")
        fmt.Scanln()
    }
}

// Helper function to convert API URL to HTML URL
func convertToHTMLURL(apiURL string) string {
    // Replace "https://api.github.com/repos/" with "https://github.com/"
    // Replace "/commits/" with "/commit/"
    htmlURL := strings.Replace(apiURL, "https://api.github.com/repos/", "https://github.com/", 1)
    htmlURL = strings.Replace(htmlURL, "/commits/", "/commit/", 1)
    return htmlURL
}