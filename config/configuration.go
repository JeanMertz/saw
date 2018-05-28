package config

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

const maxLimit = 50

type Configuration struct {
	Group      string
	Prefix     string
	Start      string
	End        string
	Filter     string
	Region     string
	Descending bool
	OrderBy    string
}

func getTime(timeStr string) (time.Time, error) {
	relative, err := time.ParseDuration(timeStr)
	if err == nil {
		return time.Now().Add(relative), nil
	}

	absolute, err := time.Parse(time.RFC3339, timeStr)
	if err == nil {
		return absolute, nil
	}

	return time.Time{}, errors.New("Could not parse relative or absolute time")
}

func (c *Configuration) DescribeLogGroupsInput() *cloudwatchlogs.DescribeLogGroupsInput {
	input := cloudwatchlogs.DescribeLogGroupsInput{}
	if c.Prefix != "" {
		input.SetLogGroupNamePrefix(c.Prefix)
	}
	return &input
}

func (c *Configuration) DescribeLogStreamsInput() *cloudwatchlogs.DescribeLogStreamsInput {
	input := cloudwatchlogs.DescribeLogStreamsInput{}
	input.SetLogGroupName(c.Group)
	input.SetDescending(c.Descending)
	input.SetOrderBy(c.OrderBy)
	if c.Prefix != "" {
		input.SetLogStreamNamePrefix(c.Prefix)
	}
	return &input
}

func (c *Configuration) FilterLogEventsInput() *cloudwatchlogs.FilterLogEventsInput {
	input := cloudwatchlogs.FilterLogEventsInput{}
	input.SetInterleaved(true)
	input.SetLogGroupName(c.Group)

	absoluteStartTime := time.Now()
	if c.Start != "" {
		st, err := getTime(c.Start)
		if err == nil {
			absoluteStartTime = st
		}
	}
	input.SetStartTime(aws.TimeUnixMilli(absoluteStartTime))

	if c.End != "" {
		et, err := getTime(c.End)
		if err == nil {
			input.SetEndTime(aws.TimeUnixMilli(et))
		}
	}

	if len(c.Filter) != 0 {
		input.SetFilterPattern(c.Filter)
	}

	return &input
}
