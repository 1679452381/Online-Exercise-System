package service

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"online_exercise_system/global"
	"online_exercise_system/models"
	"online_exercise_system/response"
	"online_exercise_system/utils"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// SubmitList
// @Tags 公共方法
// @Summary 用户提交列表
// @Param page query string false "page"
// @Param size query string false "size"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query  string false "user_identity"
// @Param status query  string false "status"
// @Success 200 {string} json "{"code":"200","msg":"",data:""}"
// @Router /submit_list [get]
func SubmitList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", global.DefaultPage))
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	size, err := strconv.Atoi(c.DefaultQuery("size", global.DefaultSize))
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	problemIdentity := c.DefaultQuery("problem_identity", "")
	userIdentity := c.DefaultQuery("user_identity", "")
	status := c.DefaultQuery("status", "")

	offset := (page - 1) * size
	//	查数据库
	//count 记录数据的条数
	var count int64
	submits := make([]*models.SubmitBasic, 0)
	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	err = tx.Count(&count).Omit("content").Offset(offset).Limit(size).Find(&submits).Error
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	//	返回结果
	//fmt.Println(problems)
	response.SuccessResponseWithData(gin.H{"list": submits, "count": count}, c)
}

// SubmitCode
// @Tags 用户方法
// @Summary 提交代码
// @Param authorization header string true "authorization"
// @Param problem_identity query string true "problem_identity"
// @Param code body string true "code"
// @Param status query  string false "status"
// @Success 200 {string} json "{"code":"200","msg":"",data:""}"
// @Router /u/admin/submit [post]
func SubmitCode(c *gin.Context) {
	problemIdentity := c.Query("problem_identity")
	if problemIdentity == "" {
		response.FailResponseWithMsg("参数不能为空", c)
		return
	}
	code, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.FailResponseWithMsg("读取代码失败", c)
		return
	}
	//保存代码
	path, err := utils.SaveCode(code)
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	u, _ := c.Get("user_claim")
	uc := u.(*utils.UserClaim)
	sb := &models.SubmitBasic{
		Model:           gorm.Model{},
		Identity:        utils.GetUUID(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    uc.Identity,
		Path:            path,
		Status:          0,
	}
	//fmt.Println(sb)

	//代码判断
	problem := &models.ProblemBasic{}
	err = utils.DB.Where("identity=?", problemIdentity).Preload("TestCases").First(&problem).Error
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	//答案错误 channel
	WA := make(chan int)
	//超内存错误 channel
	OOM := make(chan int)
	//编译错误 channel
	CE := make(chan int)
	//通过个数
	passCount := 0
	var msg string
	var lock sync.Mutex
	for _, testCase := range problem.TestCases {
		go func() {
			cmd := exec.Command("go", "run", path)
			var stderr, out bytes.Buffer
			cmd.Stderr = &stderr //标准错误
			cmd.Stdout = &out    //标准输出
			stdinPipe, err := cmd.StdinPipe()
			if err != nil {
				log.Fatal(err)
			}
			defer stdinPipe.Close()
			io.WriteString(stdinPipe, testCase.Input)

			var bm runtime.MemStats
			runtime.ReadMemStats(&bm)
			//根据测试案例运行代码，拿到结果与标准结果对比
			if err := cmd.Run(); err != nil {
				log.Println(err, stderr.String())
				if err.Error() == "exit status 2" {
					CE <- 1
					msg = stderr.String()
					return
				}
			}
			var em runtime.MemStats
			runtime.ReadMemStats(&em)
			//答案错误
			if testCase.Output != out.String() {
				WA <- 1
				msg = "答案错误"
				return
			}
			//运行超内存
			if (em.Alloc/1024 - bm.Alloc/1024) > uint64(problem.MaxMem) {
				OOM <- 1
				msg = "超内存"
				return
			}
			lock.Lock()
			passCount++
			lock.Unlock()
			msg = "运行通过"
		}()
	}
	//0-待判断 1-答案正确 2-答案错误 3-运行超时 4-超出内存 5-编译错误
	select {
	case <-WA:
		sb.Status = 2
	case <-OOM:
		sb.Status = 4
	case <-CE:
		sb.Status = 5
	case <-time.After(time.Millisecond * time.Duration(problem.MaxRuntime)):
		if passCount == len(problem.TestCases) {
			sb.Status = 1
		} else {
			sb.Status = 3
		}
	}

	//保存数据库
	err = utils.DB.Create(&sb).Error
	if err != nil {
		response.FailResponseWithMsg("提交失败", c)
		return
	}
	response.SuccessResponseWithData(gin.H{
		"status":  sb.Status,
		"msg":     msg,
		"problem": problem,
	}, c)
}
