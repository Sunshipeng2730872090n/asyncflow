package task

import (
	"errors"
	"github.com/niuniumart/asyncflow/flowsvr/src/ctrl/ctrlmodel"
	"github.com/niuniumart/asyncflow/taskutils/rpc/model"
	"net/http"

	"github.com/niuniumart/asyncflow/flowsvr/src/constant"
	"github.com/niuniumart/asyncflow/flowsvr/src/db"
	"github.com/niuniumart/gosdk/tools"

	"github.com/gin-gonic/gin"
	"github.com/niuniumart/gosdk/handler"
	"github.com/niuniumart/gosdk/martlog"
)

// SetTaskHandler 接口处理handler
type SetTaskHandler struct {
	Req    model.SetTaskReq
	Resp   model.SetTaskResp
	UserId string
}

// SetTask 接口
func SetTask(c *gin.Context) {
	var hd SetTaskHandler
	defer func() {
		hd.Resp.Msg = constant.GetErrMsg(hd.Resp.Code)
		c.JSON(http.StatusOK, hd.Resp)
	}()
	// 获取用户Id
	hd.UserId = c.Request.Header.Get(constant.HEADER_USERID)
	if err := c.ShouldBind(&hd.Req); err != nil {
		martlog.Errorf("SetTask shouldBind err %s", err.Error())
		hd.Resp.Code = constant.ERR_SHOULD_BIND
		return
	}
	handler.Run(&hd)
}

// HandleInput 参数检查
func (p *SetTaskHandler) HandleInput() error {
	if p.Req.TaskId == "" {
		martlog.Errorf("input invalid")
		p.Resp.Code = constant.ERR_INPUT_INVALID
		return constant.ERR_HANDLE_INPUT
	} else {
		p.Req.TaskData.TaskId = p.Req.TaskId
	}
	return nil
}

// HandleProcess 处理函数
func (p *SetTaskHandler) HandleProcess() error {
	var err error
	var Task *db.Task
	Task, err = db.TaskNsp.Find(db.DB, p.Req.TaskId)
	if err != nil {
		martlog.Errorf("db.TaskPosNsp.Find Task Err %s", err.Error())
		p.Resp.Code = constant.ERR_GET_TASK_INFO
		return err
	}
	if Task == nil {
		martlog.Errorf("db.TaskPosNsp.Find Task failed. TaskId:%s", p.Req.TaskId)
		p.Resp.Code = constant.ERR_GET_TASK_INFO
		return errors.New("db.TaskPosNsp.Find Task failed. TaskId:" + p.Req.TaskId)
	}
	// 填充内容
	err = ctrlmodel.FillTaskModel(&p.Req.TaskData, Task, "")
	if err != nil {
		martlog.Errorf("FillTaskModel Err %s. TaskData:%s.Task:%s.", err.Error(),
			tools.GetFmtStr(p.Req.TaskData), tools.GetFmtStr(Task))
		p.Resp.Code = constant.ERR_SET_TASK
		return err
	}
	err = Task.UpdateTask(db.DB)
	if err != nil {
		martlog.Errorf("UpdateTask err %s.Task :%.", err.Error(), tools.GetFmtStr(Task))
		p.Resp.Code = constant.ERR_SET_TASK
		return err
	}
	return nil
}
