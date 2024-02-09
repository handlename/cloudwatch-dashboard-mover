package mover

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fujiwara/tfstate-lookup/tfstate"
	"golang.org/x/net/context"
)

type MoverInput struct {
	DashboardDefinitionPath string

	TFStateFromPath string
	TFStateToPath   string

	AWSProfileFrom string
	AWSPRofileTo   string
}

func ParseMoverInput() *MoverInput {
	input := MoverInput{}

	flag.StringVar(&input.DashboardDefinitionPath, "dashboard-definition", "", "path to the dashboard definition")
	flag.StringVar(&input.TFStateFromPath, "tfstate-from", "", "path to the terraform state file")
	flag.StringVar(&input.TFStateToPath, "tfstate-to", "", "path to the terraform state file")
	flag.StringVar(&input.AWSProfileFrom, "aws-profile-from", "", "aws profile to use for the from state")
	flag.StringVar(&input.AWSPRofileTo, "aws-profile-to", "", "aws profile to use for the to state")
	flag.Parse()

	return &input
}

type MoveTarget struct {
	TFStatePath string
	AWSProfile  string
}

type Mover struct {
	DashboardDefinitionPath string

	From MoveTarget
	To   MoveTarget
}

func NewMover(input *MoverInput) *Mover {
	return &Mover{
		DashboardDefinitionPath: input.DashboardDefinitionPath,
		From: MoveTarget{
			TFStatePath: input.TFStateFromPath,
			AWSProfile:  input.AWSProfileFrom,
		},
		To: MoveTarget{
			TFStatePath: input.TFStateToPath,
			AWSProfile:  input.AWSPRofileTo,
		},
	}
}

func (m *Mover) Replace(ctx context.Context) ([]byte, error) {
	definition, err := os.ReadFile(m.DashboardDefinitionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", m.DashboardDefinitionPath, err)
	}

	tfStateFrom, err := m.readTFState(ctx, m.From.TFStatePath, m.From.AWSProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to read tfstate for FROM from %s: %w", m.From.TFStatePath, err)
	}

	tfStateTo, err := m.readTFState(ctx, m.To.TFStatePath, m.To.AWSProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to read tfstate for TO from %s: %w", m.To.TFStatePath, err)
	}

	// list all resources
	tfNames, err := tfStateFrom.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}

	extractor := NewExtractor()

	for _, tfName := range tfNames {
		// aws id of move from
		awsIDFrom, err := extractor.Extract(tfName, tfStateFrom)
		if err != nil {
			return nil, fmt.Errorf("failed to extract aws id of %s for FROM: %w", tfName, err)
		}
		if awsIDFrom == "" {
			log.Printf("[debug] no aws id found for %s in FROM", tfName)
			continue
		}
		log.Printf("[info] extracted aws id of %s for FROM: %s", tfName, awsIDFrom)

		// aws id of move to
		awsIDTo, err := extractor.Extract(tfName, tfStateTo)
		if err != nil {
			return nil, fmt.Errorf("failed to extract aws id of %s for TO: %w", tfName, err)
		}
		if awsIDTo == "" {
			log.Printf("[debug] no aws id found for %s in TO", tfName)
			continue
		}
		log.Printf("[info] extracted aws id of %s for TO: %s", tfName, awsIDTo)

		// replace id
		log.Printf("[info] replace %s with %s for %s", awsIDFrom, awsIDTo, tfName)
		definition = bytes.ReplaceAll(definition, []byte(awsIDFrom), []byte(awsIDTo))
	}

	// output
	return definition, nil
}

func (m *Mover) readTFState(ctx context.Context, path string, profile string) (*tfstate.TFState, error) {
	return RunWithAWSProfile[*tfstate.TFState](profile, func() (*tfstate.TFState, error) {
		return tfstate.ReadURL(ctx, path)
	})
}
