package mover

import (
	"fmt"
	"log"
	"strings"

	"github.com/fujiwara/tfstate-lookup/tfstate"
)

type Extractor struct {
	handlers []ExtractHandler
}

type ExtractHandler struct {
	TargetName string // terraform resource name
	Extract    func(*tfstate.Object) string
}

func (ext *Extractor) Extract(tfName string, state *tfstate.TFState) (string, error) {
	for _, h := range ext.handlers {
		if !strings.HasPrefix(tfName, h.TargetName) {
			continue
		}

		obj, err := state.Lookup(tfName)
		if err != nil {
			return "", fmt.Errorf("failed to lookup %s: %w", tfName, err)
		}

		return h.Extract(obj), nil
	}

	log.Printf("[debug] no resource type matched for %s", tfName)

	return "", nil
}

func NewExtractor() *Extractor {
	ext := Extractor{}

	ext.handlers = []ExtractHandler{
		{
			TargetName: "aws_cloudfront_distribution",
			Extract: func(o *tfstate.Object) string {
				return o.Value.(map[string]interface{})["id"].(string)
			},
		},
		{
			TargetName: "aws_lb_target_group",
			Extract: func(o *tfstate.Object) string {
				return o.Value.(map[string]interface{})["arn_suffix"].(string)
			},
		},
	}

	return &ext
}
