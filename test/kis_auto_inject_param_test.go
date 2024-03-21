package test

import (
	"context"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/flow"
	"kis-flow/kis"
	"kis-flow/test/faas"
	"kis-flow/test/proto"
	"testing"
)

func TestAutoInjectParam(t *testing.T) {
	ctx := context.Background()

	kis.Pool().FaaS("AvgStuScore", faas.AvgStuScore)
	kis.Pool().FaaS("PrintStuAvgScore", faas.PrintStuAvgScore)

	source1 := config.KisSource{
		Name: "Test",
		Must: []string{},
	}

	avgStuScoreConfig := config.NewFuncConfig("AvgStuScore", common.C, &source1, nil)
	if avgStuScoreConfig == nil {
		panic("AvgStuScore is nil")
	}

	printStuAvgScoreConfig := config.NewFuncConfig("PrintStuAvgScore", common.C, &source1, nil)
	if printStuAvgScoreConfig == nil {
		panic("printStuAvgScoreConfig is nil")
	}

	myFlowConfig1 := config.NewFlowConfig("cal_stu_avg_score", common.FlowEnable)

	flow1 := flow.NewKisFlow(myFlowConfig1)

	// 4. 拼接Functioin 到 Flow 上
	if err := flow1.Link(avgStuScoreConfig, nil); err != nil {
		panic(err)
	}
	if err := flow1.Link(printStuAvgScoreConfig, nil); err != nil {
		panic(err)
	}

	// 3. 提交原始数据
	_ = flow1.CommitRow(&faas.AvgStuScoreIn{
		proto.StuScores{
			StuId:  100,
			Score1: 1,
			Score2: 2,
			Score3: 3,
		},
	})
	_ = flow1.CommitRow(`{"stu_id":101}`)
	_ = flow1.CommitRow(faas.AvgStuScoreIn{
		proto.StuScores{
			StuId:  100,
			Score1: 1,
			Score2: 2,
			Score3: 3,
		},
	})

	// 4. 执行flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}
