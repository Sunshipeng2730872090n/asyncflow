package main

import (
	"encoding/json"
	"fmt"
	"github.com/niuniumart/asyncflow/taskutils/constant"
	"github.com/niuniumart/asyncflow/worker/src/initialise"
	"github.com/niuniumart/asyncflow/worker/src/tasksdk"
	"github.com/niuniumart/gosdk/martlog"
	"github.com/niuniumart/gosdk/response"
	"github.com/niuniumart/gosdk/tools"
)

func main() {
	larkTask := tasksdk.TaskHandler{
		TaskType: "lark",
		NewProc:  func() tasksdk.TaskIntf { return new(LarkTask) },
	}
	tasksdk.RegisterHandler(&larkTask)
	initialise.InitResource()
	tasksdk.InitSvr("http://127.0.0.1:41555", "")
	var taskMgr = tasksdk.TaskMgr{
		TaskType: "lark",
	}
	taskMgr.Schedule()
}

type LarkReq struct {
	Msg      string
	FromAddr string
	ToAddr   string
}

// SeriesClaimTask 系列发行任务结构
type LarkTask struct {
	tasksdk.TaskBase
	ContextData *LarkTaskContext
}

// SeriesClaimTaskContext 系列发行任务的上下文
type LarkTaskContext struct {
	ReqBody *LarkReq
	Stage   string
	UserId  string
}

// ContextLoad 解析上下文
func (p *LarkTask) ContextLoad() error {
	martlog.Infof("run lark task %s", p.TaskId)
	err := json.Unmarshal([]byte(p.TaskContext), &p.ContextData)
	if err != nil {
		martlog.Errorf("json unmarshal for context err %s", err.Error())
		return response.RESP_JSON_UNMARSHAL_ERROR
	}
	if p.ContextData.ReqBody == nil {
		p.ContextData.ReqBody = new(LarkReq)
	}
	return nil
}

// HandleProcess 处理函数
func (p *LarkTask) HandleProcess() error {
	fmt.Println("task ", tools.GetFmtStr(*p))
	switch p.TaskStage {
	case "sendmsg":
		p.ContextData.Stage = "record"
		p.SetContextLocal(p.ContextData)
		fallthrough
	case "record":
		fmt.Println("come here")
		p.ContextData.Stage = "record"
		p.Base().Status = int(constant.TASK_STATUS_SUCC)

	default:
		p.Base().Status = int(constant.TASK_STATUS_FAILED)
	}
	return nil
}

// HandleFinish 任务完成函数
func (p *LarkTask) HandleFinish() {
}
