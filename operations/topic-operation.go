package operations

import (
	"github.com/spf13/cobra"
)

type Topic struct {
	Name       string
	Partitions []Partition `json:",omitempty" yaml:",omitempty"`
	Configs    []Config    `json:",omitempty" yaml:",omitempty"`
}

type Partition struct {
	Id           int32
	OldestOffset int64   `json:"oldestOffset" yaml:"oldestOffset"`
	NewestOffset int64   `json:"newestOffset" yaml:"newestOffset"`
	Leader       string  `json:",omitempty" yaml:",omitempty"`
	Replicas     []int32 `json:",omitempty" yaml:",omitempty,flow"`
	ISRs         []int32 `json:"inSyncReplicas,omitempty" yaml:"inSyncReplicas,omitempty,flow"`
}

type requestedTopicFields struct {
	partitionId       bool
	partitionOffset   bool
	partitionLeader   bool
	partitionReplicas bool
	partitionISRs     bool
	config            bool
}

var allFields = requestedTopicFields{partitionId: true, partitionOffset: true, partitionLeader: true, partitionReplicas: true, partitionISRs: true, config: true}

type Config struct {
	Name  string
	Value string
}

type GetTopicsFlags struct {
	OutputFormat string
}

type CreateTopicFlags struct {
	Partitions        int32
	ReplicationFactor int16
	ValidateOnly      bool
	Configs           []string
}

type AlterTopicFlags struct {
	Partitions        int32
	ReplicationFactor int16
	ValidateOnly      bool
	Configs           []string
}

type DescribeTopicFlags struct {
	PrintConfigs bool
	OutputFormat string
}

type TopicOperation struct {
}

func (operation *TopicOperation) DescribeTopic(topic string, flags DescribeTopicFlags) error {
	return nil
}

func CompleteTopicNames(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {

	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	topics, err := (&TopicOperation{}).ListTopicsNames()

	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return topics, cobra.ShellCompDirectiveNoFileComp
}

func (operation *TopicOperation) ListTopicsNames() ([]string, error) {
	topics := []string{"Test"}

	return topics, nil
}
