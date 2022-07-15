package aws

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	tea "github.com/charmbracelet/bubbletea"
)

type LogGroup struct {
	Name string
}

type LogStream struct {
	Name     string
	Contents string
}

func (l LogGroup) Title() string       { return l.Name }
func (l LogGroup) Description() string { return "" }
func (l LogGroup) FilterValue() string { return l.Name }

type ListLogGroupsMsg struct{}
type LogGroupsMsg []LogGroup
type LoadLogStreamsMsg LogGroup
type LogStreamsMsg []LogStream

type Service struct {
	CloudWatchLogsClient *cloudwatchlogs.Client
}

func ListLogGroups() tea.Msg {
	return ListLogGroupsMsg{}
}

func LoadLogStreams(lg LogGroup) tea.Cmd {
	return func() tea.Msg {
		return LoadLogStreamsMsg(lg)
	}
}

func (s Service) ListLogGroups() tea.Msg {
	opts := cloudwatchlogs.DescribeLogGroupsInput{}
	groups := []LogGroup{}

	paginator := cloudwatchlogs.NewDescribeLogGroupsPaginator(s.CloudWatchLogsClient, &opts)

	for paginator.HasMorePages() {
		res, err := paginator.NextPage(context.TODO())

		if err != nil {
			log.Fatal(err)
		}

		for _, group := range res.LogGroups {
			groups = append(groups, newLogGroup(group))
		}
	}

	return LogGroupsMsg(groups)
}

func (s Service) LoadLogStreams(lg LogGroup) tea.Cmd {
	return func() tea.Msg {
		streamOpts := cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: &lg.Name,
			Descending:   ptrTo(true),
			Limit:        ptrTo[int32](1),
		}
		streams := []LogStream{}

		res, err := s.CloudWatchLogsClient.DescribeLogStreams(context.TODO(), &streamOpts)
		if err != nil {
			log.Fatal(err)
		}

		for _, stream := range res.LogStreams {
			events := []types.OutputLogEvent{}
			eventOpts := cloudwatchlogs.GetLogEventsInput{
				LogGroupName:  &lg.Name,
				LogStreamName: stream.LogStreamName,
				StartFromHead: ptrTo(false),
				Limit:         ptrTo[int32](50),
			}

			eventRes, err := s.CloudWatchLogsClient.GetLogEvents(context.TODO(), &eventOpts)
			if err != nil {
				log.Fatal(err)
			}

			for _, event := range eventRes.Events {
				events = append(events, event)
			}

			streams = append(streams, newLogStream(stream, events))
		}

		return LogStreamsMsg(streams)
	}
}

func newLogGroup(l types.LogGroup) LogGroup {
	return LogGroup{Name: *l.LogGroupName}
}

func newLogStream(l types.LogStream, es []types.OutputLogEvent) LogStream {
	messages := []string{}
	for _, e := range es {
		messages = append(messages, *e.Message)
	}

	return LogStream{Name: *l.LogStreamName, Contents: strings.Join(messages, "\n")}
}

func ptrTo[T any](v T) *T {
	return &v
}
