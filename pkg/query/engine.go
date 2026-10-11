package query

import (
	"strconv"
	"strings"

	"github.com/sujanv/LogSift/pkg/filter"
	"github.com/sujanv/LogSift/pkg/model"
)

type PipelineExecutor struct {
	filterPred filter.Predicate
	topK       *TopKAggregator
}

func CompilePipeline(stages []PipelineStage) (*PipelineExecutor, error) {
	pe := &PipelineExecutor{}

	for _, s := range stages {
		switch strings.ToLower(s.Command) {
		case "filter":
			if len(s.Args) > 0 {
				pred, err := filter.ParseExpression(s.Args[0])
				if err != nil {
					return nil, err
				}
				pe.filterPred = pred
			}
		case "topk":
			if len(s.Args) > 0 {
				field := s.Args[0]
				limit := 10
				if len(s.Args) > 1 {
					if n, err := strconv.Atoi(s.Args[1]); err == nil {
						limit = n
					}
				}
				pe.topK = NewTopKAggregator(field, limit)
			}
		}
	}

	return pe, nil
}

func (pe *PipelineExecutor) Process(r *model.Record) bool {
	if pe.filterPred != nil && !pe.filterPred.Match(r) {
		return false
	}
	if pe.topK != nil {
		pe.topK.Ingest(r)
	}
	return true
}

func (pe *PipelineExecutor) TopKResult() []TopKItem {
	if pe.topK != nil {
		return pe.topK.Result()
	}
	return nil
}
