package jira

import (
	"fmt"
	"os"

	"github.com/andygrunwald/go-jira"
	gojira "github.com/andygrunwald/go-jira"
	"github.com/google/go-github/v47/github"
)

// func PrintGithubIssue(issue *github.Issue, oneline bool, color bool) {
//
//     if oneline {
//         if color {
//             // print the idea in yellow, then reset the rest of the line
//             fmt.Printf("\033[33m%5d\033[0m \033[32m%s\033[0m %s\n", issue.GetNumber(), issue.GetState(), issue.GetTitle())
//         } else {
//             fmt.Printf("%5d %s %s\n", issue.GetNumber(), issue.GetState(), issue.GetTitle())
//         }
//     } else {
//         // fmt.Println(*issue.ID)
//         fmt.Printf("Issue:\t%d\n", issue.GetNumber())
//         // fmt.Println(*issue.Title)
//         fmt.Printf("State:\t%s\n", issue.GetState())
//         if issue.GetAssignee() != nil {
//             fmt.Printf("Assignee:\t%s\n", *issue.GetAssignee().Login)
//         }
//
//         // NOTE: This should be the jira body
//         // fmt.Printf("Title:\t%s\n", issue.GetTitle())
//         fmt.Printf("\n   %s\n\n", issue.GetTitle())
//         // fmt.Printf("Body:\n\t%s\n", issue.GetBody())
//
//         // Look through the labels
//         // important soon should be Urgent
//         // important long term should be Medium
//         // fmt.Println(issue.Labels)
//     }
// }

func getToken() string {
	token, ok := os.LookupEnv("JIRA_TOKEN")
	if !ok {
		fmt.Println("please supply your JIRA_TOKEN")
		os.Exit(1)
	}
	return token
}

func CloneIssueToJira(issue *github.Issue, dryRun bool) {
	token := getToken()

	tp := gojira.BearerAuthTransport{
		Token: token,
	}

	// tp := gojira.BasicAuthTransport{
	//     Username: "username",
	//     Password: "token",
	// }

	jiraClient, err := gojira.NewClient(tp.Client(), "https://issues.redhat.com")
	if err != nil {
		panic(err)
	}

	ji := jira.Issue{
		Fields: &gojira.IssueFields{
			// Assignee: &gojira.User{
			//     Name: "myuser",
			// },
			// Reporter: &gojira.User{
			//     Name: "youruser",
			// },
			Description: fmt.Sprintf("%s\n\nUpstream Github issue: %s\n", issue.GetBody(), issue.GetURL()),
			Type: gojira.IssueType{
				Name: "Story",
			},
			Project: gojira.Project{
				Key: "OSDK",
			},
			Summary: issue.GetTitle(),
		},
	}

	if dryRun {
		fmt.Printf("dryrun: Cloning %d to jira\n", issue.GetNumber())
		fmt.Printf("%#v\n", &ji)
		fmt.Println(ji.Fields.Summary)
		fmt.Println(ji.Fields.Description)
	} else {
		fmt.Printf("FORREALZ! Cloning %d to jira\n", issue.GetNumber())
		daIssue, _, err := jiraClient.Issue.Create(&ji)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		// TODO there's a nil pointer from this call
		fmt.Printf("%s: %+v\n", daIssue.Key, daIssue.Fields.Summary)
	}

}
